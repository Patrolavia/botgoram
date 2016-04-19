// This file is part of Botgoram
// Botgoram is free software: see LICENSE.txt for more details.

package telegram

import "net/url"

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
	AsHTML     ParseMode = "HTML"
)

// KeyboardButton represents a button of custom keyboard
type KeyboardButton struct {
	Text            string `json:"text"`                       // Text of the button. If none of the optional fields are used, it will be sent to the bot as a message when the button is pressed
	RequestContact  bool   `json:"request_contact,omitempty"`  // Optional. If True, the user's phone number will be sent as a contact when the button is pressed. Available in private chats only
	RequestLocation bool   `json:"request_location,omitempty"` // Optional. If True, the user's current location will be sent when the button is pressed. Available in private chats only
}

// InlineKeyboardButton represents a button of custom inline keyboard
type InlineKeyboardButton struct {
	Text string `json:"text"`                    // Label text on the button
	URL  string `json:"url,omitempty"`           // Optional. HTTP url to be opened when button is pressed
	Data string `json:"callback_data,omitempty"` // Optional. Data to be sent in a callback query to the bot when button is pressed
	// SwitchChat is optional. If set, pressing the button will prompt the user to select one of their chats,
	// open that chat and insert the bot‘s username and the specified inline query in the input field.
	// Can be empty, in which case just the bot’s username will be inserted.
	//
	// Note: This offers an easy way for users to start using your bot in inline mode when
	// they are currently in a private chat with it.
	// Especially useful when combined with switch_pm… actions – in this case the user will be
	// automatically returned to the chat they switched from, skipping the chat selection screen.
	SwitchChat string `json:"switch_inline_query,omitempty"`
}

// ReplyMarkup enables ui features on client when message is received.
type ReplyMarkup struct {
	Keyboard       [][]KeyboardButton       `json:"keyboard,omitempty"`          // Array of button rows, each represented by an Array of Strings
	Resize         bool                     `json:"resize_keyboard,omitempty"`   // Optional. Requests clients to resize the keyboard vertically for optimal fit
	Once           bool                     `json:"one_time_keyboard,omitempty"` // Optional. Requests clients to hide the keyboard as soon as it's been used. Defaults to false.
	Selective      bool                     `json:"selective,omitempty"`         // Optional. Use this parameter if you want to custom keyboard for specific users only.
	Hide           bool                     `json:"hide_keyboard,omitempty"`     // Requests clients to hide the custom keyboard
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard,omitempty"`   // Array of button rows, each represented by an Array of InlineKeyboardButton objects
	ForceReply     bool                     `json:"force_reply,omitempty"`       // Shows reply interface to the user, as if they manually selected the bot‘s message and tapped ’Reply'
}

// Values converts these optional parameters to url.Values. This implements internal option interface
func (r *ReplyMarkup) values() url.Values {
	ret := url.Values{}

	optJSON(ret, "keyboard", r.Keyboard)
	optBool(ret, "resize_keyboard", r.Resize)
	optBool(ret, "one_time_keyboard", r.Once)
	optBool(ret, "selective", r.Selective)
	optBool(ret, "hide_keyboard", r.Hide)
	optJSON(ret, "inline_keyboard", r.InlineKeyboard)

	return ret
}

// Options represents optional features when sending message.
type Options struct {
	ParseMode ParseMode // Optional
	NoPreview bool      // Optional. Disables link previews for links in this message
	ReplyTo   int       // Optional. If the message is a reply, ID of the original message
	*ReplyMarkup
	// Sends the message silently.
	// iOS users will not receive a notification, Android users will receive a notification with no sound.
	Silent bool
}

// Values converts these optional parameters to url.Values. This implements internal option interface
func (o *Options) values() url.Values {
	ret := o.ReplyMarkup.values()

	optStr(ret, "parse_mode", string(o.ParseMode))
	optBool(ret, "disable_web_page_preview", o.NoPreview)
	optBool(ret, "disable_notification", o.Silent)
	optInt(ret, "reply_to_message_id", o.ReplyTo)

	return ret
}
