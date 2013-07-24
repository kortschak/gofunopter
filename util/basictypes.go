package gofunopter

import (
	"fmt"
)

//TODO: Be more careful with resetting and error checking during optimization

var _ = fmt.Println

type BasicConvergence struct {
	Str string
}

func (b BasicConvergence) ConvergenceType() string {
	return b.Str
}

// Counts up and converges if there is a maximum
type Counter struct {
	max   int // Maximum allowable value of the counter
	curr  int // current value of the counter
	total int // Total number at the end of the optimization run
	conv  Convergence
	name  string
	disp  bool
}

// No counter iterate method because we don't know how many we want to add on each iteration

func NewCounter(name string, max int, conv Convergence, disp bool) *Counter {
	return &Counter{
		name: name,
		max:  max,
		conv: conv,
		disp: disp,
	}
}

func (c *Counter) Disp() bool {
	return c.disp
}

func (c *Counter) SetDisp(val bool) {
	c.disp = val
}

func (c *Counter) Max() int {
	return c.max
}

func (c *Counter) SetMax(val int) {
	c.max = val
}

func (c *Counter) Add(delta int) {
	c.curr += delta
}

func (c *Counter) Converged() Convergence {
	if c.curr > c.max {
		return c.conv
	}
	return nil
}

func (c *Counter) SetResult() {
	c.total = c.curr
}

func (c *Counter) Curr() int {
	return c.curr
}

func (c *Counter) Opt() int {
	return c.total
}

func (c *Counter) AppendHeadings(strs []string) []string {
	return append(strs, c.name)
}

func (c *Counter) AppendValues(vals []interface{}) []interface{} {
	return append(vals, c.curr)
}

var StepBoundsAbsTol Convergence = StepConvergence{"Step bounds absolute tolerance reached"}
var StepBoundsRelTol Convergence = StepConvergence{"Step bounds absolute tolerance reached"}
