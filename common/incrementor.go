package common

import (
	"gofunopter/convergence"
	"gofunopter/display"
)

// A Incrementor is a type for incrementing a value. Has methods for testing
// if the value is over a set max and for displaying the value.
type Incrementor struct {
	max   int // Maximum allowable value of the Incrementor
	curr  int // current value of the Incrementor
	total int // Total number at the end of the optimization run
	conv  convergence.C
	name  string
	disp  bool
}

// No Incrementor iterate method because we don't know how many we want to add on each iteration

func NewIncrementor(name string, max int, conv convergence.C, disp bool) *Incrementor {
	return &Incrementor{
		name: name,
		max:  max,
		conv: conv,
		disp: disp,
	}
}

// Add increments the value delta to the Incrementor
func (i *Incrementor) Add(delta int) {
	i.curr += delta
}

func (i *Incrementor) AddToDisplay(d []*display.Struct) []*display.Struct {
	if i.disp {
		d = append(d, &display.Struct{Value: i.curr, Heading: i.name})
	}
	return d
}

// Converged checks if the value of the Incrementor is greater
// than the maximum set value
func (i *Incrementor) Converged() convergence.C {
	if i.curr >= i.max {
		return i.conv
	}
	return nil
}

// Curr returns the current value of the Incrementor
func (i *Incrementor) Curr() int {
	return i.curr
}

// Disp returns true if the value will be displayed during optimization
func (i *Incrementor) Disp() bool {
	return i.disp
}

// SetDisp is a toggle to display the value of the Incrementor during optimization
func (i *Incrementor) SetDisp(val bool) {
	i.disp = val
}

func (i *Incrementor) Initiailize(val int) {
	i.curr = val
}

// Max returns the current setting for the  maximum value of the Incrementor
func (i *Incrementor) Max() int {
	return i.max
}

// SetMax sets the maximum value of the Incrementor
func (i *Incrementor) SetMax(val int) {
	i.max = val
}

// Opt returns the value of the Incrementor at the end of the optimization
func (i *Incrementor) Opt() int {
	return i.total
}

// SetResult sets the result of the Incrementor at the end of the optimization
func (i *Incrementor) SetResult() {
	i.total = i.curr
}

// Incrementor doesn't need to implement Reset because there are no settings
