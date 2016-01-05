// Example bot for testing Telegram APIs
//
// How to use
//
// You have to set an environment variable named "TELEGRAM_TOKEN", fill your bot token with it, then build and run.
// Send "/help" to see what commmand you can use.
// Send text, file or something else to test reply function.
//
// Test video is downloaded from https://www.youtube.com/watch?v=SyOvMDYD4PE which licensed under CC BY license.
//
// Test image is logo of patrolavia studio, copyright 2015-, patrolavia studio.
//
// Test audio is downloaded from http://soundbible.com/royalty-free-sounds-1.html
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Patrolavia/botgoram/telegram"
)

var token string

func init() {
	token = os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		log.Fatal("You have to set environment variable TELEGRAM_TOKEN first!")
	}
}

func main() {
	bot := newBot(token)
	msgs := make(chan *telegram.Message)
	queries := make(chan *telegram.InlineQuery)
	chosen := make(chan *telegram.ChosenInlineResult)
	if err := bot.Run(msgs, queries, chosen, 30); err != nil {
		log.Fatalf("Error running bot: %s", err)
	}
	log.Print("Bot started.")

	go doqueries(queries, bot)
	go dochosen(chosen, bot)
	domsg(msgs, bot)
}

func doqueries(queries chan *telegram.InlineQuery, bot telegram.API) {
	for q := range queries {
		log.Printf("Get inline query: %#v", *q)
		r := make([]telegram.InlineQueryResult, 2)
		a := telegram.NewArticleResult("1", "article", "Query: "+q.Query)
		a.URL("https://patrolavia.com")
		a.Thumb("https://martingallagher.com/images/gopher.svg", 0, 0)
		r[0] = a
		b := telegram.NewPhotoResult(
			"2",
			"https://patrolavia.com/logo128.png",
			"https://patrolavia.com/logo64.png",
		)
		b.Title("photo for " + q.Query)
		r[1] = b
		if err := bot.AnswerInlineQuery(q, r, 1, false, ""); err != nil {
			log.Fatalf("Error when sending inline query result: %s", err)
		}
	}
}

func dochosen(chosen chan *telegram.ChosenInlineResult, bot telegram.API) {
	for c := range chosen {
		log.Printf("Got chosen: %#v", *c)
	}
}

func domsg(msgs chan *telegram.Message, bot telegram.API) {
	for msg := range msgs {
		log.Printf("Got message: %#v", *msg)
		switch msg.Text {
		case "/help":
			bot.SendMessage(msg.Chat, `/help Show this message.
/text Send text message.
/doc Send test photo as document.
/photo Send test photo.
/audio Send test audio.
/voice Send test audio as voice.
/video Send test video.
/loc Send test location.
/forward Forward last message.

or you can send message to me, I will reply it with some debug message.`, nil)
		case "/text":
			bot.SendMessage(msg.Chat, `Hello, `+msg.Sender.FirstName, nil)
		case "/doc":
			bot.SendDocument(msg.Chat, &telegram.File{Filename: "test.png"}, nil)
		case "/photo":
			bot.SendPhoto(msg.Chat, &telegram.File{Filename: "test.png"}, "Test photo", nil)
		case "/audio":
			bot.SendChatAction(msg.Chat, telegram.UploadAudio)
			bot.SendAudio(msg.Chat, &telegram.File{Filename: "test.ogg"}, 7, "", "Test audio", nil)
		case "/voice":
			bot.SendChatAction(msg.Chat, telegram.RecordAudio)
			bot.SendVoice(msg.Chat, &telegram.File{Filename: "test.ogg"}, 7, nil)
		case "/video":
			bot.SendChatAction(msg.Chat, telegram.UploadVideo)
			bot.SendVideo(msg.Chat, &telegram.File{Filename: "test.mp4"}, 108, "Test video", nil)
		case "/loc":
			bot.SendLocation(msg.Chat, &telegram.Location{24.1501297, 120.6863541}, nil)
		case "/forward":
			bot.ForwardMessage(msg.Chat, msg.Chat, msg.ID)
		default:
			reply(msg, bot)
		}
	}
}

func reply(msg *telegram.Message, bot telegram.API) {
	senderType := "User"
	if msg.Chat.Type != telegram.TYPECHAT {
		senderType = "Chatroom"
	}
	msgType := msg.Type().String() + " "
	if msg.ReplyTo != nil {
		msgType += "(reply)"
	}
	if msg.Forward != nil {
		msgType += "(forward from " + msg.Forward.From.Name() + ")"
	}
	senderName := msg.Sender.Name()
	opt := &telegram.Options{ReplyTo: msg.ID}
	bot.SendMessage(msg.Chat, fmt.Sprintf(`Chatroom title or user name: %s
Send from: %s
Sender: %s
Message type: %s`, msg.Chat.Name(), senderType, senderName, msgType), opt)
}
