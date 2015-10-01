package telegram

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

// API represents all Telegram Bot APIs
type API interface {
	Me() (*User, error)
	SendMessage(victim *User, text string, opt *Options) (*Message, error)
	ForwardMessage(victim, from *User, message_id int) (*Message, error)
	SendPhoto(victim *User, file *File, caption string, opt *Options) (*Message, error)
	SendAudio(victim *User, file *File, duration int, performer, title string, opt *Options) (*Message, error)
	SendDocument(victim *User, file *File, opt *Options) (*Message, error)
	SendSticker(victim *User, file *File, opt *Options) (*Message, error)
	SendVideo(victim *User, file *File, duration int, caption string, opt *Options) (*Message, error)
	SendVoice(victim *User, file *File, duration int, opt *Options) (*Message, error)
	SendLocation(victim *User, location *Location, opt *Options) (*Message, error)
	SendChatAction(victim *User, action ChatAction) error
	GetProfilePhotos(victim *User, offset, limit int) (*UserProfilePhotos, error)
	GetAllProfilePhotos(victim *User) (*UserProfilePhotos, error)
	DownloadFile(file *File) (io.Reader, error)
	GetUpdates(offset, limit, timeout int) ([]Update, error)
	SetWebhook(hook_url string, cert []byte) error
}

type api struct {
	token  string
	client *http.Client
}

// New creates an API instance using default http client
func New(token string) API {
	return &api{token, http.DefaultClient}
}

// NewWithClient creates an API instance using specified http client, useful if your application is running in restricted environment like Google App Engine.
func NewWithClient(token string, client *http.Client) API {
	return &api{token, client}
}

func (a *api) uri(method string) string {
	return fmt.Sprintf("https://api.telegram.org/bot%s/%s", a.token, method)
}

func (a *api) sendCommand(method string, params url.Values) (ret []byte, err error) {
	url := a.uri(method)
	encoding := "application/x-www-form-urlencoded"

	resp, err := a.client.Post(url, encoding, strings.NewReader(params.Encode()))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	ret, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return
}

func (a api) sendFile(method, name string, file *File, params url.Values) (ret []byte, err error) {
	if file.Id != "" {
		params.Set(name, file.Id)
		return a.sendCommand(method, params)
	}
	s, err := file.GetReader()
	if err != nil {
		return
	}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(name, filepath.Base(file.Filename))
	if err != nil {
		return
	}

	if _, err = io.Copy(part, s); err != nil {
		return
	}

	for field, values := range params {
		if len(values) > 0 {
			writer.WriteField(field, values[0])
		}
	}

	if err = writer.Close(); err != nil {
		return
	}

	url := a.uri(method)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return
	}

	req.Header.Add("Content-Type", writer.FormDataContentType())

	resp, err := a.client.Do(req)
	if err != nil {
		return
	}

	ret, err = ioutil.ReadAll(resp.Body)
	return
}
