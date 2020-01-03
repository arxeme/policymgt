package premium

import "time"

// State enum
type State int

// Premium state consts
const (
	Initialized  State = 0
	Pending      State = 1
	Processing   State = 3
	Paid         State = 4
	Refunding    State = 5
	RefundFailed State = 6
	Refunded     State = 7
)

// Premium struct
type Premium struct {
	ID        uint64
	CreatedAt time.Time
	TxnID     uint64
	TxnAt     time.Time
	RefID     string
}
