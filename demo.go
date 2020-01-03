package main

import (
	"time"

	"github.com/arxeme/policymgt/policy"
)

type transitPolicy struct {
	OrderID string
	*policy.Policy
}

func checkStartDate(p *policy.Policy) {
	if time.Now().After(p.StartAt) {
		policy.Start(p.ID)
	}
}

func isActivated(p *policy.Policy) bool {
	return true
}

func initTransitPolicy() {
	policy.FSM.Init()

	///////////////////////////////////////////////////////////
	// config the state machine

	// declare the states & transitions
	policy.FSM.EnableTransition(policy.Pay, policy.Initialized, policy.Paid)
	policy.FSM.AutoTransition(policy.Activate, policy.Paid, policy.Ready)
	// undecleared transition/states will be disabled

	// add trigger
	policy.FSM.AfterEnter(policy.Ready, checkStartDate)
	policy.FSM.BeforeTransit(policy.Ready, isActivated)
}

func main() {
	///////////////////////////////////////////////////////////
	// Business logic is handled outside.
	initTransitPolicy()

	// when need generate the policy:
	p := &transitPolicy{"order_id_123456", policy.Create("Bukalapak Transit")}

	// when we need transit the state machine
	policy.Activate(p.ID)

}
