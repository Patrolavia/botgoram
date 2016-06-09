// This file is part of Botgoram
// Botgoram is free software: see LICENSE.txt for more details.

package botgoram

import (
	"math/rand"
	"testing"

	"github.com/Patrolavia/telegram"
)

func makeTestUser(name string) *telegram.Victim {
	return &telegram.Victim{
		ID:        int64(rand.Int()),
		FirstName: name,
	}
}

// test with 2 worker, 2 sender, no overflow, no block
func TestManagerWithTwoUser(t *testing.T) {
	u1 := makeTestUser("user1")
	u2 := makeTestUser("user2")
	ch := make(chan *telegram.Message)

	m := newManager(bySender, 2, ch)
	m1 := &telegram.Message{
		ID:   1,
		Text: "test",
		From: u1,
		Chat: u1,
	}
	m2 := &telegram.Message{
		ID:   2,
		Text: "test",
		From: u2,
		Chat: u2,
	}
	go m.feed([]*telegram.Message{m1, m2})

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
	ch := make(chan *telegram.Message)

	m := newManager(bySender, 2, ch)
	m1 := &telegram.Message{
		ID:   1,
		Text: "test",
		From: u1,
		Chat: u1,
	}
	m2 := &telegram.Message{
		ID:   2,
		Text: "test",
		From: u1,
		Chat: u1,
	}
	go m.feed([]*telegram.Message{m1, m2})

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
