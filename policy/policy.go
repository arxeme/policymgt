package policy

import (
	"errors"
	"fmt"
	"reflect"
	"time"
)

// State enum
type State int

// Policy state consts
const (
	Initialized State = 0
	Paid        State = 1
	Activated   State = 3
	Ready       State = 4
	Active      State = 5
	Expired     State = 6
	Terminated  State = 7
	Cancelling  State = 8
	Cancelled   State = 9
)

// Policy - basic info
type Policy struct {
	ID           uint64
	CreatedAt    time.Time
	PlanID       string
	ActivationID uint64
	Premium      uint64
	PaymentID    uint64
	StartAt      time.Time
	ExpireAt     time.Time

	state State
}

func load(policyID uint64) *Policy {
	return &Policy{ID: policyID}
}

///////////////////////////////////////////////////////////////

type targetState struct {
	state         State
	prerequisites [](func() bool)
	entered       [](func())
}

func transitionTag(t Transition) string {
	return fmt.Sprint(reflect.ValueOf(t).Pointer())
}

type fsm struct {
	availableTransition map[State](map[string]State)
	enabledTransition   map[string](map[State]targetState)
}

// Init - initialize the FSM
func (f *fsm) Init() {
	f.availableTransition = make(map[State](map[string]State))
	f.availableTransition[Initialized][transitionTag(Pay)] = Paid
	f.availableTransition[Activated][transitionTag(Pay)] = Ready
	f.availableTransition[Initialized][transitionTag(Activate)] = Activated
	f.availableTransition[Paid][transitionTag(Activate)] = Ready
	f.availableTransition[Ready][transitionTag(Start)] = Active

	f.enabledTransition = make(map[string](map[State]targetState))
}

// EnableTransition - enable a transition in FSM
func (f *fsm) EnableTransition(tsn Transition, src, dst State) {
	if tsns, exist := f.availableTransition[src]; !exist {
		panic("transition is not available")
	} else if _, exist := tsns[transitionTag(tsn)]; !exist {
		panic("transition is not available")
	}
	t, enabled := f.enabledTransition[transitionTag(tsn)]
	if enabled {
		if _, exist := t[src]; exist {
			panic("transition already exists")
		} else {
			t[src] = targetState{state: dst}
		}
	} else {
		tgts := make(map[State]targetState)
		tgts[src] = targetState{state: dst}
		f.enabledTransition[transitionTag(tsn)] = tgts
		tsn = func(policyID uint64) (State, error) {
			p := load(policyID)
			if _, enabled := tgts[p.state]; !enabled {
				return src, errors.New("transition is not enabled")
			}
			// TODO: check prerequisite
			p.state = dst
			// TODO: trigger entered events
			return p.state, nil
		}
	}
}

// AutoTransition - automatically transit from src to dst
func (f *fsm) AutoTransition(trans Transition, src, dst State) {
}

// AfterEnter - trigger fn after entered state s
func (f *fsm) AfterEnter(s State, fn func(*Policy)) {
}

// BeforeTransit - check prerequisite fn before transit to state s
func (f *fsm) BeforeTransit(s State, fn func(*Policy) bool) {
}

// FSM for policy
var FSM fsm

///////////////////////////////////////////////////////////////

// Create - create a policy under certain product ID
func Create(plan string) *Policy {
	p := &Policy{
		ID:           1, // assign policy ID
		CreatedAt:    time.Now(),
		PlanID:       plan,
		ActivationID: 0,
		Premium:      10, // get from plan
		PaymentID:    0,
		StartAt:      time.Now(), // set the start time
		ExpireAt:     time.Now(), // set the expire time
		state:        Initialized,
	}
	return p
}

// Transition - function signature
type Transition func(policyID uint64) (State, error)

// Pay - Initialized to Paid
//	   - Activated to Ready
var Pay Transition

// Activate - Initialized to Activated
//		    - Paid to Ready
var Activate Transition

// Start - Ready to Active
var Start Transition
