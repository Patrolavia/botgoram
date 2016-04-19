// This file is part of Botgoram
// Botgoram is free software: see LICENSE.txt for more details.

package telegram

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strconv"
)

// abstract the optional parameters
type options interface {
	values() url.Values
}

func optconv(o options) (ret url.Values) {
	ret = url.Values{}
	if o != nil {
		ret = o.values()
	}

	return ret
}

// ReturnValue represents returned values from API method calls
type ReturnValue interface {
	OK() bool
	Err() error
}

type abstractReturnValue struct {
	Ok          bool   `json:"ok"`
	Description string `json:"description,omitempty"`
	ErrorCode   int    `json:"error_code,omitemprt"`
}

func (v *abstractReturnValue) OK() bool {
	return v.Ok
}

func (v *abstractReturnValue) Err() (ret error) {
	if !v.Ok {
		ret = fmt.Errorf("%s (%d)", v.Description, v.ErrorCode)
	}
	return
}

func parseReturnValue(data []byte, val ReturnValue) (err error) {
	if err = json.Unmarshal(data, val); err == nil {
		err = val.Err()
	}
	return
}

// BoolReturnValue represents method return values of bool
type BoolReturnValue struct {
	abstractReturnValue
	Result bool `json:"result,omitempty"`
}

// Parse parses raw data
func (v *BoolReturnValue) Parse(data []byte) (ret bool, err error) {
	err = parseReturnValue(data, v)
	return v.Result, err
}

// MessageReturnValue is structure of returned data from api methods if it returns message
type MessageReturnValue struct {
	abstractReturnValue
	Result *Message `json:"result,omitempty"`
}

// Parse parses raw data
func (v *MessageReturnValue) Parse(data []byte) (ret *Message, err error) {
	err = parseReturnValue(data, v)
	return v.Result, err
}

// UserReturnValue is structure of returned data from api methods if it returns user
type UserReturnValue struct {
	abstractReturnValue
	Result *User `json:"result,omitempty"`
}

// Parse parses raw data
func (v *UserReturnValue) Parse(data []byte) (ret *User, err error) {
	err = parseReturnValue(data, v)
	return v.Result, err
}

// UserProfilePhotosReturnValue is structure of returned data from api methods if it returns user profile photos
type UserProfilePhotosReturnValue struct {
	abstractReturnValue
	Result *UserProfilePhotos `json:"result,omitempty"`
}

// Parse parses raw data
func (v *UserProfilePhotosReturnValue) Parse(data []byte) (ret *UserProfilePhotos, err error) {
	err = parseReturnValue(data, v)
	return v.Result, err
}

func (a *api) Me() (u *User, err error) {
	params := url.Values{}
	data, err := a.sendCommand("getMe", params)
	if err != nil {
		return
	}

	result := &UserReturnValue{}
	return result.Parse(data)
}

func (a *api) SendMessage(victim Recipient, text string, opt *Options) (m *Message, err error) {
	params := optconv(opt)
	params.Set("chat_id", victim.Identifier())
	params.Set("text", text)
	data, err := a.sendCommand("sendMessage", params)
	if err != nil {
		return
	}

	result := &MessageReturnValue{}
	return result.Parse(data)
}

func (a *api) ForwardMessage(victim, from Recipient, messageID int) (m *Message, err error) {
	params := url.Values{}
	params.Set("chat_id", victim.Identifier())
	params.Set("from_chat_id", from.Identifier())
	params.Set("message_id", strconv.Itoa(messageID))
	data, err := a.sendCommand("forwardMessage", params)
	if err != nil {
		return
	}

	result := &MessageReturnValue{}
	return result.Parse(data)
}

func (a *api) SendPhoto(victim Recipient, file *File, caption string, opt *Options) (m *Message, err error) {
	params := optconv(opt)
	params.Set("chat_id", victim.Identifier())
	params.Set("caption", caption)
	data, err := a.sendFile("sendPhoto", "photo", file, params)
	if err != nil {
		return
	}

	result := &MessageReturnValue{}
	return result.Parse(data)
}

func (a *api) SendAudio(victim Recipient, file *File, duration int,
	performer, title string, opt *Options) (m *Message, err error) {

	params := optconv(opt)
	params.Set("chat_id", victim.Identifier())
	optInt(params, "duration", duration)
	optStr(params, "performer", performer)
	optStr(params, "title", title)

	data, err := a.sendFile("sendAudio", "audio", file, params)
	if err != nil {
		return
	}

	result := &MessageReturnValue{}
	return result.Parse(data)
}

