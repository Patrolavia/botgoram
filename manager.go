// This file is part of Botgoram
// Botgoram is free software: see LICENSE.txt for more details.

package botgoram

import (
	"log"
	"sync"

	"github.com/Patrolavia/telegram"
)

type msgq struct {
	msg  *telegram.Message
	next *msgq
}

type manager struct {
	size         int
	runningUsers map[string]bool
	root         *msgq
	qsize        int
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
		0,
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
	defer m.cond.Signal()
	defer m.lock.Unlock()

	delete(m.runningUsers, m.getUID(msg).Identifier())
	// delete msg from Q
	if m.root == nil {
		log.Fatal("botgoram: There is no queued message to be deleted! There must be something wrong in botgoram.")
	}

	if m.root.msg == msg {
		m.root = m.root.next
		m.qsize--
		return
	}

	prev := m.root
	cur := m.root.next
	for cur != nil {
		if cur.msg != msg {
			prev = cur
			cur = cur.next
			continue
		}

		prev.next = cur.next
		m.qsize--
		return
	}

	log.Fatal("botgoram: I can't find matched message to delete from queue. There must be something wrong in botgoram.")
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
	cur := m.root

	for cur != nil {
		if !m.runningUsers[m.getUID(cur.msg).Identifier()] {
			m.runningUsers[m.getUID(cur.msg).Identifier()] = true
			return cur.msg
		}
		cur = cur.next
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
		m.feed(msg)
	}
}

func (m *manager) add(msg *telegram.Message) {
	if m.root == nil {
		m.root = &msgq{msg, nil}
		return
	}

	cur := m.root
	for cur.next != nil {
		cur = cur.next
	}
	cur.next = &msgq{msg, nil}
}

func (m *manager) feed(msg *telegram.Message) {
	m.lock.Lock()
	m.add(msg)
	m.qsize++
	m.lock.Unlock()
	m.cond.Signal()

	m.lock.Lock()
	defer m.lock.Unlock()
	for m.qsize >= m.size || len(m.runningUsers) >= m.size {
		m.cond.Wait()
	}
}
