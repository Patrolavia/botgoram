// This file is part of Botgoram
// Botgoram is free software: see LICENSE.txt for more details.

package botgoram

import (
	"testing"

	"github.com/Patrolavia/telegram"
)

func TestTypeTransitors(t *testing.T) {
	st := newState("")
	m := func() *telegram.Message {
		return &telegram.Message{
			From: &telegram.Victim{},
			Chat: &telegram.Victim{},
		}
	}
	factory := func(result string) Transitor {
		f := func(msg *telegram.Message, state State) (next string, err error) {
			return result, nil
		}
		return f
	}
	types := []string{
		TextMsg,
		FileMsg,
		AudioMsg,
		PhotoMsg,
		StickerMsg,
		VideoMsg,
		VoiceMsg,
		ContactMsg,
		LocationMsg,
		VenueMsg,
	}
	for _, t := range types {
		st.Register(t, factory(t))
	}
	st.RegisterCommand("Command", factory("Command"))

	// contact
	msg := m()
	msg.Contact = &telegram.Contact{}
	if next, _ := st.test(msg); next != ContactMsg {
		t.Errorf("While testing contact transitor: get next state %s", next)
	}

	// location
	msg = m()
	msg.Location = &telegram.Location{}
	if next, _ := st.test(msg); next != LocationMsg {
		t.Errorf("While testing location transitor: get next state %s", next)
	}

	// sticker
	msg = m()
	msg.Sticker = &telegram.Sticker{}
	if next, _ := st.test(msg); next != StickerMsg {
		t.Errorf("While testing sticker transitor: get next state %s", next)
	}

	// photo
	msg = m()
	msg.Photo = []telegram.PhotoSize{telegram.PhotoSize{}}
	if next, _ := st.test(msg); next != PhotoMsg {
		t.Errorf("While testing photo transitor: get next state %s", next)
	}

	// video
	msg = m()
	msg.Video = &telegram.Video{}
	if next, _ := st.test(msg); next != VideoMsg {
		t.Errorf("While testing video transitor: get next state %s", next)
	}

	// voice
	msg = m()
	msg.Voice = &telegram.Voice{}
	if next, _ := st.test(msg); next != VoiceMsg {
		t.Errorf("While testing voice transitor: get next state %s", next)
	}

	// audio
	msg = m()
	msg.Audio = &telegram.Audio{}
	if next, _ := st.test(msg); next != AudioMsg {
		t.Errorf("While testing audio transitor: get next state %s", next)
	}

	// document
	msg = m()
	msg.Document = &telegram.Document{}
	if next, _ := st.test(msg); next != FileMsg {
		t.Errorf("While testing document transitor: get next state %s", next)
	}

	// text
	msg = m()
	msg.Text = "asd"
	if next, _ := st.test(msg); next != TextMsg {
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
	if next, _ := st.test(msg); next != TextMsg {
		t.Errorf("While testing command transitor: not fallback to text, get next state %s", next)
	}
}
