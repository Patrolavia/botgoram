// This file is part of Botgoram
// Botgoram is free software: see LICENSE.txt for more details.

// Package fake is mocked telegram API
package fake

import (
	"io"
	"math/rand"
	"time"

	"github.com/Patrolavia/botgoram/telegram"
)

// User creates an user.
func User(id int64) *telegram.User {
	return &telegram.User{ID: id}
}

// Chat creates an chat.
func Chat(id int64, chatType telegram.ChatType) *telegram.Chat {
	return &telegram.Chat{User: User(id), Type: chatType}
}

// MockChat converts a telegram.Recipient to Chat
func MockChat(r telegram.Recipient) *telegram.Chat {
	if chat, ok := r.AsChat(); ok {
		return chat
	}

	u, _ := r.AsUser()
	return Chat(u.ID, telegram.TYPEGROUP)
}

// API implements API interface and does exactly nothing.
// You can pass a channel to provide custom message data, which will be used in GetUpdates method.
// Most identifiers are random-generated, use with cares.
type API struct {
	BotUser     *telegram.User
	MessagePipe chan *telegram.Message
	id          int
}

// Me returns the bot user you set in API
func (f *API) Me() (*telegram.User, error) {
	return f.BotUser, nil
}

// SendMessage returns a new text message as if you sent the request to server
func (f *API) SendMessage(victim telegram.Recipient, text string, opt *telegram.Options) (*telegram.Message, error) {
	return &telegram.Message{
		ID:     rand.Int(),
		Sender: f.BotUser,
		Chat:   MockChat(victim),
		Text:   text,
	}, nil
}

// ForwardMessage returns a message as if you sent the request to server
func (f *API) ForwardMessage(victim, from telegram.Recipient, messageID int) (*telegram.Message, error) {
	return &telegram.Message{
		ID:     messageID,
		Sender: f.BotUser,
		Chat:   MockChat(victim),
		Forward: &telegram.Forward{
			From:      f.BotUser,
			Timestamp: time.Now().Unix(),
		},
	}, nil
}

// SendPhoto returns a message as if you sent the request to server
func (f *API) SendPhoto(victim telegram.Recipient, file *telegram.File, caption string, opt *telegram.Options) (*telegram.Message, error) {
	return &telegram.Message{
		ID:     rand.Int(),
		Sender: f.BotUser,
		Chat:   MockChat(victim),
		Photo: []*telegram.PhotoSize{
			&telegram.PhotoSize{File: file},
		},
		Caption: caption,
	}, nil
}

// SendAudio returns a message as if you sent the request to server
func (f *API) SendAudio(victim telegram.Recipient, file *telegram.File, duration int, performer, title string, opt *telegram.Options) (*telegram.Message, error) {
	return &telegram.Message{
		ID:     rand.Int(),
		Sender: f.BotUser,
		Chat:   MockChat(victim),
		Audio: &telegram.Audio{
			File:      file,
			Duration:  duration,
			Performer: performer,
			Title:     title,
		},
	}, nil
}

// SendDocument returns a message as if you sent the request to server
func (f *API) SendDocument(victim telegram.Recipient, file *telegram.File, opt *telegram.Options) (*telegram.Message, error) {
	return &telegram.Message{
		ID:     rand.Int(),
		Sender: f.BotUser,
		Chat:   MockChat(victim),
		Document: &telegram.Document{
			File: file,
		},
	}, nil
}

// SendSticker returns a message as if you sent the request to server
func (f *API) SendSticker(victim telegram.Recipient, file *telegram.File, opt *telegram.Options) (*telegram.Message, error) {
	return &telegram.Message{
		ID:     rand.Int(),
		Sender: f.BotUser,
		Chat:   MockChat(victim),
		Sticker: &telegram.Sticker{
			File: file,
		},
	}, nil
}

