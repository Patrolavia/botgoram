// This file is part of Botgoram
// Botgoram is free software: see LICENSE.txt for more details.

package botgoram

import (
	"sync"

	"github.com/Patrolavia/telegram"
)

type manager struct {
	size         int
	runningUsers map[string]bool
	msgQ         []*telegram.Message
	lock         sync.Locker
	cond         *sync.Cond
	getUID       func(*telegram.Message) *telegram.Victim
	msgs         chan *telegram.Message
}

func newManager(f func(*telegram.Message) *telegram.Victim, size int, msgs chan *telegram.Message) *manager {
	l := new(sync.Mutex)
	return &manager{
		size,
		make(map[string]bool),
		nil,
		l,
		sync.NewCond(l),
		f,
		msgs,
	}
}

func (m *manager) GetUID(msg *telegram.Message) *telegram.Victim {
	return m.getUID(msg)
}

func (m *manager) Commit(msg *telegram.Message) {
	m.lock.Lock()
	delete(m.runningUsers, m.getUID(msg).Identifier())
	// delete msg from Q
	for k, mm := range m.msgQ {
		if mm.ID == msg.ID {
			tmp := m.msgQ[0:k]
			if k < len(m.msgQ)-1 {
				tmp = append(tmp, m.msgQ[k+1:]...)
			}
			m.msgQ = tmp
		}
	}
	m.lock.Unlock()
	m.cond.Signal()
}

func (m *manager) Rollback(msg *telegram.Message) {
	if ok := m.runningUsers[m.getUID(msg).Identifier()]; !ok {
		return
	}
	m.lock.Lock()
	delete(m.runningUsers, m.getUID(msg).Identifier())
	m.lock.Unlock()
	m.cond.Signal()
}

func (m *manager) getFirstNew() (ret *telegram.Message) {
	for _, msg := range m.msgQ {
		if !m.runningUsers[m.getUID(msg).Identifier()] {
			m.runningUsers[m.getUID(msg).Identifier()] = true
			return msg
		}
	}
	return
}

func (m *manager) Begin() *telegram.Message {
	m.lock.Lock()
	defer m.lock.Unlock()

	msg := m.getFirstNew()
	for ; msg == nil; msg = m.getFirstNew() {
		m.cond.Wait()
	}
	return msg
}

func (m *manager) Run() {
	for msg := range m.msgs {
		m.feed([]*telegram.Message{msg})
	}
}

func (m *manager) feed(msgs []*telegram.Message) {
	m.lock.Lock()
	m.msgQ = msgs
	m.lock.Unlock()
	m.cond.Signal()

	m.lock.Lock()
	defer m.lock.Unlock()
	for len(m.msgQ) > 0 || len(m.runningUsers) >= m.size {
		m.cond.Wait()
	}
}
