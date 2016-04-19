// This file is part of Botgoram
// Botgoram is free software: see LICENSE.txt for more details.

package telegram

// InputMessageContent represents the content of a message to be sent as a result of an inline query.
type InputMessageContent map[string]interface{}

// SetText make this an InputTextMessageContent. See https://core.telegram.org/bots/api#inputtextmessagecontent
func (c *InputMessageContent) SetText(text string, mode ParseMode, noPreview bool) {
	*c = map[string]interface{}{"message_text": text}
	mapStr(*c, "parse_mode", string(mode))
	mapBool(*c, "disable_web_preview", noPreview)
}

// SetLocation make this an InputLocationMessageContent. See https://core.telegram.org/bots/api#inputlocationmessagecontent
func (c *InputMessageContent) SetLocation(lat, lng float64) {
	*c = map[string]interface{}{
		"latitude":  lat,
		"longitude": lng,
	}
}

// SetVenue make this an InputVenueMessageContent. See https://core.telegram.org/bots/api#inputvenuemessagecontent
func (c *InputMessageContent) SetVenue(lat, lng float64, title, address string, foursq string) {
	*c = map[string]interface{}{
		"latitude":  lat,
		"longitude": lng,
		"title":     title,
		"address":   address,
	}
	mapStr(*c, "foursquare_id", foursq)
}

// SetContact make this an InputContactMessageContent. See https://core.telegram.org/bots/api#inputcontactmessagecontent
func (c *InputMessageContent) SetContact(phone, firstName, lastName string) {
	*c = map[string]interface{}{
		"phone_number": phone,
		"first_name":   firstName,
	}
	mapStr(*c, "last_name", lastName)
}

// InlineQueryResult represents all one result of an inline query.
// Telegram clients currently support 19 result types,
// see https://core.telegram.org/bots/api#inlinequeryresult for detail
type InlineQueryResult interface {
	ForceType() string // rewrite result type to correct value
}

// AbstractInlineQueryResult abstracts some common fields of inline query result
type AbstractInlineQueryResult struct {
	Type        string               `json:"type"`
	ID          string               `json:"id"`
	ReplyMarkup *ReplyMarkup         `json:"reply_markup,omitempty"`
	Title       string               `json:"title,omitempty"`
	Caption     string               `json:"caption,omitempty"`
	Description string               `json:"description,omitempty"`
	Content     *InputMessageContent `json:"input_message_content,omitempty"`
}

// InlineQueryResultArticle represents a link to an article or web page.
type InlineQueryResultArticle struct {
	*AbstractInlineQueryResult
	URL       string `json:"url,omitempty"`          // Optional. URL of the result
	HideURL   bool   `json:"hide_url,omitempty"`     // Optional. Pass True, if you don't want the URL to be shown in the message
	Thumbnail string `json:"thumb_url,omitempty"`    // Optional. Url of the thumbnail for the result
	Width     int    `json:"thumb_width,omitempty"`  // Optional. Thumbnail width
	Height    int    `json:"thumb_height,omitempty"` // Optional. Thumbnail height
}

// ForceType forces this result type to article
func (r *InlineQueryResultArticle) ForceType() {
	r.Type = "article"
}

// InlineQueryResultPhoto represents a link to a photo.
// By default, this photo will be sent by the user with optional caption.
// Alternatively, you can use input_message_content to send a message with the specified content instead of the photo.
type InlineQueryResultPhoto struct {
	*AbstractInlineQueryResult
	Photo     string `json:"photo_url"`              // A valid URL of the photo. Photo must be in jpeg format. Photo size must not exceed 5MB
	Thumbnail string `json:"thumb_url"`              // Url of the thumbnail for the result
	Width     int    `json:"photo_width,omitempty"`  // Optional. Width of the photo
	Height    int    `json:"photo_height,omitempty"` // Optional. Height of the photo
}

// ForceType forces this result type to photo
func (r *InlineQueryResultPhoto) ForceType() {
	r.Type = "photo"
}

