package database

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConnection_PingOK(t *testing.T) {
	connTickDelay = 1 * time.Millisecond
	c := &connection{
		host:      "dummy",
		connected: false,
		quitCh:    make(chan bool),
		s:         &dummySQL{},
		pingFn: func() error {
			return nil
		},
	}
	go c.loop()

	time.Sleep(3 * time.Millisecond)
	assert.True(t, c.connected, "should be connected=true after loop, which calls pingFn, runs")
}

func TestConnection_PingFail(t *testing.T) {
	connTickDelay = 1 * time.Millisecond
	c := &connection{
		host:      "dummy",
		connected: true,
		quitCh:    make(chan bool),
		s:         &dummySQL{},
		pingFn: func() error {
			return errors.New("fail")
		},
	}
	go c.loop()

	time.Sleep(2 * time.Millisecond)
	assert.False(t, c.connected, "should be connected=false after loop, which calls pingFn, runs")
}

func TestConnection_Loop(t *testing.T) {
	connTickDelay = 1 * time.Millisecond
	ok := true
	c := &connection{
		host:      "dummy",
		connected: true,
		quitCh:    make(chan bool),
		s:         &dummySQL{},
		pingFn: func() error {
			if ok {
				return nil
			}
			return errors.New("fail")
		},
	}
	go c.loop()

	// at the beginning, the connection is ok
	time.Sleep(2 * time.Millisecond)
	assert.True(t, c.connected, "should be connected=true")

	// then it fails
	ok = false
	time.Sleep(2 * time.Millisecond)
	assert.False(t, c.connected, "should be connected=false")

	// then it is ok again
	ok = true
	time.Sleep(2 * time.Millisecond)
	assert.True(t, c.connected, "should be connected=true again")
}

func TestConnection_Quit(t *testing.T) {
	var closed bool
	c := &connection{
		host:      "dummy",
		connected: true,
		quitCh:    make(chan bool),
		s:         &dummySQL{},
		pingFn: func() error {
			return errors.New("fail")
		},
		closeFn: func() {
			closed = true
		},
	}
	go c.loop()

	c.quit()
	time.Sleep(10 * time.Millisecond) // wait quit to finish running
	assert.True(t, closed, "quit does not update closed")
	assert.False(t, c.connected, "still connected after quit")
}

func TestWrapPingFn_WithTimeoutFn(t *testing.T) {
	pingTimeout = 5 * time.Millisecond
	timeoutFn := wrapPingFn(func() error {
		time.Sleep(10 * time.Millisecond)
		return errors.New("something wrong")
	})
	err := timeoutFn()
	assert.Equal(t, err, errPingTimeout, "should return error ping timeout")
}

func TestWrapPingFn_WithNonTimeoutFn(t *testing.T) {
	pingTimeout = 10 * time.Millisecond

	pingErr := errors.New("something wrong")
	nonTimeoutFn := wrapPingFn(func() error {
		time.Sleep(5 * time.Millisecond)
		return pingErr
	})
	err := nonTimeoutFn()
	assert.Equal(t, err, pingErr, "should return error from pingFn")
}

func TestWrapPingFn_WithSuccessFn(t *testing.T) {
	pingTimeout = 10 * time.Millisecond
	okFn := wrapPingFn(func() error {
		return nil
	})
	err := okFn()
	assert.Nil(t, err, "should return nil error")
}
