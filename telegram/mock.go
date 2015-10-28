package telegram

import (
	"io"
	"math/rand"
	"time"
)

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
		Chat:   victim.ToChat(),
		Text:   text,
	}, nil
}

// ForwardMessage returns a message as if you sent the request to server
func (f *FakeAPI) ForwardMessage(victim, from Recipient, messageID int) (*Message, error) {
	return &Message{
		ID:     messageID,
		Sender: f.BotUser,
		Chat:   victim.ToChat(),
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
		Chat:   victim.ToChat(),
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
		Chat:   victim.ToChat(),
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
		Chat:   victim.ToChat(),
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
		Chat:   victim.ToChat(),
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
		Chat:   victim.ToChat(),
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
		Chat:   victim.ToChat(),
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
		Chat:     victim.ToChat(),
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
			u = append(u, Update{offset, msg})
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
