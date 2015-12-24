package telegram

import (
	"errors"
	"io"
	"os"
	"strconv"
)

// Recipient is abstract parent for User and Chat
type Recipient interface {
	Identifier() string
	// Name returns user's name and username (if has one) if regular user, title if chat group.
	Name() string
}

// User struct represents a Telegram user or chat group.
type User struct {
	ID        int64  `json:"id"`                   // Unique identifier for this user or bot
	FirstName string `json:"first_name,omitempty"` // User‘s or bot’s first name
	LastName  string `json:"last_name,omitempty"`  // Optional. User‘s or bot’s last name
	Username  string `json:"username,omitempty"`   // Optional. User‘s or bot’s username
}

func (u User) Identifier() string {
	return strconv.FormatInt(u.ID, 10)
}

// Name returns user's name and username (if has one) if regular user, title if chat group.
func (u User) Name() string {
	ret := u.FirstName
	if u.LastName != "" {
		ret += " " + u.LastName
	}
	if u.Username != "" {
		ret += "(" + u.Username + ")"
	}
	return ret
}

type ChatType string

const (
	TYPECHAT       ChatType = "private"
	TYPEGROUP      ChatType = "group"
	TYPESUPERGROUP ChatType = "supergroup"
	TYPECHANNEL    ChatType = "channel"
)

// type Chat represents a chat
type Chat struct {
	*User
	Title string   `json:"title,omitempty"` // Group name
	Type  ChatType `json:"type"`            // Type of chat, can be either “private”, “group”, "supergroup" or “channel”
}

func (c Chat) Identifier() string {
	if c.Type == `channel` {
		return "@" + c.User.Username
	}
	return c.User.Identifier()
}

func (c Chat) Name() string {
	return c.Title
}

// File represents a regular file for sending.
// You should avoid re-upload same file by remembering file id.
//
// Sending local file
//
// A File struct with empty ID field is considered as local file, calling API methods with local file will upload it.
// API methods first try read from Stream field, so you can generate something (an image for example) and send it to your user without cache it in disk.
// Then it tries to open file specified in Filename field.
// If neither is possible, it would return in error.
type File struct {
	ID       string    `json:"file_id"`             // Unique identifier for this file
	MimeType string    `json:"mime_type,omitempty"` // Optional. MIME type of the file as defined by sender
	Size     int       `json:"file_size,omitempty"` // Optional. File size
	Filename string    `json:"file_name,omitempty"` // Optional. Local file name
	Stream   io.Reader // Optional.
}

// GetReader returns an io.Reader if it has information about local file.
func (f *File) GetReader() (io.Reader, error) {
	if f.Stream != nil {
		return f.Stream, nil
	}
	if f.Filename == "" {
		return nil, errors.New("telegram: No local file information in this File struct")
	}
	return os.Open(f.Filename)
}

// PhotoSize represents one size of a photo or a file / sticker thumbnail.
type PhotoSize struct {
	*File
	Width  int `json:"width"`  // Photo width
	Height int `json:"height"` // Photo height
}

// Document represents a general file (as opposed to photos, voice messages and audio files).
type Document struct {
	*File
	Thumb *PhotoSize `json:"thumb"` // Optional. Document thumbnail as defined by sender
}

// Audio represents an audio file to be treated as music by the Telegram clients.
type Audio struct {
	*File
	Duration  int    `json:"duration"`            // Duration of the audio in seconds as defined by sender
	Performer string `json:"performer,omitempty"` // Optional. Performer of the audio as defined by sender or by audio tags
	Title     string `json:"title,omitempty"`     // Optional. Title of the audio as defined by sender or by audio tags
}

// Sticker represents a sticker.
type Sticker struct {
	*File
	Width  int        `json:"width"`           // Sticker width
	Height int        `json:"height"`          // Sticker height
	Thumb  *PhotoSize `json:"thumb,omitempty"` // Optional. Sticker thumbnail in .webp or .jpg format
}

// Video represents a video file.
type Video struct {
	*File
	Width    int        `json:"width"`           // Video width as defined by sender
	Height   int        `json:"height"`          // Video height as defined by sender
	Duration int        `json:"duration"`        // Duration of the video in seconds as defined by sender
	Thumb    *PhotoSize `json:"thumb,omitempty"` // Optional. Video thumbnail
}

// Voice represents a voice note.
type Voice struct {
	*File
	Duration int `json:"duration"` // Duration of the audio in seconds as defined by sender
}

// Contact represents a phone contact.
type Contact struct {
	Number    string `json:"phone_number"`      // Contact's phone number
	FirstName string `json:"first_name"`        // Contact's first name
	LastName  string `json:"last_name"`         // Contact's last name
	UserID    int    `json:"user_id,omitempty"` // Optional. Contact's user identifier in Telegram
}

// Location represents a point on the map.
type Location struct {
	Latitude  float64 `json:"latitude"`  // Latitude as defined by sender
	Longitude float64 `json:"longitude"` // Longitude as defined by sender
}

