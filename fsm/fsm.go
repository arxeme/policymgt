package fsm

import (
	"errors"
	"fmt"
)

// Undefined - default value for state/transition, should define from 1
const Undefined int = 0

// Statable - Struct with state
type Statable interface {
	GetState() int
	SetState(int)
}

// EventFunc - Function signature of an event
type EventFunc func(Statable) error

// Event - Functions to check before a transition can happen
type Event struct {
	do         EventFunc
	srouce     int
	transition int
}

// NewEvent - Function to create new Event.
// Optional parameters: 1st one is Transition, 2nd is srouce state
func NewEvent(f EventFunc, ops ...int) *Event {
	var e Event
	switch len(ops) {
	case 0:
		e = Event{do: f, srouce: Undefined, transition: Undefined}
	case 1:
		e = Event{do: f, srouce: Undefined, transition: ops[0]}
	case 2:
		e = Event{do: f, srouce: ops[1], transition: ops[0]}
	default:
		panic("only up to 2 optional parameters are accepted")
	}
	return &e
}

type target struct {
	state         int
	prerequisites [](*Event)
	triggers      [](*Event)
}

// Controller - controller of a FSM
type Controller struct {
	states      map[int](*target)
	transitions map[int](map[int](*target))
}

// NewController - Create a controller
func NewController() *Controller {
	c := Controller{transitions: make(map[int](map[int](*target)))}
	return &c
}

// AddTransition - Create a transition in the FSM
func (c *Controller) AddTransition(src, dst, tsn int) error {
	// check if transition exists
	var ts map[int](*target)
	var exist bool
	if ts, exist = c.transitions[src]; !exist {
		ts = make(map[int](*target))
		c.transitions[src] = ts
	}
	if _, exist := ts[tsn]; exist {
		return errors.New("transision already exists")
	}
	// check if states exist
	if _, exist := c.states[src]; !exist {
		c.states[src] = &target{state: src}
	}
	if ptgt, exist := c.states[dst]; exist {
		ts[tsn] = ptgt
	} else {
		ptgt = &target{state: dst}
		ts[tsn] = ptgt
		c.states[dst] = ptgt
	}
	return nil
}

// AddPrerequisite - Add check before transition
func (c *Controller) AddPrerequisite(state int, e *Event) error {
	if tgt, exist := c.states[state]; exist {
		tgt.prerequisites = append(tgt.prerequisites, e)
	}
	return errors.New("target state is undefined")
}

// AddTrigger - Add triggered function after transition
func (c *Controller) AddTrigger(state int, e *Event) error {
	if tgt, exist := c.states[state]; exist {
		tgt.triggers = append(tgt.triggers, e)
	}
	return errors.New("target state is undefined")
}

// Transit - Trigger transition tsn on object s
func (c *Controller) Transit(s Statable, tsn int) (int, error) {
	src := s.GetState()
	if ts, exist := c.transitions[src]; exist {
		if tgt, exist := ts[tsn]; exist {
			// check prerequisits
			for _, pr := range tgt.prerequisites {
				// if source / transition match, then do the check
				if (pr.srouce == Undefined || pr.srouce == src) &&
					(pr.transition == Undefined || pr.transition == tsn) {
					if err := pr.do(s); err != nil {
						return src, fmt.Errorf("%s(%s)",
							"transition failed in prerequisite check", err.Error())
					}
				}
			}
			// transit
			s.SetState(tgt.state)
			// triggers
			for _, tg := range tgt.triggers {
				if (tg.srouce == Undefined || tg.srouce == src) &&
					(tg.transition == Undefined || tg.transition == tsn) {
					tg.do(s)
				}
			}
			return tgt.state, nil
		}
	}
	return src, errors.New("transition is not available for current state")
}
