package main

import (
	"errors"
	"sync"

	"github.com/Patrolavia/botgoram/telegram"
)

// Bot is a bot fetching update throuth long polling method.
type Bot interface {
	telegram.API
	Run(msgs chan *telegram.Message, queries chan *telegram.InlineQuery, chosen chan *telegram.ChosenInlineResult, timeout int) error // timeout is long polling timeout in seconds
	Err() error                                                                                                                       // check if any error occured when fetching updates. Error will be cleared after retriving.
}

type bot struct {
	telegram.API
	*sync.Mutex
	err error
}

func newBot(token string) Bot {
	return &bot{telegram.New(token), &sync.Mutex{}, nil}
}

func (b *bot) Run(msgs chan *telegram.Message, queries chan *telegram.InlineQuery, chosen chan *telegram.ChosenInlineResult, timeout int) error {
	u, err := b.Me()
	if u == nil || err != nil {
		return errors.New("Unable get bot information, is token valid?")
	}
	go func() {
		offset := 0
		for {
			updates, e := b.GetUpdates(offset, 0, timeout)
			if e != nil {
				b.Lock()
				b.err = e
				b.Unlock()
			}
			for _, update := range updates {
				if update.ID >= offset {
					offset = update.ID + 1
				}
				if update.Message != nil && msgs != nil {
					msgs <- update.Message
				}
				if update.InlineQuery != nil && queries != nil {
					queries <- update.InlineQuery
				}
				if update.ChosenInlineResult != nil && chosen != nil {
					chosen <- update.ChosenInlineResult
				}

			}
		}
	}()
	return nil
}

func (b *bot) Err() (err error) {
	err = b.err
	if err != nil {
		b.Lock()
		b.err = nil
		b.Unlock()
	}
	return
}