func (a *api) SendDocument(victim Recipient, file *File, opt *Options) (m *Message, err error) {
	params := optconv(opt)

	params.Set("chat_id", victim.Identifier())

	data, err := a.sendFile("sendDocument", "document", file, params)
	if err != nil {
		return
	}

	result := &MessageReturnValue{}
	return result.Parse(data)
}

func (a *api) SendSticker(victim Recipient, file *File, opt *Options) (m *Message, err error) {
	params := optconv(opt)

	params.Set("chat_id", victim.Identifier())

	data, err := a.sendFile("sendSticker", "sticker", file, params)
	if err != nil {
		return
	}

	result := &MessageReturnValue{}
	return result.Parse(data)
}

func (a *api) SendVideo(victim Recipient, file *File,
	duration int, caption string, opt *Options) (m *Message, err error) {

	params := optconv(opt)

	params.Set("chat_id", victim.Identifier())
	optInt(params, "duration", duration)
	optStr(params, "caption", caption)

	data, err := a.sendFile("sendVideo", "video", file, params)
	if err != nil {
		return
	}

	result := &MessageReturnValue{}
	return result.Parse(data)
}

func (a *api) SendVoice(victim Recipient, file *File, duration int, opt *Options) (m *Message, err error) {
	params := optconv(opt)

	params.Set("chat_id", victim.Identifier())

	optInt(params, "duration", duration)

	data, err := a.sendFile("sendVoice", "voice", file, params)
	if err != nil {
		return
	}

	result := &MessageReturnValue{}
	return result.Parse(data)
}

func (a *api) SendLocation(victim Recipient, location *Location, opt *Options) (m *Message, err error) {
	params := optconv(opt)

	params.Set("chat_id", victim.Identifier())
	optFloat(params, "latitude", location.Latitude)   // not optional, just lazy
	optFloat(params, "longitude", location.Longitude) // not optional, just lazy

	data, err := a.sendCommand("sendLocation", params)
	if err != nil {
		return
	}

	result := &MessageReturnValue{}
	return result.Parse(data)
}

func (a *api) SendChatAction(victim Recipient, action ChatAction) (err error) {
	params := url.Values{}
	params.Set("action", string(action))
	_, err = a.sendCommand("sendChatAction", params)
	return
}

func (a *api) GetProfilePhotos(victim *User, offset, limit int) (p *UserProfilePhotos, err error) {
	params := url.Values{}
	params.Set("chat_id", victim.Identifier())
	optInt(params, "offset", offset)
	optInt(params, "limit", limit)
	data, err := a.sendCommand("getUserProfilePhotos", params)
	if err != nil {
		return
	}

	result := &UserProfilePhotosReturnValue{}
	return result.Parse(data)
}

func (a *api) GetAllProfilePhotos(victim *User) (p *UserProfilePhotos, err error) {
	params := url.Values{}
	params.Set("chat_id", victim.Identifier())
	data, err := a.sendCommand("getUserProfilePhotos", params)
	if err != nil {
		return
	}

	result := &UserProfilePhotosReturnValue{}
	return result.Parse(data)
}

func (a *api) DownloadFile(file *File) (r io.Reader, err error) {
	if file.ID == "" {
		return r, errors.New("file_id not specified, not remote file?")
	}
	type fileToken struct {
		ID   string `json:"file_id"`
		Size int    `json:"file_size,omitempty"`
		Path string `json:"file_path"`
	}

	params := url.Values{}
	params.Set("file_id", file.ID)
	data, err := a.sendCommand("getFile", params)
	if err != nil {
		return
	}
	token := fileToken{}
	if err = json.Unmarshal(data, &token); err != nil {
		return
	}
	resp, err := a.client.Get(a.uri("getFile"))
	if err == nil {
		r = resp.Body
	}
	return
}

func (a *api) GetUpdates(offset, limit, timeout int) (u []Update, err error) {
	params := url.Values{}

	optInt(params, "offset", offset)
	optInt(params, "limit", limit)
	optInt(params, "timeout", timeout)

	data, err := a.sendCommand("getUpdates", params)
	if err != nil {
		return
	}
	var res updates
	if err = json.Unmarshal(data, &res); err == nil {
		u = res.Result
	}
	return
}

