/*
Package botgoram is a state-based Telegram bot framework.


State based

We think the work flow for bot is like a Finite State Machine (FSM): given
current state, transit to another state according to the input. With botgoram,
we write code to choose right state, and migrate from one to another.

Why FSM

There are some algorithms to simplify a FSM. So you can design your bot as FSM,
apply the optimization, and happilly announce your high-efficiency bot.

Transitor

A transitor parses the current state and/or message to decide which state we
should move to. There SHOULD NOT be side-effects in it.

Action

Action is the code executed when you leaving/entering a state. You should put
the code with side-effects here, like modifying state data, sending telegram
messages, recording something interesting into db or so.

Initial State

We use special id ""(empty string) for initial state. You cannot register enter/leave actions
for initial state: It is just the beginning (and the end) of the path travelling
through states.
*/
package botgoram
