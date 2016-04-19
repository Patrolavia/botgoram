// This file is part of Botgoram
// Botgoram is free software: see LICENSE.txt for more details.

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
	// Main API methods. See https://core.telegram.org/bots/api#available-methods
	Me() (*User, error)
	SendMessage(victim Recipient, text string, opt *Options) (*Message, error)
	ForwardMessage(victim, from Recipient, messageID int) (*Message, error)
	SendPhoto(victim Recipient, file *File, caption string, opt *Options) (*Message, error)
	SendAudio(victim Recipient, file *File, duration int, performer, title string, opt *Options) (*Message, error)
	SendDocument(victim Recipient, file *File, opt *Options) (*Message, error)
	SendSticker(victim Recipient, file *File, opt *Options) (*Message, error)
	SendVideo(victim Recipient, file *File, duration int, caption string, opt *Options) (*Message, error)
	SendVoice(victim Recipient, file *File, duration int, opt *Options) (*Message, error)
	SendLocation(victim Recipient, location *Location, opt *Options) (*Message, error)
	SendChatAction(victim Recipient, action ChatAction) error
	GetProfilePhotos(victim *User, offset, limit int) (*UserProfilePhotos, error)
	GetAllProfilePhotos(victim *User) (*UserProfilePhotos, error)
	DownloadFile(file *File) (io.Reader, error)

	// Getting updates. See https://core.telegram.org/bots/api#getting-updates
	GetUpdates(offset, limit, timeout int) ([]Update, error)
	SetWebhook(hookURL string, cert []byte) error

	// API methods to update bot message. See https://core.telegram.org/bots/2-0-intro#updating-messages
	EditText(victim Recipient, msg *Message, text string, opt *Options) (*Message, error)
	EditInlineText(victim Recipient, id, text string, opt *Options) error
	EditCaption(victim Recipient, msg *Message, caption string, markup *ReplyMarkup) (*Message, error)
	EditInlineCaption(victim Recipient, id, caption string, markup *ReplyMarkup) error
	EditMarkup(victim Recipient, msg *Message, markup *ReplyMarkup) (*Message, error)
	EditInlineMarkup(victim Recipient, id string, markup *ReplyMarkup) error

	// AnswerCallbackQuery is used to send answers to callback queries sent from inline keyboards.
	// The answer will be displayed to the user as a notification at the top of the chat screen or as an alert.
	// On success, True is returned.
	AnswerCallbackQuery(id string, text string, alert bool) error

	// AnswerInlineQuery sends answers to an inline query. On success, True is returned.
	// No more than 50 results per query are allowed.
	//
	//   * cacheTime is maximum amount of time in seconds that the result of the inline query may be cached on the server. defaults to 300.
	//   * personal controls whether results may be cached on the server side only for the user that sent the query.
	//     By default, results may be returned to any user who sends the same query
	//   * next is the offset that a client should send in the next query with the same text to receive more results.
	//     Pass an empty string if there are no more results or if you don‘t support pagination. Offset length can’t exceed 64 bytes.
	AnswerInlineQuery(query *InlineQuery, results []InlineQueryResult, options *InlineQueryOptions) (err error)
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
	if file.ID != "" {
		params.Set(name, file.ID)
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
	defer resp.Body.Close()

	ret, err = ioutil.ReadAll(resp.Body)
	return
}