// InlineQueryResultGif represents a link to an animated GIF file.
// By default, this animated GIF file will be sent by the user with optional caption.
// Alternatively, you can use input_message_content to send a message with the specified content instead of the animation.
type InlineQueryResultGif struct {
	*AbstractInlineQueryResult
	GIF       string `json:"gif_url"`              // A valid URL for the GIF file. File size must not exceed 1MB
	Thumbnail string `json:"thumb_url"`            // Url of the thumbnail for the result
	Width     int    `json:"gif_width,omitempty"`  // Optional. Width of the GIF
	Height    int    `json:"gif_height,omitempty"` // Optional. Height of the GIF
}

// ForceType forces this result type to gif
func (r *InlineQueryResultGif) ForceType() {
	r.Type = "gif"
}

// InlineQueryResultMpeg4Gif represents a link to a video animation (H.264/MPEG-4 AVC video without sound).
// By default, this animated MPEG-4 file will be sent by the user with optional caption.
// Alternatively, you can use input_message_content to send a message with the specified content instead of the animation.
type InlineQueryResultMpeg4Gif struct {
	*AbstractInlineQueryResult
	Mpeg4     string `json:"mpeg4_url"`              // A valid URL for the MP4 file. File size must not exceed 1MB
	Thumbnail string `json:"thumb_url"`              // Url of the thumbnail for the result
	Width     int    `json:"mpeg4_width,omitempty"`  // Optional. Video width
	Height    int    `json:"mpeg4_height,omitempty"` // Optional. Video height
}

// ForceType forces this result type to mpeg4_gif
func (r *InlineQueryResultMpeg4Gif) ForceType() {
	r.Type = "mpeg4_gif"
}

// InlineQueryResultVideo represents a link to a page containing an embedded video player or a video file.
// By default, this video file will be sent by the user with an optional caption.
// Alternatively, you can use input_message_content to send a message with the specified content instead of the video.
type InlineQueryResultVideo struct {
	*AbstractInlineQueryResult
	Video     string `json:"video_url"`                // A valid URL for the embedded video player or video file
	MIME      string `json:"mime_type"`                // Mime type of the content of video url, “text/html” or “video/mp4”
	Thumbnail string `json:"thumb_url"`                // Url of the thumbnail for the result
	Width     int    `json:"video_width,omitempty"`    // Optional. Video width
	Height    int    `json:"video_height,omitempty"`   // Optional. Video height
	Duration  int    `json:"video_duration,omitempty"` // Optional. Video duration in seconds
}

// ForceType forces this result type to video
func (r *InlineQueryResultVideo) ForceType() {
	r.Type = "video"
}

// InlineQueryResultAudio represents a link to an mp3 audio file.
// By default, this audio file will be sent by the user.
// Alternatively, you can use input_message_content to send a message with the specified content instead of the audio.
type InlineQueryResultAudio struct {
	*AbstractInlineQueryResult
	Audio     string `json:"audio_url"`                // A valid URL for the audio file
	Performer string `json:"performer,omitempty"`      // Optional. Performer
	Duration  int    `json:"audio_duration,omitempty"` // Optional. Audio duration in seconds
}

// ForceType forces this result type to audio
func (r *InlineQueryResultAudio) ForceType() {
	r.Type = "audio"
}

// InlineQueryResultVoice represents a link to a voice recording in an .ogg container encoded with OPUS.
// By default, this voice recording will be sent by the user.
// Alternatively, you can use input_message_content to send a message with the specified content instead
// of the the voice message.
type InlineQueryResultVoice struct {
	*AbstractInlineQueryResult
	Voice    string `json:"voice_url"`                // A valid URL for the voice recording
	Duration int    `json:"voice_duration,omitempty"` // Optional. Recording duration in seconds
}

// ForceType forces this result type to voice
func (r *InlineQueryResultVoice) ForceType() {
	r.Type = "voice"
}

