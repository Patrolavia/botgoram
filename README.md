# Botgoram - State-based telegram bot framework, in Go

Botgoram is state-based telegram bot framework written in go. It is inspired by [tucnak/telebot](https://github.com/tucnak/telebot).

Current it is still under development, not usable now.

### State based

We think the work flow for bot is like a [Finite State Machine](https://en.wikipedia.org/wiki/Finite-state_machine): given current state, transit to next state acording to the input. We write code to choose right state, and migrate from one state to another.

The data stored in state will be serialized to JSON format for persistence.

## Synopsis

```go
import (
	"fmt"
	"log"

	"github.com/Patrolavia/botgoram/telegram"
)

// MyStateData defines what data to store in state.
type MyStateData struct {
	Name string
}

func main() {
	fsm := botgoram.New("my_telegram_token", new(MyStateData))

	fsm.InitialState.RegisterCommand("/greet", func(msg *telegram.Message, cur botgoram.State) string {
		cur.Bot().SendMessage(msg.Sender, "Please input your name.", nil)
		return "/greet:ask name"
	})

	fsm.NewState("/greet:ask name").RegisterTextMessage(func(msg *telegram.Message, cur botgoram.State) string {
		if len(msg.Text) < 4 {
			cur.Bot().SendMessage(msg.Sender, "Name too short, at least 4 characters.", nil)
			return cur.Name()
		}
		cur.Data().Name = msg.Text
		cur.Bot().SendMessage(msg.Sender, "Please input your title.", nil)
		return "/greet:ask title"
	})

	fsm.NewState("/greet:ask title").RegisterTextMessage(func(msg *telegram.Message, cur botgoram.State) string {
		cur.Bot().SendMessage(msg.Sender, fmt.Sprintf("Hello. %s %s", msg.Text, cur.Data.Name), nil)
		return "" // back to initial state
	})

	if err := fsm.Start(); err != nil {
		log.Fatal("something goes wrong with your state map: %s", err)
	}
}
```

## License

Any version of MIT, GPL or LGPL.
