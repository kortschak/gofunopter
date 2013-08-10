// Package common defines a OptCommon Optimizer type for expressing
// basic tolerances

package common

import (
	"github.com/btracey/gofunopter/common/display"
	"github.com/btracey/gofunopter/common/status"
	"math"
	"time"
)

// OptCommon defines a structure common to all optimizers
// It includes a type to monitor the iterations, a type to
// monitor the function evaluations, and a type to monitor the runtime
type OptCommon struct {
	Iter     *Iterations
	FunEvals *FunctionEvaluations
	Time     *Time
	*display.Display
	stat status.Status
}

// CommonSettings is a list of settings for the OptCommon structure
// See NewCommonSettings for a list of default values
type CommonSettings struct {
	*display.DisplaySettings
	MaximumIterations          int           // Sets the maximum number of iterations that can occur
	MaximumFunctionEvaluations int           // Sets the maximum number of function evaluations that can occur
	MaxRuntime                 time.Duration // Sets the maximum runtime that can elapse
	DisplayIterations          bool          // A toggle if the iteration number should display during the optimization
	DisplayFunctionEvaluations bool          // A toggle if the function evaluations should display during the optimization
	DisplayRuntime             bool          // A toggle if the runtime should display during the optimization
}

// NewCommonSettings creates the default common settings structure
func NewCommonSettings() *CommonSettings {
	return &CommonSettings{
		DisplaySettings:            display.NewDisplaySettings(),
		MaximumIterations:          math.MaxInt32 - 1,                // Defaults to no maximum iterations
		MaximumFunctionEvaluations: math.MaxInt32 - 1,                // Defaults to no maximum function evaluations
		MaxRuntime:                 time.Duration(math.MaxInt64 - 1), // Defaults to no maximum runtime
		DisplayIterations:          true,
		DisplayFunctionEvaluations: true,
		DisplayRuntime:             false,
	}
}

// NewOptCommon creates a new OptCommon structure. Should be called by optimization method
func NewOptCommon() *OptCommon {
	c := &OptCommon{
		Iter:     NewIterations(),
		FunEvals: NewFunctionEvaluations(),
		Time:     NewTime(),
		Display:  display.NewDisplay(),
	}
	return c
}

// SetSettings takes the settings from CommonSettings and translates them
// into the relevant data types
func (c *OptCommon) SetSettings(s *CommonSettings) {

	c.Display.SetSettings(s.DisplaySettings)
	c.Iter.SetMax(s.MaximumIterations)
	c.FunEvals.SetMax(s.MaximumFunctionEvaluations)
	c.Time.SetMax(s.MaxRuntime)
	c.Iter.SetDisp(s.DisplayIterations)
	c.FunEvals.SetDisp(s.DisplayFunctionEvaluations)
	c.Time.SetDisp(s.DisplayRuntime)
}

// CommonResult is a list of results from the common structure
type CommonResult struct {
	Iterations          int           // Total number of iterations taken by the optimizer
	FunctionEvaluations int           // Total number of function evaluations taken by the optimizer
	Runtime             time.Duration // Total runtime elapsed during the optimization
	Status              status.Status // How did the optimizer end
}

// AddToDisplay adds the structures to the display
func (c *OptCommon) AddToDisplay(d []*display.Struct) []*display.Struct {
	//return append(d, &display.Struct{Value: c.curr, Heading: c.name})
	// AddToDisplay can't change because it needs to satisfy displayer interface
	d = c.Iter.AddToDisplay(d)
	d = c.FunEvals.AddToDisplay(d)
	d = c.Time.AddToDisplay(d)
	return d
}

// GetOptCommon is to allow optimizers to easily
// satisfy the Optimizer interface
func (c *OptCommon) GetOptCommon() *OptCommon {
	return c
}

// CommonConverged checks if any of the elements of common have converged
func (c *OptCommon) CommonStatus() status.Status {
	return status.CheckStatus(c.Iter, c.FunEvals, c.Time)
}

// CommonInitialize initializes the elements of the common structure at the start of a run.
func (c *OptCommon) CommonInitialize() {
	c.FunEvals.Initialize()
	c.Iter.Initialize()
	c.Time.Initialize()
}

func (c *OptCommon) SetResult(s status.Status) {
	c.Iter.SetResult()
	c.FunEvals.SetResult()
	c.Time.SetResult()
	c.stat = s
}

// CommonResult sets the results from the
func (c *OptCommon) CommonResult() *CommonResult {
	r := &CommonResult{
		Iterations:          c.Iter.Opt(),
		FunctionEvaluations: c.FunEvals.Opt(),
		Runtime:             c.Time.Opt(),
		Status:              c.stat,
	}
	return r
}

type Iterations struct {
	*Incrementor
}

func NewIterations() *Iterations {
	return &Iterations{
		Incrementor: NewIncrementor("Iter", math.MaxInt32-1, status.MaximumIterations, true),
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
		Incrementor: NewIncrementor("FunEval", math.MaxInt32-1, status.MaximumFunctionEvaluations, true),
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

// Converged returns a status if the elapsed run time is
// longer than the maximum allowed
func (t *Time) Status() status.Status {
	if time.Since(t.init) > t.max {
		return status.MaximumRuntime
	}
	return status.Continue
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