func (a *api) SetWebhook(hookURL string, cert []byte) (err error) {
	params := url.Values{}
	optStr(params, "url", hookURL)
	if cert != nil {
		f := &File{Filename: "server.cert", Stream: bytes.NewReader(cert)}
		_, err = a.sendFile("setWebhook", "certificate", f, params)
	} else {
		_, err = a.sendCommand("setWebhook", params)
	}
	return
}

// AnswerIQResult is structure of returned data from api method "answerInlineQuery"
type AnswerIQResult struct {
	Ok          bool   `json:"ok"`
	Result      bool   `json:"result,omitempty"`
	ErrorCode   int    `json:"error_code,omitempty"`
	Description string `json:"description,omitempty"`
}

func (a *api) doEdit(params url.Values, method string) (ret *Message, err error) {
	data, err := a.sendCommand(method, params)
	if err != nil {
		return
	}

	var res MessageReturnValue
	return res.Parse(data)
}

func (a *api) doInlineEdit(params url.Values, method string) (err error) {
	data, err := a.sendCommand(method, params)
	if err != nil {
		return
	}

	var res BoolReturnValue
	ok, err := res.Parse(data)
	if err != nil {
		return
	}

	if !ok {
		err = fmt.Errorf("Calling to %s failed!", method)
	}
	return
}

func (a *api) EditText(victim Recipient, msg *Message, text string, opt *Options) (ret *Message, err error) {
	params := optconv(opt)

	params.Set("chat_id", victim.Identifier())
	params.Set("message_id", strconv.Itoa(msg.ID))
	params.Set("text", text)

	return a.doEdit(params, "editMessageText")
}

func (a *api) EditInlineText(victim Recipient, id, text string, opt *Options) error {
	params := optconv(opt)

	params.Set("chat_id", victim.Identifier())
	params.Set("inline_message_id", id)
	params.Set("text", text)

	return a.doInlineEdit(params, "editMessageText")
}

func (a *api) EditCaption(victim Recipient, msg *Message, caption string, markup *ReplyMarkup) (*Message, error) {
	params := optconv(markup)

	params.Set("chat_id", victim.Identifier())
	params.Set("message_id", strconv.Itoa(msg.ID))
	params.Set("caption", caption)

	return a.doEdit(params, "editMessageCaption")
}

func (a *api) EditInlineCaption(victim Recipient, id, caption string, markup *ReplyMarkup) error {
	params := optconv(markup)

	params.Set("chat_id", victim.Identifier())
	params.Set("inline_message_id", id)
	params.Set("caption", caption)

	return a.doInlineEdit(params, "editMessageCaption")
}

func (a *api) EditMarkup(victim Recipient, msg *Message, markup *ReplyMarkup) (*Message, error) {
	params := optconv(markup)

	params.Set("chat_id", victim.Identifier())
	params.Set("message_id", strconv.Itoa(msg.ID))

	return a.doEdit(params, "editMessageReplyMarkup")
}

func (a *api) EditInlineMarkup(victim Recipient, id string, markup *ReplyMarkup) error {
	params := optconv(markup)

	params.Set("chat_id", victim.Identifier())
	params.Set("inline_message_id", id)

	return a.doInlineEdit(params, "editMessageReplyMarkup")
}

func (a *api) AnswerCallbackQuery(id string, text string, alert bool) (err error) {
	params := url.Values{}

	params.Set("callback_query_id", id)
	optStr(params, "text", text)
	optBool(params, "show_alert", alert)

	data, err := a.sendCommand("answerCallbackQuery", params)
	if err != nil {
		return
	}

	var result BoolReturnValue
	_, err = result.Parse(data)
	return
}

func (a *api) AnswerInlineQuery(query *InlineQuery, results []InlineQueryResult, opts *InlineQueryOptions) (err error) {
	for _, r := range results {
		r.ForceType()
	}

	params := optconv(opts)
	params.Set("inline_query_id", query.ID)
	optJSON(params, "results", results)

	data, err := a.sendCommand("answerInlineQuery", params)
	if err != nil {
		return
	}

	var result BoolReturnValue
	_, err = result.Parse(data)
	return
}
