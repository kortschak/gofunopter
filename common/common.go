// Package common defines a OptCommon Optimizer type for expressing
// basic tolerances

package common

import (
	"gofunopter/convergence"
	"gofunopter/display"
	"math"
	"time"
)

// OptCommon defines a structure common to all optimizers
// All optimizers should embed OptCommon
// MaxIter is the maximum number of allowable iterations.
// MaxRuntime is the maximum runtime of the optimizer
// MaxFunEvals is the maximum allowable number of function calls
// (enforced on a per-iteration) basis, so may go over this limit
// depending on the optimizer used
type OptCommon struct {
	iter     *Iterations
	funEvals *FunctionEvaluations
	time     *Time
	disp     bool
}

func NewOptCommon() *OptCommon {
	c := &OptCommon{
		iter:     NewIterations(),
		funEvals: NewFunctionEvaluations(),
		time:     NewTime(),
		disp:     true,
	}
	return c
}

func (c *OptCommon) AddToDisplay(d []*display.Struct) []*display.Struct {
	//return append(d, &display.Struct{Value: c.curr, Heading: c.name})
	d = c.iter.AddToDisplay(d)
	d = c.funEvals.AddToDisplay(d)
	d = c.time.AddToDisplay(d)
	return d
}

// OptCommon is to allow optimizers to easily
// satisfy the Optimizer interface
func (c *OptCommon) GetOptCommon() *OptCommon {
	return c
}

// Converged checks if any of the elements of common have converged
func (c *OptCommon) Converged() convergence.C {
	return convergence.CheckConvergence(c.iter, c.funEvals, c.time)
}

func (c *OptCommon) Disp() bool {
	return c.disp
}

func (c *OptCommon) SetDisp(b bool) {
	c.disp = b
}

// FunEvals is to allow access to the FunEvals struct
func (c *OptCommon) FunEvals() *FunctionEvaluations {
	return c.funEvals
}

// Iter is to allow access to the Iterations struct
func (c *OptCommon) Iter() *Iterations {
	return c.iter
}

func (c *OptCommon) Time() *Time {
	return c.time
}

// Initialize the common structure at the start of a run.
func (c *OptCommon) Initialize() {
	c.time.Initialize()
	c.funEvals.Initialize()
	c.iter.Initialize()
}

func (c *OptCommon) SetResult() {
	//SetResults(c.iter, c.funEvals, c.runtime)
	c.iter.SetResult()
	c.funEvals.SetResult()
	c.time.SetResult()
}

type Iterations struct {
	*Incrementor
}

func NewIterations() *Iterations {
	return &Iterations{
		Incrementor: NewIncrementor("Iter", math.MaxInt32-1, convergence.Iterations, true),
	}
}

func (i *Iterations) Initialize() error {
	i.Incrementor.curr = 0
	return nil
}

func (i *Iterations) Iterate() error {
	i.Add(1)
	return nil
}

type FunctionEvaluations struct {
	*Incrementor
}

func NewFunctionEvaluations() *FunctionEvaluations {
	return &FunctionEvaluations{
		Incrementor: NewIncrementor("FunEvals", math.MaxInt32-1, convergence.FunEvals, true),
	}
}

func (f *FunctionEvaluations) Initialize() error {
	f.Incrementor.curr = 0
	return nil
}

// Time controls the runtime of the optimizer
type Time struct {
	max  time.Duration
	init time.Time
	opt  time.Duration
	Name string
	disp bool
}

func NewTime() *Time {
	return &Time{
		max: math.MaxInt64,
		// display defaults to off
	}
}

func (t *Time) AddToDisplay(d []*display.Struct) []*display.Struct {
	if t.disp {
		d = append(d, &display.Struct{
			Value:   time.Since(t.init),
			Heading: "Time",
		})
	}
	return d
}

// Converged returns a convergence if the elapsed run time is
// longer than the maximum allowed
func (t *Time) Converged() convergence.C {
	if time.Since(t.init) > t.max {
		return convergence.Time
	}
	return nil
}

func (t *Time) Disp() bool {
	return t.disp
}

func (t *Time) SetDisp(val bool) {
	t.disp = val
}

// Returns the initial value of the start time
func (t *Time) Init() time.Time {
	return t.init
}

// Initialize sets the time of the start of the optimization
func (t *Time) Initialize() error {
	t.init = time.Now()
	return nil
}

// Returns the maximum allowed elapsed runtime of the optimization
func (t *Time) Max() time.Duration {
	return t.max
}

// SetMax sets the maximum allowable elapsed time
// for the optimization. This is only checked between
// iterations
func (t *Time) SetMax(d time.Duration) {
	t.max = d
}

// Opt returns the total elapsed time of the optimization
func (t *Time) Opt() time.Duration {
	return t.opt
}

// SetResult sets the total elapsed time at the end of
// the optimization
func (t *Time) SetResult() {
	t.opt = time.Since(t.init)
}

/*

// TODO: Replace all of the Disps with a disper

*/
