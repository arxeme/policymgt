package activation

import "github.com/arxeme/policymgt/fsm"

import "errors"

// State - Activation state enum
type State int

// Activation state consts
const (
	Initialized State = 1 // Initial status
	Pending     State = 3 // Pending at customer, to submit required information
	Processing  State = 2 // Pending at operation team, to process the claim
	Activated   State = 4 // Claim is rejected by operation team
)

// Transition - Activation transition enum
type Transition int

// Activation transition consts
const (
	Start   Transition = 1
	Submit  Transition = 2
	Succeed Transition = 3
	Fail    Transition = 4
)

var controller *fsm.Controller

// Init - Construct the FSM
func Init() {
	controller = fsm.NewController()

	addTransition(Initialized, Pending, Start)
	addTransition(Pending, Processing, Submit)
	addTransition(Processing, Pending, Fail)
	addTransition(Processing, Activated, Succeed)

	addTrigger(Processing, fsm.NewEvent(afterSubmit))
}

// Transit - Transit claim state
func Transit(a *Activation, tsn Transition) {
	controller.Transit(*a, int(tsn))
}

func addTransition(src, dst State, tsn Transition) error {
	return controller.AddTransition(int(src), int(dst), int(tsn))
}

func addPrerequisite(state State, e *fsm.Event) error {
	return controller.AddPrerequisite(int(state), e)
}

func addTrigger(state State, e *fsm.Event) error {
	return controller.AddTrigger(int(state), e)
}

// Activation - structure represent premium
type Activation struct {
	state State
	retry int
}

// GetState - Implenmentation of interface fsm.Statable.GetState()
func (a Activation) GetState() int {
	return int(a.state)
}

// SetState - Implenmentation of interface fsm.Statable.SetState()
func (a Activation) SetState(s int) {
	a.state = State(s)
}

func afterSubmit(s fsm.Statable) error {
	if a, ok := s.(*Activation); ok {
		a.retry++
		return nil
	}
	return errors.New("illegal parameters: type Activation expected")
}
