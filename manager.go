package botgoram

import (
	"sync"

	"github.com/Patrolavia/botgoram/telegram"
)

type manager struct {
	size          int
	running_users map[int]bool
	msgQ          []*telegram.Message
	lock          sync.Locker
	cond          *sync.Cond
	getUid        func(*telegram.Message) *telegram.User
}

func newManager(f func(*telegram.Message) *telegram.User, size int) *manager {
	l := new(sync.Mutex)
	return &manager{
		size,
		make(map[int]bool),
		make([]*telegram.Message, 0, size),
		l,
		sync.NewCond(l),
		f,
	}
}

func (m *manager) GetUid(msg *telegram.Message) *telegram.User {
	return m.getUid(msg)
}

func (m *manager) Commit(msg *telegram.Message) {
	m.lock.Lock()
	delete(m.running_users, m.getUid(msg).Id)
	// delete msg from Q
	for k, mm := range m.msgQ {
		if mm.Id == msg.Id {
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
	if ok := m.running_users[m.getUid(msg).Id]; !ok {
		return
	}
	m.lock.Lock()
	delete(m.running_users, m.getUid(msg).Id)
	m.lock.Unlock()
	m.cond.Signal()
}

func (m *manager) getFirstNew() (ret *telegram.Message) {
	for _, msg := range m.msgQ {
		if !m.running_users[m.getUid(msg).Id] {
			m.running_users[m.getUid(msg).Id] = true
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

func (m *manager) feed(msg *telegram.Message) {
	defer m.cond.Signal()
	m.lock.Lock()
	defer m.lock.Unlock()
	for len(m.msgQ) >= m.size {
		m.cond.Wait()
	}
	m.msgQ = append(m.msgQ, msg)
}
