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

func (f *FakeAPI) Me() (*User, error) {
	return f.BotUser, nil
}

func (f *FakeAPI) SendMessage(victim *User, text string, opt *Options) (*Message, error) {
	return &Message{
		Id:     rand.Int(),
		Sender: f.BotUser,
		Chat:   victim,
		Text:   text,
	}, nil
}

func (f *FakeAPI) ForwardMessage(victim, from *User, message_id int) (*Message, error) {
	return &Message{
		Id:     message_id,
		Sender: f.BotUser,
		Chat:   victim,
		Forward: &Forward{
			From:      f.BotUser,
			Timestamp: time.Now().Unix(),
		},
	}, nil
}

func (f *FakeAPI) SendPhoto(victim *User, file *File, caption string, opt *Options) (*Message, error) {
	return &Message{
		Id:     rand.Int(),
		Sender: f.BotUser,
		Chat:   victim,
		Photo: []*PhotoSize{
			&PhotoSize{File: file},
		},
		Caption: caption,
	}, nil
}

func (f *FakeAPI) SendAudio(victim *User, file *File, duration int, performer, title string, opt *Options) (*Message, error) {
	return &Message{
		Id:     rand.Int(),
		Sender: f.BotUser,
		Chat:   victim,
		Audio: &Audio{
			File:      file,
			Duration:  duration,
			Performer: performer,
			Title:     title,
		},
	}, nil
}

func (f *FakeAPI) SendDocument(victim *User, file *File, opt *Options) (*Message, error) {
	return &Message{
		Id:     rand.Int(),
		Sender: f.BotUser,
		Chat:   victim,
		Document: &Document{
			File: file,
		},
	}, nil
}

func (f *FakeAPI) SendSticker(victim *User, file *File, opt *Options) (*Message, error) {
	return &Message{
		Id:     rand.Int(),
		Sender: f.BotUser,
		Chat:   victim,
		Sticker: &Sticker{
			File: file,
		},
	}, nil
}

func (f *FakeAPI) SendVideo(victim *User, file *File, duration int, caption string, opt *Options) (*Message, error) {
	return &Message{
		Id:     rand.Int(),
		Sender: f.BotUser,
		Chat:   victim,
		Video: &Video{
			File:     file,
			Duration: duration,
		},
		Caption: caption,
	}, nil
}

func (f *FakeAPI) SendVoice(victim *User, file *File, duration int, opt *Options) (*Message, error) {
	return &Message{
		Id:     rand.Int(),
		Sender: f.BotUser,
		Chat:   victim,
		Video: &Video{
			File:     file,
			Duration: duration,
		},
	}, nil
}

func (f *FakeAPI) SendLocation(victim *User, location *Location, opt *Options) (*Message, error) {
	return &Message{
		Id:       rand.Int(),
		Sender:   f.BotUser,
		Chat:     victim,
		Location: location,
	}, nil
}

func (f *FakeAPI) SendChatAction(victim *User, action ChatAction) error {
	return nil
}

func (f *FakeAPI) GetProfilePhotos(victim *User, offset, limit int) (*UserProfilePhotos, error) {
	return new(UserProfilePhotos), nil
}

func (f *FakeAPI) GetAllProfilePhotos(victim *User) (*UserProfilePhotos, error) {
	return new(UserProfilePhotos), nil
}

func (f *FakeAPI) DownloadFile(file *File) (io.Reader, error) {
	return file.GetReader()
}

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

func (f *FakeAPI) SetWebhook(hook_url string, cert []byte) error {
	return nil
}
