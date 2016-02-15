// This file is part of Botgoram
// Botgoram is free software: see LICENSE.txt for more details.

package telegram

import "encoding/json"

// InlineQuery represents an incoming inline query. When the user sends an empty query, your bot could return some default or trending results.
type InlineQuery struct {
	ID     string `json:"id"`     // Unique identifier for this query
	From   *User  `json:"from"`   // Sender
	Query  string `json:"query"`  // Text of the query
	Offset string `json:"offset"` // Offset of the results to be returned, can be controlled by the bot
}

// ChosenInlineResult represents a result of an inline query that was chosen by the user and sent to their chat partner.
type ChosenInlineResult struct {
	ID    string `json:"result_id"` // The unique identifier for the result that was chosen
	From  *User  `json:"from"`      // The user that chose the result
	Query string `json:"query"`     // The query that was used to obtain the result
}

// InlineQueryResult represents all one result of an inline query.
// Telegram clients currently support results of the following 5 types: article, photo, git, mpeg4gif and video.
type InlineQueryResult interface {
	Type() string
	ID() string
	MarshalJSON() ([]byte, error)
}

// InlineQueryResultArticle represents a link to an article or web page.
type InlineQueryResultArticle interface {
	InlineQueryResult

	Description(desc string)
	ParseMode(mode ParseMode)
	DisableWebPreview()

	URL(u string)
	HideURL()
	Thumb(thumb string, width, height int) // width and height is optional, use number <= 0 to bypass them
}

// NewArticleResult creates a new article result for inline query.
func NewArticleResult(id, title, message string) InlineQueryResultArticle {
	return &iqr{
		id:      id,
		resType: "article",
		data: map[string]interface{}{
			"title":        title,
			"message_text": message,
		},
	}
}

// InlineQueryResultPhoto represents a link to a photo.
// By default, this photo will be sent by the user with optional caption.
// Alternatively, you can provide message_text to send it instead of photo.
type InlineQueryResultPhoto interface {
	InlineQueryResult
	Title(t string)
	Description(desc string)
	Message(m string)
	ParseMode(mode ParseMode)
	DisableWebPreview()
	Caption(c string)

	Photo(photo string, width, height int) // width and height is optional, use number <= 0 to bypass them
}

// NewPhotoResult creates a new photo result for inline query.
func NewPhotoResult(id, photo, thumb string) InlineQueryResultPhoto {
	return &iqr{
		id:      id,
		resType: "photo",
		data: map[string]interface{}{
			"photo_url": photo,
			"thumb_url": thumb,
		},
	}
}

// InlineQueryResultGif represents a link to an animated GIF file.
// By default, this animated GIF file will be sent by the user with optional caption.
// Alternatively, you can provide message_text to send it instead of the animation.
type InlineQueryResultGif interface {
	InlineQueryResult
	Title(t string)
	Description(desc string)
	Message(m string)
	ParseMode(mode ParseMode)
	DisableWebPreview()
	Caption(c string)

	Gif(gif string, width, height int) // width and height is optional, use number <= 0 to bypass them
}

// NewGifResult creates a new gif result for inline query.
func NewGifResult(id, gif, thumb string) InlineQueryResultGif {
	return &iqr{
		id:      id,
		resType: "gif",
		data: map[string]interface{}{
			"gif_url":   gif,
			"thumb_url": thumb,
		},
	}
}

// InlineQueryResultMpeg4Gif represents a link to a video animation (H.264/MPEG-4 AVC video without sound).
// By default, this animated MPEG-4 file will be sent by the user with optional caption.
// Alternatively, you can provide message_text to send it instead of the animation.
type InlineQueryResultMpeg4Gif interface {
	InlineQueryResult
	Title(t string)
	Description(desc string)
	Message(m string)
	ParseMode(mode ParseMode)
	DisableWebPreview()
	Caption(c string)

	Mpeg4(mp4 string, width, height int) // width and height is optional, use number <= 0 to bypass them
}

// NewMpeg4Result creates a new mpeg4 result for inline query.
func NewMpeg4Result(id, mpeg4, thumb string) InlineQueryResultMpeg4Gif {
	return &iqr{
		id:      id,
		resType: "mpeg4_gif",
		data: map[string]interface{}{
			"mpeg4_url": mpeg4,
			"thumb_url": thumb,
		},
	}
}

// InlineQueryResultVideo represents link to a page containing an embedded video player or a video file.
type InlineQueryResultVideo interface {
	InlineQueryResult
	Title(t string)
	Description(desc string)
	Message(m string)
	ParseMode(mode ParseMode)
	DisableWebPreview()
	Caption(c string)

	Video(video string, width, height, duration int) // width, height and duration is optional, use number <= 0 to bypass them
}

// NewVideoResult creates a new video result for inline query.
func NewVideoResult(id, video, thumb string) InlineQueryResultVideo {
	return &iqr{
		id:      id,
		resType: "video",
		data: map[string]interface{}{
			"video_url": video,
			"thumb_url": thumb,
		},
	}
}

type iqr struct {
	id      string
	resType string
	data    map[string]interface{} // inline query result implementation
}

func (i *iqr) Type() string {
	return i.resType
}

func (i *iqr) ID() string {
	return i.id
}

func (i *iqr) MarshalJSON() ([]byte, error) {
	i.data["id"] = i.id
	i.data["type"] = i.resType
	return json.Marshal(i.data)
}

func (i *iqr) ParseMode(mode ParseMode) {
	i.data["parse_mode"] = string(mode)
}

func (i *iqr) DisableWebPreview() {
	i.data["disable_web_page_preview"] = true
}

func (i *iqr) URL(u string) {
	i.data["url"] = u
}

func (i *iqr) HideURL() {
	i.data["hide_url"] = true
}

func (i *iqr) Description(desc string) {
	i.data["description"] = desc
}

func (i *iqr) Title(t string) {
	i.data["title"] = t
}

func (i *iqr) graphics(prefix, u string, w, h int) {
	i.data[prefix+"_url"] = u
	if w > 0 {
		i.data[prefix+"_width"] = w
	}
	if h > 0 {
		i.data[prefix+"_height"] = h
	}
}

func (i *iqr) Thumb(thumb string, width, height int) {
	i.graphics("thumb", thumb, width, height)
}

func (i *iqr) Photo(photo string, width, height int) {
	i.graphics("photo", photo, width, height)
}

func (i *iqr) Gif(gif string, width, height int) {
	i.graphics("gif", gif, width, height)
}

func (i *iqr) Mpeg4(mpeg4 string, width, height int) {
	i.graphics("mpeg4", mpeg4, width, height)
}

func (i *iqr) Video(video string, width, height, duration int) {
	i.graphics("video", video, width, height)
	if duration > 0 {
		i.data["video_duration"] = duration
	}
}

func (i *iqr) Caption(c string) {
	i.data["caption"] = c
}

func (i *iqr) Message(m string) {
	i.data["message_text"] = m
}
