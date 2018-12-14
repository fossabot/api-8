package database

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConnectionSet_Quit(t *testing.T) {
	masterConn := &connection{
		connected: true,
		host:      "master",
		quitCh:    make(chan bool),
		s:         &dummySQL{},
		pingFn: func() error {
			return nil
		},
	}
	go masterConn.loop()
	slaveConns := []*connection{
		&connection{
			connected: true,
			host:      "slave1",
			quitCh:    make(chan bool),
			s:         &dummySQL{},
		},
	}
	go slaveConns[0].loop()
	cs := &connectionSet{
		masterConn: masterConn,
		slaveConns: slaveConns,
		nSlaves:    1,
		quitCh:     make(chan bool),
		asMutex:    new(sync.RWMutex),
	}
	go cs.loop()

	cs.quit()
	time.Sleep(5 * time.Millisecond)
	assert.False(t, masterConn.connected, "failed to close master connection")
	for _, conn := range slaveConns {
		assert.False(t, conn.connected, "failed to close slave connection")
	}
}

func TestConnectionSet_UpdateActiveSlaves(t *testing.T) {
	slaveConns := []*connection{
		&connection{
			connected: true,
			host:      "slave1",
			quitCh:    make(chan bool),
			s:         &dummySQL{},
			pingFn: func() error {
				return nil
			},
		},
	}
	go slaveConns[0].loop()

	cs := &connectionSet{
		masterConn: nil,
		slaveConns: slaveConns,
		nSlaves:    1,
		quitCh:     make(chan bool),
		asMutex:    new(sync.RWMutex),
	}
	cs.updateActiveSlaves()

	assert.Equal(t, cs.activeSlaves[0], cs.slaveConns[0], "activeSlaves should have the same value with salveConns after running updateActiveSlaves")
}

func TestConnectionSet_Loop(t *testing.T) {
	csTickDelay = 1 * time.Millisecond
	slaveConns := []*connection{
		&connection{
			connected: true,
			host:      "slave1",
			quitCh:    make(chan bool),
			s:         &dummySQL{},
			pingFn: func() error {
				return nil
			},
		},
	}
	go slaveConns[0].loop()
	cs := &connectionSet{
		masterConn: nil,
		slaveConns: slaveConns,
		nSlaves:    1,
		quitCh:     make(chan bool),
		asMutex:    new(sync.RWMutex),
	}

	closed := false
	go func() {
		<-cs.quitCh
		closed = true
	}()
	go cs.loop()

	time.Sleep(5 * time.Millisecond)
	assert.Equal(t, cs.activeSlaves[0], cs.slaveConns[0], "loop does not update activeSlaves")

	cs.quit()
	time.Sleep(5 * time.Millisecond)
	assert.True(t, closed, "loop does not close connectionSet after calling quit")
}

func TestConnectionSet_LoopConnStatusChanges(t *testing.T) {
	csTickDelay = 1 * time.Millisecond
	slaveConns := []*connection{
		&connection{
			connected: true,
			host:      "slave1",
			quitCh:    make(chan bool),
			s:         &dummySQL{},
			pingFn: func() error {
				return nil
			},
		},
		&connection{
			connected: true,
			host:      "slave2",
			quitCh:    make(chan bool),
			s:         &dummySQL{},
			pingFn: func() error {
				return nil
			},
		},
	}
	cs := &connectionSet{
		masterConn: nil,
		slaveConns: slaveConns,
		nSlaves:    2,
		quitCh:     make(chan bool),
		asMutex:    new(sync.RWMutex),
	}
	go cs.loop()

	// at the beginning there are 2 active slaves
	time.Sleep(5 * time.Millisecond)
	assert.Len(t, cs.activeSlaves, 2, "all connected slaves must be in active slaves")

	// then 1 of them died
	slaveConns[0].connected = false
	time.Sleep(5 * time.Millisecond)
	assert.Len(t, cs.activeSlaves, 1, "should just have 1 active slaves")

	// the last 1 dies too :(
	slaveConns[1].connected = false
	time.Sleep(5 * time.Millisecond)
	assert.Len(t, cs.activeSlaves, 0, "should just have 0 active slaves")

	// then all of them alive again
	slaveConns[0].connected = true
	slaveConns[1].connected = true
	time.Sleep(5 * time.Millisecond)
	assert.Len(t, cs.activeSlaves, 2, "all slaves should be connected")
}

func TestConnectionSet_Reader(t *testing.T) {
	masterConn := &connection{
		connected: true,
		host:      "master",
		quitCh:    make(chan bool),
		s:         &dummySQL{},
		pingFn: func() error {
			return nil
		},
	}
	slaveConns := []*connection{
		&connection{
			connected: true,
			host:      "slave1",
			quitCh:    make(chan bool),
			s:         &dummySQL{},
		},
	}
	cs := &connectionSet{
		masterConn: masterConn,
		slaveConns: slaveConns,
		nSlaves:    1,
		quitCh:     make(chan bool),
		asMutex:    new(sync.RWMutex),
	}
	cs.updateActiveSlaves()

	// should return sql from slave
	assert.NotNil(t, cs.reader(), "nil reader")
	assert.Equal(t, cs.reader(), slaveConns[0].s, "reader does not return sql from connection from slaveConns")

	// all slaves are dead, should return sql from master
	cs.slaveConns = nil
	cs.updateActiveSlaves()
	assert.NotNil(t, cs.reader(), "nil reader")
	assert.Equal(t, cs.reader(), masterConn.s, "reader does not return masterConn sql when all slaves are gone")
}

func TestConnectionSet_Writer(t *testing.T) {
	masterConn := &connection{
		connected: true,
		host:      "master",
		quitCh:    make(chan bool),
		s:         &dummySQL{},
		pingFn: func() error {
			return nil
		},
	}
	slaveConns := []*connection{
		&connection{
			connected: true,
			host:      "slave1",
			quitCh:    make(chan bool),
			s:         &dummySQL{},
		},
	}
	cs := &connectionSet{
		masterConn: masterConn,
		slaveConns: slaveConns,
		nSlaves:    1,
		quitCh:     make(chan bool),
		asMutex:    new(sync.RWMutex),
	}
	cs.updateActiveSlaves()

	assert.NotNil(t, cs.writer(), "nil writer")
	assert.Equal(t, cs.writer(), masterConn.s, "writer does not return masterConn sql")

	cs.slaveConns = nil
	cs.updateActiveSlaves()
	assert.NotNil(t, cs.writer(), "nil writer")
	assert.Equal(t, cs.writer(), masterConn.s, "writer does not return masterConn sql")
}
