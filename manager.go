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
		nil,
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

func (m *manager) Run(api telegram.API, timeout int) {
	offset := 0
	// Set limit to double of max goroutines
	// This can improve performance without using too much memory.
	limit := m.size * 2
	for {
		updates, err := api.GetUpdates(offset, limit, timeout)

		// just try again if no pending message or any error
		if err != nil || len(updates) < 1 {
			continue
		}

		msgs := make([]*telegram.Message, len(updates))
		for k, v := range updates {
			msgs[k] = v.Message
			if offset <= v.Id {
				offset = v.Id + 1
			}
		}
		m.feed(msgs)
	}
}

func (m *manager) feed(msgs []*telegram.Message) {
	m.lock.Lock()
	m.msgQ = msgs
	m.lock.Unlock()
	m.cond.Signal()

	m.lock.Lock()
	defer m.lock.Unlock()
	for len(m.msgQ) > 0 || len(m.running_users) >= m.size {
		m.cond.Wait()
	}
}
