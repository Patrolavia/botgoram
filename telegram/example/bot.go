package main

import (
	"errors"
	"sync"

	"ronmi.tw/git/Patrolavia/botgoram/telegram"
)

// Bot is a bot fetching update throuth long polling method.
type Bot interface {
	telegram.API
	Run(channel chan *telegram.Message, timeout int) error // timeout is long polling timeout in seconds
	Err() error                                            // check if any error occured when fetching updates. Error will be cleared after retriving.
}

type bot struct {
	telegram.API
	*sync.Mutex
	err error
}

func NewBot(token string) Bot {
	return &bot{telegram.New(token), &sync.Mutex{}, nil}
}

func (b *bot) Run(channel chan *telegram.Message, timeout int) error {
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
				if update.Id >= offset {
					offset = update.Id + 1
				}
				channel <- update.Message
			}
		}
	}()
	return nil
}

func (b *bot) Err() (err error) {
	err = b.err
	if b.err != nil {
		b.Lock()
		b.err = nil
		b.Unlock()
	}
	return
}
