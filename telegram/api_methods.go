package telegram

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/url"
)

func optconv(opt *Options, u *User) (params url.Values, err error) {
	params = url.Values{}
	if opt != nil {
		params, err = opt.encode()
	}
	params.Set("chat_id", itoa(u.ID))
	return
}

func (a *api) Me() (u *User, err error) {
	params := url.Values{}
	data, err := a.sendCommand("getMe", params)
	if err != nil {
		return
	}

	u = &User{}
	err = json.Unmarshal(data, u)
	return
}

func (a *api) SendMessage(victim *User, text string, opt *Options) (m *Message, err error) {
	params, err := optconv(opt, victim)
	if err != nil {
		return
	}
	params.Set("chat_id", itoa(victim.ID))
	params.Set("text", text)
	data, err := a.sendCommand("sendMessage", params)
	if err != nil {
		return
	}

	m = &Message{}
	err = json.Unmarshal(data, m)
	return
}

func (a *api) ForwardMessage(victim, from *User, messageID int) (m *Message, err error) {
	params := url.Values{}
	params.Set("chat_id", itoa(victim.ID))
	params.Set("from_chat_id", itoa(from.ID))
	params.Set("message_id", itoa(messageID))
	data, err := a.sendCommand("forwardMessage", params)
	if err != nil {
		return
	}

	m = &Message{}
	err = json.Unmarshal(data, m)
	return
}

func (a *api) SendPhoto(victim *User, file *File, caption string, opt *Options) (m *Message, err error) {
	params, err := optconv(opt, victim)
	if err != nil {
		return
	}
	params.Set("caption", caption)
	data, err := a.sendFile("sendPhoto", "photo", file, params)
	if err != nil {
		return
	}
	m = &Message{}
	err = json.Unmarshal(data, m)
	return
}

func (a *api) SendAudio(victim *User, file *File, duration int,
	performer, title string, opt *Options) (m *Message, err error) {

	params, err := optconv(opt, victim)
	if err != nil {
		return
	}
	if duration > 0 {
		params.Set("duration", itoa(duration))
	}
	if performer != "" {
		params.Set("performer", performer)
	}
	if title != "" {
		params.Set("title", title)
	}
	data, err := a.sendFile("sendAudio", "audio", file, params)
	if err != nil {
		return
	}
	m = &Message{}
	err = json.Unmarshal(data, m)
	return
}

func (a *api) SendDocument(victim *User, file *File, opt *Options) (m *Message, err error) {
	params, err := optconv(opt, victim)
	if err != nil {
		return
	}
	data, err := a.sendFile("sendDocument", "document", file, params)
	if err != nil {
		return
	}
	m = &Message{}
	err = json.Unmarshal(data, m)
	return
}

func (a *api) SendSticker(victim *User, file *File, opt *Options) (m *Message, err error) {
	params, err := optconv(opt, victim)
	if err != nil {
		return
	}
	data, err := a.sendFile("sendSticker", "sticker", file, params)
	if err != nil {
		return
	}
	m = &Message{}
	err = json.Unmarshal(data, m)
	return
}

func (a *api) SendVideo(victim *User, file *File,
	duration int, caption string, opt *Options) (m *Message, err error) {

	params, err := optconv(opt, victim)
	if err != nil {
		return
	}
	if duration > 0 {
		params.Set("duration", itoa(duration))
	}
	if caption != "" {
		params.Set("caption", caption)
	}

	data, err := a.sendFile("sendVideo", "video", file, params)
	if err != nil {
		return
	}
	m = &Message{}
	err = json.Unmarshal(data, m)
	return
}

func (a *api) SendVoice(victim *User, file *File, duration int, opt *Options) (m *Message, err error) {
	params, err := optconv(opt, victim)
	if err != nil {
		return
	}
	if duration > 0 {
		params.Set("duration", itoa(duration))
	}

	data, err := a.sendFile("sendVoice", "voice", file, params)
	if err != nil {
		return
	}
	m = &Message{}
	err = json.Unmarshal(data, m)
	return
}

func (a *api) SendLocation(victim *User, location *Location, opt *Options) (m *Message, err error) {
	params, err := optconv(opt, victim)
	if err != nil {
		return
	}
	params.Set("latitude", ftoa(location.Latitude))
	params.Set("longitude", ftoa(location.Longitude))
	data, err := a.sendCommand("sendLocation", params)
	if err != nil {
		return
	}
	m = &Message{}
	err = json.Unmarshal(data, m)
	return
}

func (a *api) SendChatAction(victim *User, action ChatAction) (err error) {
	params := url.Values{}
	params.Set("chat_id", itoa(victim.ID))
	params.Set("action", string(action))
	_, err = a.sendCommand("sendChatAction", params)
	return
}

func (a *api) GetProfilePhotos(victim *User, offset, limit int) (p *UserProfilePhotos, err error) {
	params := url.Values{}
	params.Set("chat_id", itoa(victim.ID))
	params.Set("offset", itoa(offset))
	params.Set("limit", itoa(limit))
	data, err := a.sendCommand("getUserProfilePhotos", params)
	if err != nil {
		return
	}
	p = &UserProfilePhotos{}
	err = json.Unmarshal(data, p)
	return
}

func (a *api) GetAllProfilePhotos(victim *User) (p *UserProfilePhotos, err error) {
	params := url.Values{}
	params.Set("chat_id", itoa(victim.ID))
	data, err := a.sendCommand("getUserProfilePhotos", params)
	if err != nil {
		return
	}
	p = &UserProfilePhotos{}
	err = json.Unmarshal(data, p)
	return
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
	params.Set("offset", itoa(offset))
	if limit > 0 {
		params.Set("limit", itoa(limit))
	}
	if timeout > 0 {
		params.Set("timeout", itoa(timeout))
	}
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
	if hookURL != "" {
		params.Set("url", hookURL)
	}
	if cert != nil {
		f := &File{Filename: "server.cert", Stream: bytes.NewReader(cert)}
		_, err = a.sendFile("setWebhook", "certificate", f, params)
	} else {
		_, err = a.sendCommand("setWebhook", params)
	}
	return
}