// UserProfilePhotos represent a user's profile pictures.
type UserProfilePhotos struct {
	Count  int            `json:"total_count"` // Total number of profile pictures the target user has
	Photos [][]*PhotoSize `json:"photos"`      // Requested profile pictures (in up to 4 sizes each)
}

// Forward denotes a message is forwarded. It is a meta-type used only in Message type
type Forward struct {
	From      *User `json:"forward_from,omitempty"` // Optional. For forwarded messages, sender of the original message
	Timestamp int64 `json:"forward_date,omitempty"` // Optional. For forwarded messages, date the original message was sent in Unix time
}

// MessageType represents type of message
type MessageType string

func (t MessageType) String() string {
	return string(t)
}

// predefined message types
const (
	CONTACT  MessageType = "Contact"
	LOCATION MessageType = "Location"
	STICKER  MessageType = "Sticker"
	PHOTO    MessageType = "Photo"
	VIDEO    MessageType = "Video"
	VOICE    MessageType = "Voice"
	AUDIO    MessageType = "Audio"
	DOCUMENT MessageType = "Document"
	TEXT     MessageType = "Text"
	STATUS   MessageType = "Status"
)

// Message represents a message.
type Message struct {
	ID                int          `json:"message_id"` // Unique message identifier
	Sender            *User        `json:"from"`       // Sender
	Timestamp         int64        `json:"date"`       // Date the message was sent in Unix time.
	Chat              *Chat        `json:"chat"`       // Conversation the message belongs to — user in case of a private message, GroupChat in case of a group
	*Forward                       // Optional
	ReplyTo           *Message     `json:"reply_to_message,omitempty"`      // Optional. For replies, the original message.
	Text              string       `json:"text,omitempty"`                  // Optional. For text messages, the actual UTF-8 text of the message
	Audio             *Audio       `json:"audio,omitempty"`                 // Optional. Message is an audio file, information about the file
	Document          *Document    `json:"document,omitempty"`              // Optional. Message is a general file, information about the file
	Photo             []*PhotoSize `json:"photo,omitempty"`                 // Optional. Message is a photo, available sizes of the photo
	Sticker           *Sticker     `json:"sticker,omitempty"`               // Optional. Message is a sticker, information about the sticker
	Video             *Video       `json:"video,omitempty"`                 // Optional. Message is a video, information about the video
	Voice             *Voice       `json:"voice,omitempty"`                 // Optional. Message is a voice message, information about the file
	Caption           string       `json:"caption,omitempty"`               // Optional. Caption for the photo or video
	Contact           *Contact     `json:"contact,omitempty"`               // Optional. Message is a shared contact, information about the contact
	Location          *Location    `json:"location,omitempty"`              // Optional. Message is a shared location, information about the location
	MemberEnter       *User        `json:"new_chat_participant,omitempty"`  // Optional. A new member was added to the group, information about them (this member may be bot itself)
	MemberLeave       *User        `json:"left_chat_participant,omitempty"` // Optional. A member was removed from the group, information about them (this member may be bot itself)
	NewTitle          string       `json:"new_chat_title,omitempty"`        // Optional. A group title was changed to this value
	NewPhoto          []*PhotoSize `json:"new_chat_photo,omitempty"`        // Optional. A group photo was change to this value
	ChatPhotoDeleted  bool         `json:"delete_chat_photo,omitempty"`     // Optional. Informs that the group photo was deleted
	GroupCreated      bool         `json:"group_created,omitempty"`         // Optional. Informs that the group has been created
	SuperGroupCreated bool         `json:"supergroup_chat_created"`         // Optional. Service message: the supergroup has been created
	ChannelCreated    bool         `json:"channel_chat_created"`            // Optional. Service message: the channel has been created
	MigrateTo         int64        `json:"migrate_to_chat_id"`              // Optional. The group has been migrated to a supergroup with the specified identifier, not exceeding 1e13 by absolute value
	MigrateFrom       int64        `json:"migrate_from_chat_id"`            // Optional. The group has been migrated to a supergroup with the specified identifier, not exceeding 1e13 by absolute value
}

// Type returns message type.
func (m Message) Type() (ret MessageType) {
	ret = STATUS
	switch {
	case m.Contact != nil:
		ret = CONTACT
	case m.Location != nil:
		ret = LOCATION
	case m.Sticker != nil:
		ret = STICKER
	case m.Photo != nil:
		ret = PHOTO
	case m.Video != nil:
		ret = VIDEO
	case m.Voice != nil:
		ret = VOICE
	case m.Audio != nil:
		ret = AUDIO
	case m.Document != nil:
		ret = DOCUMENT
	case m.Text != "":
		ret = TEXT
	}
	return
}

// Update represents an incoming update.
type Update struct {
	ID      int      `json:"update_id"` // The update‘s unique identifier.
	Message *Message `json:"message"`   // Optional. New incoming message of any kind — text, photo, sticker, etc.
}

type updates struct {
	Success bool     `json:"ok"`
	Result  []Update `json:"result"`
}
