package database

import (
	"math/rand"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var csTickDelay = 30 * time.Second

// connectionSet abstracts all database connections we are currently possessing
// there are one connection to master and n number connections to slaves
type connectionSet struct {
	masterConn *connection
	slaveConns []*connection
	nSlaves    int // so we don't neet to always count slaveConns

	asMutex      *sync.RWMutex
	activeSlaves []*connection

	quitCh chan bool
}

func newConnectionSet(masterConn *connection, slaveConns []*connection) *connectionSet {
	cs := &connectionSet{
		masterConn:   masterConn,
		slaveConns:   slaveConns,
		nSlaves:      len(slaveConns),
		asMutex:      new(sync.RWMutex),
		activeSlaves: nil,
		quitCh:       make(chan bool),
	}
	cs.updateActiveSlaves()
	go cs.loop()
	return cs
}

func (c *connectionSet) loop() {
	ticker := time.NewTicker(csTickDelay)
	for {
		select {
		case <-ticker.C:
			c.updateActiveSlaves()
		case <-c.quitCh:
			if c.masterConn != nil {
				c.masterConn.quit()
			}
			for _, conn := range c.slaveConns {
				conn.quit()
			}
			return
		}
	}
}

func (c *connectionSet) updateActiveSlaves() {
	slaves := make([]*connection, 0, c.nSlaves)
	for i, conn := range c.slaveConns {
		conn.cMutex.RLock()
		if conn.connected {
			slaves[i] = conn
		}
		conn.cMutex.RUnlock()
	}

	// TODO: should add a simple mechanism to detect changes in activeSlaves
	// so we don't have to always lock read access to activeSlaves
	// everytime this method runs

	c.asMutex.Lock()
	c.activeSlaves = slaves
	c.asMutex.Unlock()
}

func (c *connectionSet) reader() sql {
	nActive := len(c.activeSlaves)
	if nActive == 0 {
		logrus.Warnln("no active database connection to slaves")
		return c.masterConn.s
	}

	c.asMutex.RLock()
	slaveConn := c.activeSlaves[rand.Intn(nActive)]
	c.asMutex.RUnlock()
	return slaveConn.s
}

func (c *connectionSet) writer() sql {
	return c.masterConn.s
}

func (c *connectionSet) quit() {
	close(c.quitCh)
}
