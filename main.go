package main

import "fmt"
import "github.com/arxeme/policymgt/policy"

func main() {
	policy.Initialize()

	demoPolicyCreate()
	demoPolicyPay()
	demoPolicyPurchase()
}

func demoPolicyCreate() {
	// input: planid, userid
	var planid, userid int64
	planid, userid = 123, 456

	p := &policy.Plan{}
	p.LoadByID(planid)
	po := p.CreatePolicy(userid)

	fmt.Println(po)
}

func demoPolicyPay() {
	// input: policyid
}

func demoPolicyPurchase() {
	// input: planid, userid
}
