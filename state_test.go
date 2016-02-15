// This file is part of Botgoram
// Botgoram is free software: see LICENSE.txt for more details.

package botgoram

import (
	"testing"

	"github.com/Patrolavia/botgoram/telegram"
)

func TestTypeTransitors(t *testing.T) {
	st := newState("")
	m := func() *telegram.Message {
		return &telegram.Message{
			Sender: &telegram.User{},
			Chat:   telegram.MockChat(&telegram.User{}),
		}
	}
	factory := func(result string) Transitor {
		f := func(msg *telegram.Message, state State) (next string, err error) {
			return string(result), nil
		}
		return f
	}
	types := []telegram.MessageType{
		telegram.CONTACT,
		telegram.LOCATION,
		telegram.STICKER,
		telegram.PHOTO,
		telegram.VIDEO,
		telegram.AUDIO,
		telegram.VOICE,
		telegram.DOCUMENT,
		telegram.TEXT,
		telegram.STATUS,
	}
	for _, t := range types {
		st.Register(t, factory(t.String()))
	}
	st.RegisterCommand("Command", factory("Command"))

	// contact
	msg := m()
	msg.Contact = &telegram.Contact{}
	if next, _ := st.test(msg); next != telegram.CONTACT.String() {
		t.Errorf("While testing contact transitor: get next state %s", next)
	}

	// location
	msg = m()
	msg.Location = &telegram.Location{}
	if next, _ := st.test(msg); next != telegram.LOCATION.String() {
		t.Errorf("While testing location transitor: get next state %s", next)
	}

	// sticker
	msg = m()
	msg.Sticker = &telegram.Sticker{}
	if next, _ := st.test(msg); next != telegram.STICKER.String() {
		t.Errorf("While testing sticker transitor: get next state %s", next)
	}

	// photo
	msg = m()
	msg.Photo = []*telegram.PhotoSize{&telegram.PhotoSize{}}
	if next, _ := st.test(msg); next != telegram.PHOTO.String() {
		t.Errorf("While testing photo transitor: get next state %s", next)
	}

	// video
	msg = m()
	msg.Video = &telegram.Video{}
	if next, _ := st.test(msg); next != telegram.VIDEO.String() {
		t.Errorf("While testing video transitor: get next state %s", next)
	}

	// voice
	msg = m()
	msg.Voice = &telegram.Voice{}
	if next, _ := st.test(msg); next != telegram.VOICE.String() {
		t.Errorf("While testing voice transitor: get next state %s", next)
	}

	// audio
	msg = m()
	msg.Audio = &telegram.Audio{}
	if next, _ := st.test(msg); next != telegram.AUDIO.String() {
		t.Errorf("While testing audio transitor: get next state %s", next)
	}

	// document
	msg = m()
	msg.Document = &telegram.Document{}
	if next, _ := st.test(msg); next != telegram.DOCUMENT.String() {
		t.Errorf("While testing document transitor: get next state %s", next)
	}

	// text
	msg = m()
	msg.Text = "asd"
	if next, _ := st.test(msg); next != telegram.TEXT.String() {
		t.Errorf("While testing text transitor: get next state %s", next)
	}

	// command
	msg = m()
	msg.Text = "Command"
	if next, _ := st.test(msg); next != "Command" {
		t.Errorf("While testing command transitor: get next state %s", next)
	}

	// command mismatch, fallback to text transitor
	msg = m()
	msg.Text = "Command2"
	if next, _ := st.test(msg); next != telegram.TEXT.String() {
		t.Errorf("While testing command transitor: not fallback to text, get next state %s", next)
	}
}
