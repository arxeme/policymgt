package policy

import "github.com/arxeme/policymgt/fsm"

// Policy state consts
type stateEnums struct {
	Initialized int `state:"1"`  // Initial status
	Activated   int `state:"3"`  // Activation completed
	Paid        int `state:"4"`  // Payment for premium completed
	Ready       int `state:"5"`  // Ready to start the protection
	Active      int `state:"11"` // In protection
	Expired     int `state:"21"` // Policy ended after the policy duration
	Terminated  int `state:"22"` // Policy ended before the policy duration because of claims or else
	Cancelling  int `state:"23"` // Request of cancelling the policy is initiated by customer
	Cancelled   int `state:"24"` // Cancellation of the policy is confirmed by operation team
}

// State - Policy state enum
var State stateEnums

// Policy transition consts
type transitionEnums struct {
	Activate      int `transition:"1"`
	Pay           int `transition:"2"`
	Start         int `transition:"3"`
	Expire        int `transition:"11"`
	Terminate     int `transition:"12"`
	CancelRequest int `transition:"13"`
	Cancel        int `transition:"14"`
}

// Transition - Policy transition enum
var Transition transitionEnums

var controller *fsm.Controller

// Init - Construct the FSM
func Init() {
	controller = fsm.NewController(&State, &Transition)

	controller.AddTransition(State.Initialized, State.Activated, Transition.Activate)
	controller.AddTransition(State.Initialized, State.Paid, Transition.Pay)
	controller.AddTransition(State.Activated, State.Ready, Transition.Pay)
	controller.AddTransition(State.Paid, State.Ready, Transition.Activate)
	controller.AddTransition(State.Ready, State.Active, Transition.Start)
	controller.AddTransition(State.Active, State.Expired, Transition.Expire)
	controller.AddTransition(State.Active, State.Terminated, Transition.Terminate)
	controller.AddTransition(State.Paid, State.Cancelling, Transition.CancelRequest)
	controller.AddTransition(State.Ready, State.Cancelling, Transition.CancelRequest)
	controller.AddTransition(State.Cancelling, State.Cancelled, Transition.Cancel)
}

// Transit - Transit claim state
func Transit(a *Policy, tsn int) {
	controller.Transit(*a, tsn)
}

// Policy - structure represent premium
type Policy struct {
	state int
}

// GetState - Implenmentation of interface fsm.Statable.GetState()
func (p Policy) GetState() int {
	return p.state
}

// SetState - Implenmentation of interface fsm.Statable.SetState()
func (p Policy) SetState(s int) {
	p.state = s
}
