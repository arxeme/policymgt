package claim

import "github.com/arxeme/policymgt/fsm"

// State - Claim state enum
type State int

// Claim state consts
const (
	Initialized             State = 1  // Initial status
	Processing              State = 2  // Pending at operation team, to process the claim
	Pending                 State = 3  // Pending at customer, to submit requested documents
	Rejected                State = 4  // Claim is rejected by operation team
	Approved                State = 5  // Claim is approved by operation team, will start reimbursement
	ReimbursementPending    State = 11 // Pending required information for reimbursement from customer
	ReimbursementProcessing State = 12 // Pending reimbursement at finance team
	ReimbursementFailed     State = 13 // Failed in processing the reimbursement
	Reimbursed              State = 14 // Reimbursement is done successfully
)

// Transition - Claim transition enum
type Transition int

// Claim transition consts
const (
	File                 Transition = 1
	ProofRequest         Transition = 2
	ProofSubmit          Transition = 3
	ProofExpire          Transition = 4
	Approve              Transition = 5
	Reject               Transition = 6
	ReimbursementStart   Transition = 11
	ReimbursementRequest Transition = 12
	ReimbursementFail    Transition = 13
	ReimbursementSucceed Transition = 14
	ReimbursementSolve   Transition = 15
)

var controller *fsm.Controller

// Init - Construct the FSM
func Init() {
	controller = fsm.NewController()

	addTransition(Initialized, Processing, File)
	addTransition(Processing, Pending, ProofRequest)
	addTransition(Pending, Processing, ProofSubmit)
	addTransition(Processing, Approved, Approve)
	addTransition(Processing, Rejected, Reject)
	addTransition(Approved, ReimbursementPending, ReimbursementStart)
	addTransition(ReimbursementPending, ReimbursementProcessing, ReimbursementRequest)
	addTransition(ReimbursementProcessing, ReimbursementFailed, ReimbursementFail)
	addTransition(ReimbursementProcessing, Reimbursed, ReimbursementSucceed)
	addTransition(ReimbursementFailed, Reimbursed, ReimbursementSolve)
}

// Transit - Transit claim state
func Transit(c *Claim, tsn Transition) {
	controller.Transit(*c, int(tsn))
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

// Claim - structure represent premium
type Claim struct {
	state State
}

// GetState - Implenmentation of interface fsm.Statable.GetState()
func (c Claim) GetState() int {
	return int(c.state)
}

// SetState - Implenmentation of interface fsm.Statable.SetState()
func (c Claim) SetState(s int) {
	c.state = State(s)
}
