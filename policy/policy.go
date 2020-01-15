package policy

import "github.com/arxeme/policymgt/fsm"

// State - Policy state enum
type State int

// Policy state consts
const (
	Initialized State = 1  // Initial status
	Activated   State = 3  // Activation completed
	Paid        State = 4  // Payment for premium completed
	Ready       State = 5  // Ready to start the protection
	Active      State = 11 // In protection
	Expired     State = 21 // Policy ended after the policy duration
	Terminated  State = 22 // Policy ended before the policy duration because of claims or else
	Cancelling  State = 23 // Request of cancelling the policy is initiated by customer
	Cancelled   State = 24 // Cancellation of the policy is confirmed by operation team
)

// Transition - Policy transition enum
type Transition int

// Policy transition consts
const (
	Activate      Transition = 1
	Pay           Transition = 2
	Start         Transition = 3
	Expire        Transition = 11
	Terminate     Transition = 12
	CancelRequest Transition = 13
	Cancel        Transition = 14
)

var controller *fsm.Controller

// Init - Construct the FSM
func Init() {
	controller = fsm.NewController()

	addTransition(Initialized, Activated, Activate)
	addTransition(Initialized, Paid, Pay)
	addTransition(Activated, Ready, Pay)
	addTransition(Paid, Ready, Activate)
	addTransition(Ready, Active, Start)
	addTransition(Active, Expired, Expire)
	addTransition(Active, Terminated, Terminate)
	addTransition(Paid, Cancelling, CancelRequest)
	addTransition(Ready, Cancelling, CancelRequest)
	addTransition(Cancelling, Cancelled, Cancel)
}

// Transit - Transit claim state
func Transit(a *Policy, tsn Transition) {
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

// Policy - structure represent premium
type Policy struct {
	state State
}

// GetState - Implenmentation of interface fsm.Statable.GetState()
func (p Policy) GetState() int {
	return int(p.state)
}

// SetState - Implenmentation of interface fsm.Statable.SetState()
func (p Policy) SetState(s int) {
	p.state = State(s)
}
