// This file is part of Botgoram
// Botgoram is free software: see LICENSE.txt for more details.

package telegram

import "sync"

// SimpleBot is designed for people who making simple bot
type SimpleBot struct {
	API     API
	offset  int
	Limit   int
	Timeout int

	// set these properties only if you need this type of update
	// unwanted update will be skipped
	Messages            chan *Message
	InlineQueries       chan *InlineQuery
	ChosenInlineResults chan *ChosenInlineResult
	CallbackQueries     chan *CallbackQuery

	err  error
	lock *sync.Mutex
}

// Run runs the bot in foreground
func (b *SimpleBot) Run() error {
	b.Start()
	return b.Wait()
}

func (b *SimpleBot) run() {
	defer b.lock.Unlock()

	for {
		updates, err := b.API.GetUpdates(b.offset, b.Limit, b.Timeout)
		if err != nil {
			b.err = err
			break
		}

		for _, u := range updates {
			switch {
			case u.Message != nil && b.Messages != nil:
				b.Messages <- u.Message
			case u.InlineQuery != nil && b.InlineQueries != nil:
				b.InlineQueries <- u.InlineQuery
			case u.ChosenInlineResult != nil && b.ChosenInlineResults != nil:
				b.ChosenInlineResults <- u.ChosenInlineResult
			case u.CallbackQuery != nil && b.CallbackQueries != nil:
				b.CallbackQueries <- u.CallbackQuery
			}
			b.offset = u.ID + 1
		}
	}
}

// Start runs the bot in background, this can execute only once
func (b *SimpleBot) Start() {
	b.lock.Lock()
	if b.err != nil {
		return
	}
	go b.run()
}

// Wait blocks until something goes wrong
func (b *SimpleBot) Wait() (err error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	err = b.err
	b.err = nil
	return err
}