// InlineQueryResultDocument represents a link to a file.
// By default, this file will be sent by the user with an optional caption.
// Alternatively, you can use input_message_content to send a message with the specified content instead of the file.
// Currently, only .PDF and .ZIP files can be sent using this method.
type InlineQueryResultDocument struct {
	*AbstractInlineQueryResult
	Document  string `json:"document_url"`           // A valid URL for the file
	MIME      string `json:"mime_type"`              // Mime type of the content of the file, either “application/pdf” or “application/zip”
	Thumbnail string `json:"thumb_url,omitempty"`    // Optional. URL of the thumbnail (jpeg only) for the file
	Width     int    `json:"thumb_width,omitempty"`  // Optional. Thumbnail width
	Height    int    `json:"thumb_height,omitempty"` // Optional. Thumbnail height
}

// ForceType forces this result type to document
func (r *InlineQueryResultDocument) ForceType() {
	r.Type = "document"
}

// InlineQueryResultLocation represents a location on a map.
// By default, the location will be sent by the user.
// Alternatively, you can use input_message_content to send a message with the specified content instead of the location.
type InlineQueryResultLocation struct {
	*AbstractInlineQueryResult
	Latitude  float64 `json:"latitude"`               // Location latitude in degrees
	Longitude float64 `json:"longitude"`              // Location longitude in degrees
	Thumbnail string  `json:"thumb_url,omitempty"`    // Optional. URL of the thumbnail (jpeg only) for the file
	Width     int     `json:"thumb_width,omitempty"`  // Optional. Thumbnail width
	Height    int     `json:"thumb_height,omitempty"` // Optional. Thumbnail height
}

// ForceType forces this result type to location
func (r *InlineQueryResultLocation) ForceType() {
	r.Type = "location"
}

// InlineQueryResultVenue represents a venue.
// By default, the venue will be sent by the user.
// Alternatively, you can use input_message_content to send a message with the specified content instead of the venue.
type InlineQueryResultVenue struct {
	*InlineQueryResultLocation
	Address      string `json:"address"`                 // Address of the venue
	FourSquareID string `json:"foursquare_id,omitempty"` // Optional. Foursquare identifier of the venue if known
}

// ForceType forces this result type to venue
func (r *InlineQueryResultVenue) ForceType() {
	r.Type = "venue"
}

// InlineQueryResultContact represents a contact with a phone number.
// By default, this contact will be sent by the user.
// Alternatively, you can use input_message_content to send a message with the specified content instead of the contact.
type InlineQueryResultContact struct {
	*AbstractInlineQueryResult
	Phone     string `json:"phone_number"`           // Contact's phone number
	FirstName string `json:"first_name"`             // Contact's first name
	LastName  string `json:"last_name,omitempty"`    // Optional. Contact's last name
	Thumbnail string `json:"thumb_url,omitempty"`    // Optional. URL of the thumbnail (jpeg only) for the file
	Width     int    `json:"thumb_width,omitempty"`  // Optional. Thumbnail width
	Height    int    `json:"thumb_height,omitempty"` // Optional. Thumbnail height
}

// ForceType forces this result type to contact
func (r *InlineQueryResultContact) ForceType() {
	r.Type = "contact"
}

// InlineQueryResultCachedPhoto represents a link to a photo stored on the Telegram servers.
// By default, this photo will be sent by the user with an optional caption.
// Alternatively, you can use input_message_content to send a message with the specified content instead of the photo.
type InlineQueryResultCachedPhoto struct {
	*AbstractInlineQueryResult
	PhotoID string `json:"photo_file_id"` // A valid file identifier of the photo
}

// ForceType forces this result type to photo
func (r *InlineQueryResultCachedPhoto) ForceType() {
	r.Type = "photo"
}

// InlineQueryResultCachedGif represents a link to an animated GIF file stored on the Telegram servers.
// By default, this animated GIF file will be sent by the user with an optional caption.
// Alternatively, you can use input_message_content to send a message with specified content instead of the animation.
type InlineQueryResultCachedGif struct {
	*AbstractInlineQueryResult
	GIFID string `json:"gif_file_id"` // A valid file identifier for the GIF file
}

// ForceType forces this result type to gif
func (r *InlineQueryResultCachedGif) ForceType() {
	r.Type = "gif"
}

