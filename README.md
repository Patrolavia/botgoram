# Botgoram - State-based telegram bot framework, in Go

Botgoram is state-based telegram bot framework written in go. It is inspired by [tucnak/telebot](https://github.com/tucnak/telebot).

[![GoDoc](https://godoc.org/github.com/Patrolavia/botgoram?status.svg)](https://godoc.org/github.com/Patrolavia/botgoram)

Current it is still under development, not usable for now.

### State based

We think the work flow for bot is like a [Finite State Machine](https://en.wikipedia.org/wiki/Finite-state_machine): given current state, transit to next state acording to the input. We write code to choose right state, and to define what to to when entering/ leaving a state.

## Synopsis

```go
// create a fsm, using message.Sender.Id as primary key,
// up to 5 users at the same time, save state data (integer) in memory.
fsm, err := botgoram.NewBySender(
	"my-token",
	botgoram.MemoryStore(func(uid int) interface{} {
		return 0
	}),
	5)

// allocate a state
fsm.AddState(
	"hello world",
	func(msg *telegram.Message, current botgoram.State, api telegram.API) error {
		api.SendMessage(current.User(), "Hello, you're entring world state", nil)
		current.Transit(botgoram.StateId("")) // go back to initial state
	},
	func(msg *telegram.Message, current botgoram.State, api telegram.API) error {
		api.SendMessage(current.User(), "Hello, you're leaving world state", nil)
	})

init_state, ok := fsm.State(botgoram.StateId(""))
if !ok {
    panic("Cannot get initial state.")
}

init_state.RegisterFallback(
	func(m *telegram.Message, d interface{}, u *telegram.User, c botgoram.StateId) (botgoram.StateId, error) {
		return botgoram.StateId("hello world") nil
})

log.Fatal(fsm.Start(30))
```

## It looks so complicate!

[Yes, the code will be much longer.](https://en.wikipedia.org/wiki/Automata-based_programming#Automata-based_style_program) But it will also eliminates a number of control structures and function calls. And program can be faster if you apply certain optimization on your state map.

## License

Any version of MIT, GPL or LGPL.
