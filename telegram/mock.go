package telegram

import (
	"io"
	"math/rand"
	"time"
)

// FakeAPI is mock object implements API interface and does exactly nothing.

type FakeAPI struct {
	BotUser *User
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

func (f *FakeAPI) GetUpdates(offset, limit, timeout int) ([]Update, error) {
	return []Update{}, nil
}

func (f *FakeAPI) SetWebhook(hook_url string, cert []byte) error {
	return nil
}
