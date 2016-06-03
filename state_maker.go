// This file is part of Botgoram
// Botgoram is free software: see LICENSE.txt for more details.

package botgoram

// TransitorMap maps a transitor to parent state.
type TransitorMap struct {
	Transitor Transitor
	State     string // state id of parent state
	// matched first. A hidden transitor is only for generating state map.
	// It denotes a call to Transit(id).
	IsHidden   bool
	IsFallback bool // if this is a fallback transitor. matched second.
	Type       string
	Command    string // ignored if it is empty string or Type is not TEXT.
	Desc       string // only for state map generating.
}

// StateMaker helps you design you own state map by
// just implement this and call FSM.MakeState() method.
type StateMaker interface {
	Name() string // state id you want to register
	Actions() (enter Action, leave Action)
	Transitors() []TransitorMap
}
