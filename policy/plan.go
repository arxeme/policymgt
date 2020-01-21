package policy

import "time"

// Plan - Represnet an insurance plan
type Plan struct {
	ID               int64
	StartAt          time.Time
	ProductID        string
	Description      string
	Restrictions     string
	PremiumPlan      string
	SubscriptionPlan string
	Coverages        string
}

// LoadByID - Load contents of a plan from DB by plan ID
func (p *Plan) LoadByID(id int64) {
	// TODO: load from DB
}

// CreatePolicy - Create a policy under current plan, return policy ID
func (p *Plan) CreatePolicy(userid int64) *Policy {
	po := &Policy{}
	// TODO: genereate policy and save to DB

	Transit(po, Transition.Create)
	return po
}
