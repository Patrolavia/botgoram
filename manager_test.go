package botgoram

import (
	"math/rand"
	"testing"

	"github.com/Patrolavia/botgoram/telegram"
)

func makeTestUser(name string) *telegram.User {
	return &telegram.User{
		Id:        rand.Int(),
		FirstName: name,
	}
}

// test with 2 worker, 2 sender, no overflow, no block
func TestManagerWithTwoUser(t *testing.T) {
	u1 := makeTestUser("user1")
	u2 := makeTestUser("user2")

	m := newManager(bySender, 2)
	m1 := &telegram.Message{
		Id:     1,
		Text:   "test",
		Sender: u1,
		Chat:   u1,
	}
	m2 := &telegram.Message{
		Id:     2,
		Text:   "test",
		Sender: u2,
		Chat:   u2,
	}
	m.feed(m1)
	m.feed(m2)

	actual := m.Begin()
	if actual != m1 {
		t.Errorf("Got different message in test 2 msg#1.")
	}
	actual = m.Begin()
	if actual != m2 {
		t.Errorf("Got different message in test 2 msg#2.")
	}
}

func TestManagerWithOneUser(t *testing.T) {
	u1 := makeTestUser("user1")

	m := newManager(bySender, 2)
	m1 := &telegram.Message{
		Id:     1,
		Text:   "test",
		Sender: u1,
		Chat:   u1,
	}
	m2 := &telegram.Message{
		Id:     2,
		Text:   "test",
		Sender: u1,
		Chat:   u1,
	}
	m.feed(m1)
	m.feed(m2)

	actual := m.Begin()
	if actual != m1 {
		t.Errorf("Got different message in test 2 msg#1.")
	}
	m.Commit(m1)
	actual = m.Begin()
	if actual != m2 {
		t.Errorf("Got different message in test 2 msg#2.")
	}
}
