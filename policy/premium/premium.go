package premium

import "github.com/arxeme/policymgt/fsm"

type state struct {
	Initialized      int `state:"1"`  // Initial status
	Pending          int `state:"2"`  // Pending customer action, to pay premium
	Processing       int `state:"3"`  // Pending at finance team, to confirm the payment
	Paid             int `state:"4"`  // Payment is confirmed
	RefundRequested  int `state:"11"` // Refund request is send by customer
	RefundProcessing int `state:"12"` // Refund is pending processing by finance team
	RefundFailed     int `state:"13"` // Refund failed
	Refunded         int `state:"14"` // Refund successful
}

type transition struct {
	Invoice       int `transition:"1"`
	Pay           int `transition:"2"`
	PaymentFail   int `transition:"3"`
	PaymendDone   int `transition:"4"`
	RefundRequest int `transition:"11"`
	RefundReject  int `transition:"12"`
	RefundApprove int `transition:"13"`
	RefundFail    int `transition:"14"`
	RefundSucceed int `transition:"15"`
	RefundCancel  int `transition:"16"`
	RefundSolve   int `transition:"17"`
}

// State - Premium state enum
var State state

// Transition - Premium transition enum
var Transition transition

var controller *fsm.Controller

// Init - Construct the FSM
func Init() {
	controller = fsm.NewController(State, Transition)

	controller.AddTransition(State.Initialized, State.Pending, Transition.Invoice)
	controller.AddTransition(State.Pending, State.Processing, Transition.Pay)
	controller.AddTransition(State.Processing, State.Pending, Transition.PaymentFail)
	controller.AddTransition(State.Processing, State.Paid, Transition.PaymendDone)
	controller.AddTransition(State.Paid, State.RefundRequested, Transition.RefundRequest)
	controller.AddTransition(State.RefundRequested, State.Paid, Transition.RefundReject)
	controller.AddTransition(State.RefundRequested, State.RefundProcessing, Transition.RefundApprove)
	controller.AddTransition(State.RefundProcessing, State.RefundFailed, Transition.RefundFail)
	controller.AddTransition(State.RefundProcessing, State.Refunded, Transition.RefundSucceed)
	controller.AddTransition(State.RefundFailed, State.Paid, Transition.RefundCancel)
	controller.AddTransition(State.RefundFailed, State.Refunded, Transition.RefundSolve)
}

// Transit - Transit claim state
func Transit(a *Premium, tsn int) {
	controller.Transit(*a, tsn)
}

// Premium - structure represent premium
type Premium struct {
	state int
}

// GetState - Implenmentation of interface fsm.Statable.GetState()
func (p Premium) GetState() int {
	return p.state
}

// SetState - Implenmentation of interface fsm.Statable.SetState()
func (p Premium) SetState(s int) {
	p.state = s
}
