# Botgoram - State-based telegram bot framework, in Go

Botgoram is state-based telegram bot framework written in go. It is inspired by [tucnak/telebot](https://github.com/tucnak/telebot). Botgoram helps when writing large, complicated, interative bots. If you only need a small, simple, command-based bot, [tucnak/telebot](https://github.com/tucnak/telebot) or [telegram api implementation in botgoram](https://godoc.org/github.com/Patrolavia/botgoram/telegram) would be your best friend.

[![GoDoc](https://godoc.org/github.com/Patrolavia/botgoram?status.svg)](https://godoc.org/github.com/Patrolavia/botgoram)
[![Go Report Card](https://goreportcard.com/badge/github.com/Patrolavia/botgoram)](https://goreportcard.com/report/github.com/Patrolavia/botgoram)
[![Build Status](https://travis-ci.org/Patrolavia/botgoram.svg?branch=master)](https://travis-ci.org/Patrolavia/botgoram)

### State based

We think the work flow for bot is like a [Finite State Machine](https://en.wikipedia.org/wiki/Finite-state_machine): given current state, transit to next state acording to the input. We write code to choose right state, and define what to do when entering/ leaving a state.

## Synopsis

See [example code on godoc.org](https://godoc.org/github.com/Patrolavia/botgoram#example-package).

## But how can I convert my business logic to a state machine

It depends. Draw a flowchart, especially a data flowchart, and treat each unit as a state might be a reasonable start. The [state pattern](https://en.wikipedia.org/wiki/State_pattern), [Automata-based programming](https://en.wikipedia.org/wiki/Automata-based_programming) on wikipedia might also give you some thoughts.

## It looks so complicate!

[Yes, the code will be much longer.](https://en.wikipedia.org/wiki/Automata-based_programming#Automata-based_style_program) But it will also eliminates a number of control structures and function calls. And program can be faster if you apply certain optimization on your state map.

## License

Any version of MIT, GPL or LGPL. See LICENSE.txt for details.
