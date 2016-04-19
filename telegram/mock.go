// This file is part of Botgoram
// Botgoram is free software: see LICENSE.txt for more details.

package telegram

import (
	"io"
	"math/rand"
	"time"
)

// FakeUser creates an user.
func FakeUser(id int64) *User {
	return &User{ID: id}
}

// FakeChat creates an chat.
func FakeChat(id int64, chatType ChatType) *Chat {
	return &Chat{User: FakeUser(id), Type: chatType}
}

// MockChat converts a Recipient to Chat
func MockChat(r Recipient) *Chat {
	if chat, ok := r.AsChat(); ok {
		return chat
	}

	u, _ := r.AsUser()
	return FakeChat(u.ID, TYPEGROUP)
}

// FakeAPI implements API interface and does exactly nothing.
// You can pass a channel to provide custom message data, which will be used in GetUpdates method.
// Most identifiers are random-generated, use with cares.
type FakeAPI struct {
	BotUser     *User
	MessagePipe chan *Message
	id          int
}

// Me returns the bot user you set in FakeAPI
func (f *FakeAPI) Me() (*User, error) {
	return f.BotUser, nil
}

// SendMessage returns a new text message as if you sent the request to server
func (f *FakeAPI) SendMessage(victim Recipient, text string, opt *Options) (*Message, error) {
	return &Message{
		ID:     rand.Int(),
		Sender: f.BotUser,
		Chat:   MockChat(victim),
		Text:   text,
	}, nil
}

// ForwardMessage returns a message as if you sent the request to server
func (f *FakeAPI) ForwardMessage(victim, from Recipient, messageID int) (*Message, error) {
	return &Message{
		ID:     messageID,
		Sender: f.BotUser,
		Chat:   MockChat(victim),
		Forward: &Forward{
			From:      f.BotUser,
			Timestamp: time.Now().Unix(),
		},
	}, nil
}

// SendPhoto returns a message as if you sent the request to server
func (f *FakeAPI) SendPhoto(victim Recipient, file *File, caption string, opt *Options) (*Message, error) {
	return &Message{
		ID:     rand.Int(),
		Sender: f.BotUser,
		Chat:   MockChat(victim),
		Photo: []*PhotoSize{
			&PhotoSize{File: file},
		},
		Caption: caption,
	}, nil
}

// SendAudio returns a message as if you sent the request to server
func (f *FakeAPI) SendAudio(victim Recipient, file *File, duration int, performer, title string, opt *Options) (*Message, error) {
	return &Message{
		ID:     rand.Int(),
		Sender: f.BotUser,
		Chat:   MockChat(victim),
		Audio: &Audio{
			File:      file,
			Duration:  duration,
			Performer: performer,
			Title:     title,
		},
	}, nil
}

// SendDocument returns a message as if you sent the request to server
func (f *FakeAPI) SendDocument(victim Recipient, file *File, opt *Options) (*Message, error) {
	return &Message{
		ID:     rand.Int(),
		Sender: f.BotUser,
		Chat:   MockChat(victim),
		Document: &Document{
			File: file,
		},
	}, nil
}

// SendSticker returns a message as if you sent the request to server
func (f *FakeAPI) SendSticker(victim Recipient, file *File, opt *Options) (*Message, error) {
	return &Message{
		ID:     rand.Int(),
		Sender: f.BotUser,
		Chat:   MockChat(victim),
		Sticker: &Sticker{
			File: file,
		},
	}, nil
}

// SendVideo returns a message as if you sent the request to server
func (f *FakeAPI) SendVideo(victim Recipient, file *File, duration int, caption string, opt *Options) (*Message, error) {
	return &Message{
		ID:     rand.Int(),
		Sender: f.BotUser,
		Chat:   MockChat(victim),
		Video: &Video{
			File:     file,
			Duration: duration,
		},
		Caption: caption,
	}, nil
}

// SendVoice returns a message as if you sent the request to server
func (f *FakeAPI) SendVoice(victim Recipient, file *File, duration int, opt *Options) (*Message, error) {
	return &Message{
		ID:     rand.Int(),
		Sender: f.BotUser,
		Chat:   MockChat(victim),
		Video: &Video{
			File:     file,
			Duration: duration,
		},
	}, nil
}

// SendLocation returns a message as if you sent the request to server
func (f *FakeAPI) SendLocation(victim Recipient, location *Location, opt *Options) (*Message, error) {
	return &Message{
		ID:       rand.Int(),
		Sender:   f.BotUser,
		Chat:     MockChat(victim),
		Location: location,
	}, nil
}

// SendChatAction does nothing but return a nil
func (f *FakeAPI) SendChatAction(victim Recipient, action ChatAction) error {
	return nil
}

// GetProfilePhotos returns as if user have go profile photo
func (f *FakeAPI) GetProfilePhotos(victim *User, offset, limit int) (*UserProfilePhotos, error) {
	return new(UserProfilePhotos), nil
}

// GetAllProfilePhotos returns as if user have go profile photo
func (f *FakeAPI) GetAllProfilePhotos(victim *User) (*UserProfilePhotos, error) {
	return new(UserProfilePhotos), nil
}

// DownloadFile returns a reader from the file data you pass
func (f *FakeAPI) DownloadFile(file *File) (io.Reader, error) {
	return file.GetReader()
}

// GetUpdates gets the message from channel
func (f *FakeAPI) GetUpdates(offset, limit, timeout int) (u []Update, err error) {
	if offset < 1 {
		offset = 1
	}
	if limit < 1 {
		limit = 1
	}
	u = make([]Update, 0, limit)
	for ; limit > 0; limit-- {
		select {
		case msg := <-f.MessagePipe:
			u = append(u, Update{ID: offset, Message: msg})
			offset++
		default:
			if timeout > 0 {
				time.Sleep(time.Duration(timeout) * time.Second)
				// ignore err, because we'll never have network problem
				uu, _ := f.GetUpdates(offset, limit, 0)
				u = append(u, uu...)
			}
			return
		}
	}
	return
}

// SetWebhook does exactly nothing but returns a nil
func (f *FakeAPI) SetWebhook(hookURL string, cert []byte) error {
	return nil
}

// EditText modifies the message you passed in, but it will NOT parse message entities
func (f *FakeAPI) EditText(victim Recipient, msg *Message, text string, opt *Options) (*Message, error) {
	msg.Text = text
	return msg, nil
}

// EditInlineText does nothing but return nil
func (f *FakeAPI) EditInlineText(victim Recipient, id, text string, opt *Options) error {
	return nil
}

// EditCaption modifies the message you passed in
func (f *FakeAPI) EditCaption(victim Recipient, msg *Message, caption string, markup *ReplyMarkup) (*Message, error) {
	msg.Caption = caption
	return msg, nil
}

// EditInlineCaption does nothing but return nil
func (f *FakeAPI) EditInlineCaption(victim Recipient, id, caption string, markup *ReplyMarkup) error {
	return nil
}

// EditMarkup returns the message
func (f *FakeAPI) EditMarkup(victim Recipient, msg *Message, markup *ReplyMarkup) (*Message, error) {
	return msg, nil
}

// EditInlineMarkup does nothing but return nil
func (f *FakeAPI) EditInlineMarkup(victim Recipient, id string, markup *ReplyMarkup) error {
	return nil
}

// AnswerInlineQuery does exactly nothing but returns a nil
func (f *FakeAPI) AnswerInlineQuery(query *InlineQuery, results []InlineQueryResult, options *InlineQueryOptions) (err error) {
	return nil
}
