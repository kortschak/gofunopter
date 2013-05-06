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
	Iter     *Iterations
	FunEvals *Counter
	Runtime  *RuntimeStruct
}

func DefaultCommon() *Common {
	return &Common{Iter: &Iterations{Counter: &Counter{Max: math.MaxInt32, Name: "iterations"}},
		FunEvals: &Counter{Max: math.MaxInt32, Name: "function evaluations"},
		Runtime:  &RuntimeStruct{Max: math.MaxInt64, Name: "runtime"},
	}
}

// Initialize the common structure at the start of a run.
func (c *Common) Initialize() {
	c.Runtime.Initialize()
}

// Check if any of the elements of the common structure have converged
func (c *Common) CheckConvergence() (str string) {
	return CheckConvergence(c.Iter, c.FunEvals, c.Runtime)
}

func (c *Common) Iterate() {
	Iterate(c.Iter)
}

func (c *Common) DisplayHeadings() []string {
	return []string{"Iter", "FunEvals"}
}

func (c *Common) DisplayValues() []string {
	return []string{"Iter", "FunEvals"}
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
		return "Maximum " + r.Name + " reached"
	}
	return ""
}

type Iterations struct {
	*Counter
}

func (i *Iterations) Iterate() {
	i.Add(1)
}