// InlineQueryResultCachedMpeg4Gif represents a link to a video animation (H.264/MPEG-4 AVC video without sound)
// stored on the Telegram servers.
// By default, this animated MPEG-4 file will be sent by the user with an optional caption.
// Alternatively, you can use input_message_content to send a message with the specified content instead of the animation.
type InlineQueryResultCachedMpeg4Gif struct {
	*AbstractInlineQueryResult
	Mpeg4ID string `json:"mpeg4_file_id"` // A valid file identifier for the MP4 file
}

// ForceType forces this result type to mpeg4_gif
func (r *InlineQueryResultCachedMpeg4Gif) ForceType() {
	r.Type = "mpeg4_gif"
}

// InlineQueryResultCachedSticker represents a link to a sticker stored on the Telegram servers.
// By default, this sticker will be sent by the user.
// Alternatively, you can use input_message_content to send a message with the specified content instead of the sticker.
type InlineQueryResultCachedSticker struct {
	*AbstractInlineQueryResult
	StickerID string `json:"sticker_file_id"` // A valid file identifier of the sticker
}

// ForceType forces this result type to sticker
func (r *InlineQueryResultCachedSticker) ForceType() {
	r.Type = "sticker"
}

// InlineQueryResultCachedDocument represents a link to a file stored on the Telegram servers.
// By default, this file will be sent by the user with an optional caption.
// Alternatively, you can use input_message_content to send a message with the specified content instead of the file.
// Currently, only pdf-files and zip archives can be sent using this method.
type InlineQueryResultCachedDocument struct {
	*AbstractInlineQueryResult
	DocumentID string `json:"document_file_id"` // A valid file identifier for the file
}

// ForceType forces this result type to document
func (r *InlineQueryResultCachedDocument) ForceType() {
	r.Type = "document"
}

// InlineQueryResultCachedVideo represents a link to a video file stored on the Telegram servers.
// By default, this video file will be sent by the user with an optional caption.
// Alternatively, you can use input_message_content to send a message with the specified content instead of the video.
type InlineQueryResultCachedVideo struct {
	*AbstractInlineQueryResult
	VideoID string `json:"video_file_id"` // A valid file identifier for the video file
}

// ForceType forces this result type to video
func (r *InlineQueryResultCachedVideo) ForceType() {
	r.Type = "video"
}

// InlineQueryResultCachedVoice represents a link to a voice message stored on the Telegram servers.
// By default, this voice message will be sent by the user.
// Alternatively, you can use input_message_content to send a message with the specified content instead of the voice message.
type InlineQueryResultCachedVoice struct {
	*AbstractInlineQueryResult
	VoiceID string `json:"voice_file_id"` // A valid file identifier for the voice message
}

// ForceType forces this result type to voice
func (r *InlineQueryResultCachedVoice) ForceType() {
	r.Type = "voice"
}

// InlineQueryResultCachedAudio represents a link to an mp3 audio file stored on the Telegram servers.
// By default, this audio file will be sent by the user.
// Alternatively, you can use input_message_content to send a message with the specified content instead of the audio.
type InlineQueryResultCachedAudio struct {
	*AbstractInlineQueryResult
	AudioID string `json:"audio_file_id"` // A valid file identifier for the audio file
}

// ForceType forces this result type to audio
func (r *InlineQueryResultCachedAudio) ForceType() {
	r.Type = "audio"
}

// ChosenInlineResult represents a result of an inline query that was chosen by the user and sent to their chat partner.
type ChosenInlineResult struct {
	ID       string    `json:"result_id"`          // The unique identifier for the result that was chosen
	From     *User     `json:"from"`               // The user that chose the result
	Location *Location `json:"location,omitempty"` // Optional. Sender location, only for bots that require user location
	// Optional. Identifier of the sent inline message.
	// Available only if there is an inline keyboard attached to the message.
	// Will be also received in callback queries and can be used to edit the message.
	MessageID string `json:"inline_messsage_id,omitempty"`
	Query     string `json:"query"` // The query that was used to obtain the result
}
