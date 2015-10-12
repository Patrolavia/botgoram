package botgoram

import (
	"errors"
	"fmt"

	"github.com/Patrolavia/botgoram/telegram"
)

var ErrStateNotFound error = errors.New("State not found.")

// Action describes what to do when enter/leaving a state.
type Action func(msg *telegram.Message, current State, api telegram.API) error

// FSM is a finite state machine.
type FSM interface {
	// Start will "power on" the FSM.
	// It will block until any worker goes wrong, stop that worker and return the error.
	// At the mean time, other worker will stay working.
	Start(timeout int) error
	// Resume the stopped worker. Calling Resume will block until any worker goes wrong.
	Resume() error
	AddState(id string, enter, leave Action) (State, error)
	State(id string) (State, bool)
	// MakeState will register a new state with StateMaker.
	// Transitors will be registered when first call to FSM.Start()
	MakeState(StateMaker) (State, error)
	// StateMap generate graphviz diagram from registered StateMaker
	StateMap(name string) (dot string)
}

func bySender(msg *telegram.Message) *telegram.User {
	return msg.Sender
}

func byChat(msg *telegram.Message) *telegram.User {
	return msg.Chat
}

type internalStateData struct {
	state State
	enter Action
	leave Action
}

type fsm struct {
	api           telegram.API
	userExtractor func(*telegram.Message) *telegram.User
	states        map[string]internalStateData
	storage       SaveLoader
	manager       *manager
	error_chan    chan error
	sm            []StateMaker
}

func newFSM(token string, ue func(*telegram.Message) *telegram.User, sl SaveLoader, size int) (ret FSM, err error) {
	api := telegram.New(token)

	// validate token
	if _, err = api.Me(); err != nil {
		return
	}

	tmp := &fsm{
		api, ue, map[string]internalStateData{
			"": internalStateData{state: newState("")},
		}, sl, newManager(ue, size),
		make(chan error, size),
		make([]StateMaker, 0),
	}
	for i := 0; i < size; i++ {
		tmp.error_chan <- nil
	}

	return tmp, err
}

// NewBySender creates a FSM associates with message sender, and test if token is valid.
func NewBySender(token string, sl SaveLoader, size int) (FSM, error) {
	return newFSM(token, bySender, sl, size)
}

// NewByChat creates a FSM associates with chatroom, and test if token is valid.
func NewByChat(token string, sl SaveLoader, size int) (FSM, error) {
	return newFSM(token, byChat, sl, size)
}

func (f *fsm) MakeState(sm StateMaker) (ret State, err error) {
	enter, leave := sm.Actions()
	if ret, err = f.AddState(sm.Name(), enter, leave); err != nil {
		return
	}

	f.sm = append(f.sm, sm)
	return
}

func (f *fsm) AddState(id string, enter, leave Action) (ret State, err error) {
	if _, ok := f.states[id]; ok {
		return ret, fmt.Errorf("State id %s is in use.", id)
	}

	ret = newState(id)
	f.states[id] = internalStateData{ret, enter, leave}
	return
}

func (f *fsm) State(id string) (ret State, ok bool) {
	res, ok := f.states[id]
	if ok {
		ret = res.state
	}
	return
}

// register statemaker's transitors
func (f *fsm) registerStateMapTransitors() error {
	for _, s := range f.sm {
		for _, t := range s.Transitors() {
			if t.IsHidden {
				continue
			}
			st, ok := f.State(t.State)
			if !ok {
				return ErrStateNotFound
			}
			switch {
			case t.IsFallback:
				st.RegisterFallback(t.Transitor)
			case t.Command != "" && t.Type == telegram.TEXT:
				st.RegisterCommand(t.Command, t.Transitor)
			default:
				st.Register(t.Type, t.Transitor)
			}
		}
	}
	f.sm = []StateMaker{}
	return nil
}

func (f *fsm) Start(timeout int) error {
	if err := f.registerStateMapTransitors(); err != nil {
		return err
	}

	// start message manager
	go f.manager.Run(f.api, timeout)

	// start worker goroutines
	return f.Resume()
}

func (f *fsm) Resume() error {
	for err := range f.error_chan {
		if err != nil {
			return err
		}
		go func() { f.error_chan <- f.work() }()
	}
	return nil
}

func (f *fsm) work() (err error) {
	msg := f.manager.Begin()
	defer f.manager.Rollback(msg)

	user := f.userExtractor(msg)
	sid, data, err := f.storage.Load(user.Id)
	if err != nil {
		return
	}

	cur_node, ok := f.states[sid]
	if !ok {
		return fmt.Errorf("Cannot load state[%s] of user#%d", sid, user.Id)
	}
	cur := cur_node.state.clone(user)
	cur.SetData(data)

	do_next := func(cur State, msg *telegram.Message) (next State, err error) {
		next_sid, err := cur.test(msg)
		if err != nil {
			return
		}

		return f.transit(msg, cur, next_sid)
	}

	next, err := do_next(cur, msg)
	if err != nil {
		return
	}

	for ;next.re(); {
		if next, err = do_next(next, msg); err != nil {
			return
		}
	}

	f.manager.Commit(msg)
	return
}

func (f *fsm) transit(msg *telegram.Message, current State, id string) (next State, err error) {
	user := current.User()
	cur_node, ok := f.states[current.Id()]
	if !ok {
		return next, fmt.Errorf("Cannot load state[%s] of user#%d", current.Id(), user.Id)
	}

	next_node, ok := f.states[id]
	if !ok {
		return next, fmt.Errorf("Cannot load next state[%s] of user#%d", id, user.Id)
	}
	next = next_node.state.clone(user)

	if cur_node.leave != nil {
		if err = cur_node.leave(msg, current, f.api); err != nil {
			return
		}
	}

	next.SetData(current.Data())

	if next_node.enter != nil {
		if err = next_node.enter(msg, next, f.api); err != nil {
			return
		}
	}

	err = f.storage.Save(user.Id, next.Id(), next.Data())
	if next.next() != nil {
		next, err = f.transit(msg, next, *next.next())
	}

	if next.re() {

	}

	return
}
