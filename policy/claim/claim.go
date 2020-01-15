package claim

import "github.com/arxeme/policymgt/fsm"

// Claim state consts
type stateEnums struct {
	Initialized             int `state:"1"`  // Initial status
	Processing              int `state:"2"`  // Pending at operation team, to process the claim
	Pending                 int `state:"3"`  // Pending at customer, to submit requested documents
	Rejected                int `state:"4"`  // Claim is rejected by operation team
	Approved                int `state:"5"`  // Claim is approved by operation team, will start reimbursement
	ReimbursementPending    int `state:"11"` // Pending required information for reimbursement from customer
	ReimbursementProcessing int `state:"12"` // Pending reimbursement at finance team
	ReimbursementFailed     int `state:"13"` // Failed in processing the reimbursement
	Reimbursed              int `state:"14"` // Reimbursement is done successfully
}

// State - Claim state enum
var State stateEnums

// Claim transition consts
type transitionEnums struct {
	File                 int `transition:"1"`
	ProofRequest         int `transition:"2"`
	ProofSubmit          int `transition:"3"`
	ProofExpire          int `transition:"4"`
	Approve              int `transition:"5"`
	Reject               int `transition:"6"`
	ReimbursementStart   int `transition:"11"`
	ReimbursementRequest int `transition:"12"`
	ReimbursementFail    int `transition:"13"`
	ReimbursementSucceed int `transition:"14"`
	ReimbursementSolve   int `transition:"15"`
}

// Transition - Claim transition enum
var Transition transitionEnums

var controller *fsm.Controller

// Init - Construct the FSM
func Init() {
	controller = fsm.NewController(State, Transition)

	controller.AddTransition(State.Initialized, State.Processing, Transition.File)
	controller.AddTransition(State.Processing, State.Pending, Transition.ProofRequest)
	controller.AddTransition(State.Pending, State.Processing, Transition.ProofSubmit)
	controller.AddTransition(State.Processing, State.Approved, Transition.Approve)
	controller.AddTransition(State.Processing, State.Rejected, Transition.Reject)
	controller.AddTransition(State.Approved, State.ReimbursementPending, Transition.ReimbursementStart)
	controller.AddTransition(State.ReimbursementPending, State.ReimbursementProcessing, Transition.ReimbursementRequest)
	controller.AddTransition(State.ReimbursementProcessing, State.ReimbursementFailed, Transition.ReimbursementFail)
	controller.AddTransition(State.ReimbursementProcessing, State.Reimbursed, Transition.ReimbursementSucceed)
	controller.AddTransition(State.ReimbursementFailed, State.Reimbursed, Transition.ReimbursementSolve)
}

// Transit - Transit claim state
func Transit(c *Claim, tsn int) {
	controller.Transit(*c, tsn)
}

// Claim - structure represent premium
type Claim struct {
	state int
}

// GetState - Implenmentation of interface fsm.Statable.GetState()
func (c Claim) GetState() int {
	return c.state
}

// SetState - Implenmentation of interface fsm.Statable.SetState()
func (c Claim) SetState(s int) {
	c.state = s
}
