package gofunopter

import (
	"math"
	"time"
)

// TODO: Replace all of the Disps with a disper

// Defines a structure common to all optimizers
// MaxIter is the maximum number of allowable iterations.
// MaxRuntime is the maximum runtime of the optimizer
// MaxFunEvals is the maximum allowable number of function calls 
// (enforced on a per-iteration) basis, so may go over this limit
// depending on the optimizer used
type Common struct {
	Iter     *Iterations
	FunEvals *FunctionEvaluations
	Runtime  *RuntimeStruct
	*DisplayStruct
	Displayer
	disp bool
}

func DefaultCommon() *Common {
	c := &Common{
		Iter:          DefaultIterations(),
		FunEvals:      DefaultFunctionEvaluations(),
		Runtime:       DefaultRuntime(),
		DisplayStruct: DefaultDisplayStruct(),
		disp:          true,
	}
	SetDisplayMethods(c)
	return c
}

func (c *Common) Disp() bool {
	return c.disp
}

func (c *Common) SetDisp(val bool) {
	c.disp = val
}

func (c *Common) Common() *Common {
	return c
}

// Initialize the common structure at the start of a run.
func (c *Common) Initialize() {
	c.Runtime.Initialize()
}

// Check if any of the elements of the common structure have converged
func (c *Common) Converged() Convergence {
	return Converged(c.Iter, c.FunEvals, c.Runtime)
}

func (c *Common) Iterate() {
	Iterate(c.Iter, c.DisplayStruct)
}

func (c *Common) AppendHeadings(strs []string) []string {
	return AppendHeadings(strs, c.Runtime, c.Iter, c.FunEvals)
}

func (c *Common) AppendValues(vals []interface{}) []interface{} {
	return AppendValues(vals, c.Runtime, c.Iter, c.FunEvals)
}

func (c *Common) SetResult() {
	SetResults(c.Iter, c.FunEvals, c.Runtime)
}

type RuntimeStruct struct {
	Max  time.Duration
	init time.Time
	opt  time.Duration
	Name string
	disp bool
}

func (r *RuntimeStruct) Disp() bool {
	return r.disp
}

func (r *RuntimeStruct) SetDisp(val bool) {
	r.disp = val
}

func DefaultRuntime() *RuntimeStruct {
	return &RuntimeStruct{
		Max:  math.MaxInt64,
		disp: false,
	}
}

func (r *RuntimeStruct) AppendHeadings(strs []string) []string {
	return append(strs, "Runtime")
}

func (r *RuntimeStruct) AppendValues(vals []interface{}) []interface{} {
	return append(vals, time.Since(r.init))
}

func (r *RuntimeStruct) Start() time.Time {
	return r.init
}

func (r *RuntimeStruct) Initialize() {
	r.init = time.Now()
}

var MaxRuntime Convergence = BasicConvergence{"Maximum runtime reached"}

func (r *RuntimeStruct) Converged() Convergence {
	if time.Since(r.init).Seconds() > r.Max.Seconds() {
		return MaxRuntime
	}
	return nil
}

func (r *RuntimeStruct) SetResult() {
	r.opt = time.Since(r.init)
}

func (r *RuntimeStruct) Opt() time.Duration {
	return r.opt
}

var MaxIter Convergence = BasicConvergence{"Maximum iterations reached"}

type Iterations struct {
	*Counter
}

func (i *Iterations) Iterate() error {
	i.Add(1)
	return nil
}

func DefaultIterations() *Iterations {
	return &Iterations{
		Counter: NewCounter("Iter", math.MaxInt32-1, MaxIter, true),
	}
}

var MaxFunEvals Convergence = BasicConvergence{"Maximum function evaluations reached"}

func DefaultFunctionEvaluations() *FunctionEvaluations {
	return &FunctionEvaluations{
		Counter: NewCounter("FunEvals", math.MaxInt32-1, MaxFunEvals, true),
	}
}

type FunctionEvaluations struct {
	*Counter
}
