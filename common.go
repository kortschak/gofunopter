package gofunopter

import (
	"math"
	"time"
)

// Defines a structure common to all optimizers
// MaxIter is the maximum number of allowable iterations.
// MaxRuntime is the maximum runtime of the optimizer
// MaxFunEvals is the maximum allowable number of function calls 
// (enforced on a per-iteration) basis, so may go over this limit
// depending on the optimizer used
type Common struct {
	Iter     *Counter
	FunEvals *Counter
	Runtime  *RuntimeStruct
	Display
}

func DefaultCommon() *Common {
	return &Common{Iter: &Counter{Max: math.MaxInt32, Name: "iterations"},
		FunEvals:   &Counter{Max: math.MaxInt32, Name: "function evaluations"},
		MaxRuntime: &FloatMax{Max: math.Inf(1)},
		Display:    DefaultDisplay(),
	}
}

// Initialize the common structure at the start of a run.
func (c *Common) Initialize() {
	c.Runtime.Initialize()
}

// Check if any of the elements of the common structure have converged
func (c *Common) CheckConvergence() (str string) {
	str = c.Iter.CheckConvergence()
	if str != "" {
		return str
	}
	str = c.FunEvals.CheckConvergence()
	if str != "" {
		return str
	}
	str = c.Runtime.CheckConvergence()
	if str != "" {
		return str
	}
	return ""
}

func (c *Common) Iterate() {
	c.Iter.Add(1)
	c.Display.Iterate()
}

type RuntimeStruct struct {
	Max  time.Duration
	init time.Time
	Name string
}

func (r *RuntimeStruct) Initialize() {
	r.init = time.Now()
}

func (r *RuntimeStruct) CheckConvergence() string {
	if time.Since(r.init).Seconds() > r.Max.Seconds() {
		return "Maximum runtime reached"
	}
}
