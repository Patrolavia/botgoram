/*
Package botgoram is a state-based Telegram bot framework.

	import (
		"fmt"
		"log"

		"github.com/Patrolavia/botgoram"
		"github.com/Patrolavia/botgoram/telegram"
	)

	func errState(fsm botgoram.FSM) botgoram.State {
		st, err := fsm.AddState(
			"error",
			func(msg *telegram.Message, current botgoram.State, api telegram.API) error { // enter
				api.SendMessage(current.User(), `/start - start using greeting bot.`, nil)
				current.Transit(botgoram.InitialState)
				return nil
			},
			nil,
		)
		if err != nil {
			log.Fatalf("Cannot add error state: %s", err)
		}
		return st
	}

	func askName(fsm botgoram.FSM) botgoram.State {
		st, err := fsm.AddState(
			"ask name",
			func(msg *telegram.Message, current botgoram.State, api telegram.API) error { // enter
				api.SendMessage(current.User(), "Please tell me your name.", nil)
				return nil
			},
			func(msg *telegram.Message, current botgoram.State, api telegram.API) error { // leave
				current.SetData(msg.Text)
				return nil
			},
		)
		if err != nil {
			log.Fatalf("Cannot register ask name state: %s", err)
		}
		return st
	}

	func askTitle(fsm botgoram.FSM) botgoram.State {
		st, err := fsm.AddState(
			"ask title",
			func(msg *telegram.Message, current botgoram.State, api telegram.API) error { // enter
				api.SendMessage(current.User(), "Please tell me your name.", nil)
				return nil
			},
			func(msg *telegram.Message, current botgoram.State, api telegram.API) error { // leave
				name := current.Data().(string)
				api.SendMessage(current.User(), fmt.Sprintf("Hello, %s %s", msg.Text, name), nil)
				return nil
			},
		)
		if err != nil {
			log.Fatalf("Cannot register ask name state: %s", err)
		}
		return st
	}

	func main() {
		token := "my-token"
		// create a FSM, stroe state data in memory, serve at most
		// 5 different users at the same time.
		fsm, err := botgoram.NewBySender(token, botgoram.MemoryStore(func(uid int) interface{} {
			// we store only string data in state
			return ""
		}), 5)
		if err != nil {
			log.Fatalf("Cannot create FSM: %s", err)
		}

		// register our command at initial state
		st, ok := fsm.State(botgoram.InitialState)
		if !ok {
			log.Fatal("Cannot get initial state")
		}
		st.RegisterCommand(
			"/start",
			func(msg *telegram.Message, data interface{}, user *telegram.User, sid string) (next string, err error) {
				return "ask name", nil
			},
		)
		// error handling, any error will go to error state
		err_handler := func(
			msg *telegram.Message,
			data interface{},
			user *telegram.User,
			sid string,
		) (next string, err error) {
			return "error", err
		}
		st.RegisterFallback(err_handler)

		// define error state: show help message, then go bak to initial state
		errState(fsm)

		// the state asking user name
		askName(fsm).Register(
			telegram.TEXT,
			func(msg *telegram.Message, data interface{}, user *telegram.User, sid string) (next string, err error) {
				next = sid
				if msg.Text != "" {
					next = "ask title"
				}
				return
			},
		)

		askTitle(fsm).Register(
			telegram.TEXT,
			func(msg *telegram.Message, data interface{}, user *telegram.User, sid string) (next string, err error) {
				next = sid
				if msg.Text != "" {
					next = botgoram.InitialState
				}
				return
			},
		)

		log.Fatal(fsm.Start(30))
	}


State based

We think the work flow for bot is like a Finite State Machine (FSM): given
current state, transit to another state according to the input. With botgoram,
we write code to choose right state, and migrate from one to another.

Why FSM

There are some algorithms to simplify a FSM. So you can design your bot as FSM,
apply the optimization, and happilly announce your high-efficiency bot.

Transitor

A transitor parses the current state and / or message to decide which state we
should move to. There SHOULD NOT be side-effects in it.

Action

Action is the code executed when you leaving / entering a state. You should put
the code with side-effects here, like modifying state data, sending telegram
messages, recording something interesting into db or so.

Initial State

We use special id "" for initial state. You cannot register enter / leave actions
for initial state: It is just the beginning (and the end) of the path travelling
through states.
*/
package botgoram
