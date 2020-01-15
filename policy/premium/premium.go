package premium

import "github.com/arxeme/policymgt/fsm"

// State - Premium state enum
type State int

// Premium state consts
const (
	Initialized      State = 1  // Initial status
	Pending          State = 3  // Pending customer action, to pay premium
	Processing       State = 4  // Pending at finance team, to confirm the payment
	Paid             State = 5  // Payment is confirmed
	RefundRequested  State = 11 // Refund request is send by customer
	RefundProcessing State = 12 // Refund is pending processing by finance team
	RefundFailed     State = 13 // Refund failed
	Refunded         State = 14 // Refund successful
)

// Transition - Premium transition enum
type Transition int

// Premium transition consts
const (
	Invoice       Transition = 1
	Pay           Transition = 2
	PaymentFail   Transition = 3
	PaymendDone   Transition = 4
	RefundRequest Transition = 11
	RefundReject  Transition = 12
	RefundApprove Transition = 13
	RefundFail    Transition = 14
	RefundSucceed Transition = 15
	RefundCancel  Transition = 16
	RefundSolve   Transition = 17
)

var controller *fsm.Controller

// Init - Construct the FSM
func Init() {
	controller = fsm.NewController()

	addTransition(Initialized, Pending, Invoice)
	addTransition(Pending, Processing, Pay)
	addTransition(Processing, Pending, PaymentFail)
	addTransition(Processing, Paid, PaymendDone)
	addTransition(Paid, RefundRequested, RefundRequest)
	addTransition(RefundRequested, Paid, RefundReject)
	addTransition(RefundRequested, RefundProcessing, RefundApprove)
	addTransition(RefundProcessing, RefundFailed, RefundFail)
	addTransition(RefundProcessing, Refunded, RefundSucceed)
	addTransition(RefundFailed, Paid, RefundCancel)
	addTransition(RefundFailed, Refunded, RefundSolve)
}

// Transit - Transit claim state
func Transit(a *Premium, tsn Transition) {
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

// Premium - structure represent premium
type Premium struct {
	state State
}

// GetState - Implenmentation of interface fsm.Statable.GetState()
func (p Premium) GetState() int {
	return int(p.state)
}

// SetState - Implenmentation of interface fsm.Statable.SetState()
func (p Premium) SetState(s int) {
	p.state = State(s)
}
