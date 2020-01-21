package activation

import "github.com/arxeme/policymgt/fsm"

import "errors"

// State - Activation state enum
var State = &struct {
	Pending    int `state:"1"` // Pending at customer, to submit required information
	Processing int `state:"2"` // Pending at operation team, to process the claim
	Activated  int `state:"4"` // Claim is rejected by operation team
}{}

// Transition - Activation transition enum
var Transition = &struct {
	Create  int `transition:"1"`
	Submit  int `transition:"2"`
	Succeed int `transition:"3"`
	Fail    int `transition:"4"`
}{}

var controller *fsm.Controller

// Initialize - Construct the FSM
func Initialize() {
	controller = fsm.NewController(State, Transition)

	controller.AddTransition(State.Pending, State.Processing, Transition.Submit)
	controller.AddTransition(State.Processing, State.Pending, Transition.Fail)
	controller.AddTransition(State.Processing, State.Activated, Transition.Succeed)

	controller.AddTrigger(State.Processing, fsm.NewEvent(afterSubmit))
}

// Transit - Transit claim state
func Transit(a *Activation, tsn int) {
	controller.Transit(*a, tsn)
}

// Activation - structure represent premium
type Activation struct {
	state int
	retry int
}

// GetState - Implenmentation of interface fsm.Statable.GetState()
func (a Activation) GetState() int {
	return a.state
}

// SetState - Implenmentation of interface fsm.Statable.SetState()
func (a Activation) SetState(s int) {
	a.state = s
}

func afterSubmit(s fsm.Statable) error {
	if a, ok := s.(*Activation); ok {
		a.retry++
		return nil
	}
	return errors.New("illegal parameters: type Activation expected")
}
