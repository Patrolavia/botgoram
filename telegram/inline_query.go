// This file is part of Botgoram
// Botgoram is free software: see LICENSE.txt for more details.

package telegram

import "net/url"

// InlineQueryOptions represents optional parameters for calling AnswerInlineQuery
type InlineQueryOptions struct {
	CacheTime int // The maximum amount of time in seconds that the result of the inline query may be cached on the server. Defaults to 300.

	// Pass True, if results may be cached on the server side only for the user that sent the query.
	// By default, results may be returned to any user who sends the same query
	Personal bool

	// Pass the offset that a client should send in the next query with the same text to receive more results.
	// Pass an empty string if there are no more results or if you don‘t support pagination.
	// Offset length can’t exceed 64 bytes.
	NextOffset string

	// If passed, clients will display a button with specified text that switches the user to a private chat
	// with the bot and sends the bot a start message with the parameter switch_pm_parameter
	SwitchPM string

	// Parameter for the start message sent to the bot when user presses the switch button
	//
	// Example: An inline bot that sends YouTube videos can ask the user to connect the bot to their YouTube
	// account to adapt search results accordingly. To do this, it displays a ‘Connect your YouTube account’
	// button above the results, or even before showing any. The user presses the button, switches to a
	// private chat with the bot and, in doing so, passes a start parameter that instructs the bot to return
	// an oauth link. Once done, the bot can offer a switch_inline button so that the user can easily return
	// to the chat where they wanted to use the bot's inline capabilities.
	SwitchParam string
}

// Values convert this option struct to url.Values. This implements internal option interface
func (i *InlineQueryOptions) values() url.Values {
	ret := url.Values{}

	optInt(ret, "cache_time", i.CacheTime)
	optBool(ret, "is_personal", i.Personal)
	optStr(ret, "next_offset", i.NextOffset)
	optStr(ret, "switch_pm_text", i.SwitchPM)
	optStr(ret, "switch_pm_parameter", i.SwitchParam)

	return ret
}

// InlineQuery represents an incoming inline query. When the user sends an empty query, your bot could return some default or trending results.
type InlineQuery struct {
	ID       string    `json:"id"`                 // Unique identifier for this query
	From     *User     `json:"from"`               // Sender
	Location *Location `json:"location,omitempty"` // Optional. Sender location, only for bots that request user location
	Query    string    `json:"query"`              // Text of the query
	Offset   string    `json:"offset"`             // Offset of the results to be returned, can be controlled by the bot
}
