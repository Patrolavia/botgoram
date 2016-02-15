// This file is part of Botgoram
// Botgoram is free software: see LICENSE.txt for more details.

package telegram

import (
	"encoding/json"
	"net/url"
	"strconv"
)

// ChatAction controls how you tell the user that something is taking some time on the bot's side.
type ChatAction string

// predefined chat actions
const (
	Typing         ChatAction = "typing"
	UploadPhoto    ChatAction = "upload_photo"
	RecordVideo    ChatAction = "record_video"
	UploadVideo    ChatAction = "upload_video"
	RecordAudio    ChatAction = "record_audio"
	UploadAudio    ChatAction = "upload_audio"
	UploadDocument ChatAction = "upload_document"
	FindLocation   ChatAction = "find_location"
)

// ParseMode controls how Telegram client render your message.
type ParseMode string

// predefined parse modes
const (
	AsText     ParseMode = ""
	AsMarkdown ParseMode = "Markdown"
)

// ReplyMarkup enables ui features on client when message is received.
type ReplyMarkup struct {
	Keyboard  [][]string `json:"keyboard,omitempty"`        // Array of button rows, each represented by an Array of Strings
	Resize    bool       `json:"resize_keyboard,omitempty"` // Optional. Requests clients to resize the keyboard vertically for optimal fit
	OneTime   bool       `json:"one_time_keyboard"`         // Optional. Requests clients to hide the keyboard as soon as it's been used. Defaults to false.
	Selective bool       `json:"selective,omitempty"`       // Optional. Use this parameter if you want to custom keyboard for specific users only.
	Hide      bool       `json:"hide_keyboard,omitempty"`   // Requests clients to hide the custom keyboard
	Reply     bool       `json:"force_reply"`               // Shows reply interface to the user, as if they manually selected the bot‘s message and tapped ’Reply'
}

// Options represents optional features when sending message.
type Options struct {
	ParseMode         ParseMode // Optional
	DisableWebPreview bool      // Optional. Disables link previews for links in this message
	ReplyTo           int       // Optional. If the message is a reply, ID of the original message
	*ReplyMarkup
}

func (opt *Options) encode() (ret url.Values, err error) {
	ret = url.Values{}

	if opt.ParseMode != AsText {
		ret.Set("parse_mode", string(opt.ParseMode))
	}
	if opt.DisableWebPreview {
		ret.Set("disable_web_page_preview", "true")
	}
	if opt.ReplyTo != 0 {
		ret.Set("reply_to_message_id", itoa(opt.ReplyTo))
	}
	if opt.ReplyMarkup != nil {
		data, err := json.Marshal(opt.ReplyMarkup)
		if err != nil {
			return ret, err
		}
		ret.Set("reply_markup", string(data))
	}
	return
}

func itoa(i int) string {
	return strconv.FormatInt(int64(i), 10)
}

func ftoa(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 32)
}