// SendVideo returns a message as if you sent the request to server
func (f *API) SendVideo(victim telegram.Recipient, file *telegram.File, duration int, caption string, opt *telegram.Options) (*telegram.Message, error) {
	return &telegram.Message{
		ID:     rand.Int(),
		Sender: f.BotUser,
		Chat:   MockChat(victim),
		Video: &telegram.Video{
			File:     file,
			Duration: duration,
		},
		Caption: caption,
	}, nil
}

// SendVoice returns a message as if you sent the request to server
func (f *API) SendVoice(victim telegram.Recipient, file *telegram.File, duration int, opt *telegram.Options) (*telegram.Message, error) {
	return &telegram.Message{
		ID:     rand.Int(),
		Sender: f.BotUser,
		Chat:   MockChat(victim),
		Video: &telegram.Video{
			File:     file,
			Duration: duration,
		},
	}, nil
}

// SendLocation returns a message as if you sent the request to server
func (f *API) SendLocation(victim telegram.Recipient, location *telegram.Location, opt *telegram.Options) (*telegram.Message, error) {
	return &telegram.Message{
		ID:       rand.Int(),
		Sender:   f.BotUser,
		Chat:     MockChat(victim),
		Location: location,
	}, nil
}

// Sendtelegram.ChatAction does nothing but return a nil
func (f *API) SendChatAction(victim telegram.Recipient, action telegram.ChatAction) error {
	return nil
}

// GetProfilePhotos returns as if user have go profile photo
func (f *API) GetProfilePhotos(victim *telegram.User, offset, limit int) (*telegram.UserProfilePhotos, error) {
	return new(telegram.UserProfilePhotos), nil
}

// GetAllProfilePhotos returns as if user have go profile photo
func (f *API) GetAllProfilePhotos(victim *telegram.User) (*telegram.UserProfilePhotos, error) {
	return new(telegram.UserProfilePhotos), nil
}

// DownloadFile returns a reader from the file data you pass
func (f *API) DownloadFile(file *telegram.File) (io.Reader, error) {
	return file.GetReader()
}

// GetUpdates gets the message from channel
func (f *API) GetUpdates(offset, limit, timeout int) (u []telegram.Update, err error) {
	if offset < 1 {
		offset = 1
	}
	if limit < 1 {
		limit = 1
	}
	u = make([]telegram.Update, 0, limit)
	for ; limit > 0; limit-- {
		select {
		case msg := <-f.MessagePipe:
			u = append(u, telegram.Update{ID: offset, Message: msg})
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
func (f *API) SetWebhook(hookURL string, cert []byte) error {
	return nil
}

// EditText modifies the message you passed in, but it will NOT parse message entities
func (f *API) EditText(victim telegram.Recipient, msg *telegram.Message, text string, opt *telegram.Options) (*telegram.Message, error) {
	msg.Text = text
	return msg, nil
}

// EditInlineText does nothing but return nil
func (f *API) EditInlineText(victim telegram.Recipient, id, text string, opt *telegram.Options) error {
	return nil
}

// EditCaption modifies the message you passed in
func (f *API) EditCaption(victim telegram.Recipient, msg *telegram.Message, caption string, markup *telegram.ReplyMarkup) (*telegram.Message, error) {
	msg.Caption = caption
	return msg, nil
}

// EditInlineCaption does nothing but return nil
func (f *API) EditInlineCaption(victim telegram.Recipient, id, caption string, markup *telegram.ReplyMarkup) error {
	return nil
}

// EditMarkup returns the message
func (f *API) EditMarkup(victim telegram.Recipient, msg *telegram.Message, markup *telegram.ReplyMarkup) (*telegram.Message, error) {
	return msg, nil
}

// EditInlineMarkup does nothing but return nil
func (f *API) EditInlineMarkup(victim telegram.Recipient, id string, markup *telegram.ReplyMarkup) error {
	return nil
}

// AnswerInlineQuery does exactly nothing but returns a nil
func (f *API) AnswerInlineQuery(query *telegram.InlineQuery, results []telegram.InlineQueryResult, options *telegram.InlineQueryOptions) (err error) {
	return nil
}
