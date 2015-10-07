package botgoram

import "github.com/Patrolavia/botgoram/telegram"

// TransitorMap maps a transitor to parent state.
type TransitorMap struct {
	Transitor Transitor
	State string // state id of parent state
	IsFallback bool // if this is a fallback transitor. matched fist.
	Type telegram.MessageType
	Command string // ignored if it is empty string or Type is not TEXT.
}

// StateMaker helps you design you own state map by
// just implement this and call FSM.MakeState() method.
type StateMaker interface {
	Name() string // state id you want to register
	Enter(msg *telegram.Message, current State, api telegram.API) error
	Leave(msg *telegram.Message, current State, api telegram.API) error
	Transitors() []TransitorMap
}
