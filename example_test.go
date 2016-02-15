package botgoram

import (
	"fmt"
	"log"
	"os"

	"github.com/Patrolavia/botgoram/telegram"
)

// HelloState defines a state, in which we reply hello message to any user who send anything to bot.
type HelloState string

// this defines what we should do when entering this state
func (h HelloState) enter(msg *telegram.Message, current State, api telegram.API) error {
	reply := fmt.Sprintf("Hello, %s.\nThe message type of previous message: %s",
		msg.Sender.FirstName, current.Data().(string))
	api.SendMessage(msg.Sender, reply, nil)
	current.SetData(msg.Type)
	current.Transit(InitialState)
	return nil
}

// this will register as fallback transitor, which will receive all messages if no other transitor.
func (h HelloState) fallbackTransitor(msg *telegram.Message, state State) (next string, err error) {
	return string(h), nil
}

// state name
func (h HelloState) Name() string {
	return string(h)
}

// fsm calls this method to register enter/leave actions.
func (h HelloState) Actions() (Action, Action) {
	return h.enter, nil
}

// fsm calls this method to register transitors and/or generate state map
func (h HelloState) Transitors() []TransitorMap {
	return []TransitorMap{
		TransitorMap{
			Transitor:  h.fallbackTransitor,
			State:      "dispatch",
			IsFallback: true,
		},
		TransitorMap{
			IsHidden: true,
			State:    "",
			Desc:     "Work done, back to initial state",
		},
	}
}

func Example() {
	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatal("Please fill your bot token into the environmental variable 'TOKEN' to use this example")
	}
	api := telegram.New(token)

	// validate token
	if _, err := api.Me(); err != nil {
		log.Fatalf("Failed to validate token: %s", err)
	}

	// create FSM
	fsm := NewBySender(
		api,
		// store state data (a string) using default memory storage
		MemoryStore(func(uid string) interface{} {
			return "No previous message"
		}),
		10, // at most process 10 users' message at the same time
	)

	fsm.MakeState(HelloState("hello"))

	// generate state map
	fmt.Println(fsm.StateMap("initial state"))

	// main loop, uncomment to do the real stuff
	/*
		for err := fsm.Start(30); err != nil; err = fsm.Resume() {
			log.Printf("There is something goes wrong: %s", err)
		}
	*/
}
