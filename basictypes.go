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

type GradConvergence struct{ Str string }

func (b GradConvergence) ConvergenceType() string {
	return b.Str
}

type LocConvergence struct{ Str string }

func (b LocConvergence) ConvergenceType() string {
	return b.Str
}

type FunConvergence struct{ Str string }

func (b FunConvergence) ConvergenceType() string {
	return b.Str
}

type StepConvergence struct{ Str string }

func (b StepConvergence) ConvergenceType() string {
	return b.Str
}

type OptimizeError struct {
	Str string
	Err error
}

func (o OptimizeError) ConvergenceType() string {
	return o.Str + o.Err.Error()
}

func (o OptimizeError) Error() string {
	return o.Str + o.Err.Error()
}

func InitializationError(err error) OptimizeError {
	return OptimizeError{
		Str: "Optimizer failed to initialize, ",
		Err: err,
	}
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

var LocAbsTol Convergence = LocConvergence{"Location absolute tolerance reached"}
var LocRelTol Convergence = LocConvergence{"Location relative tolerance reached"}
var ObjAbsTol Convergence = FunConvergence{"Function absolute tolerance reached"}
var ObjRelTol Convergence = FunConvergence{"Function relative tolerance reached"}
var GradAbsTol Convergence = GradConvergence{"Gradient absolute tolerance reached"}
var GradRelTol Convergence = GradConvergence{"Gradient relative tolerance reached"}
var StepAbsTol Convergence = StepConvergence{"Step absolute tolerance reached"}
var StepRelTol Convergence = StepConvergence{"Step relative tolerance reached"}
var StepBoundsAbsTol Convergence = StepConvergence{"Step bounds absolute tolerance reached"}
var StepBoundsRelTol Convergence = StepConvergence{"Step bounds absolute tolerance reached"}

const DefaultGradAbsTol = 1E-6
const DefaultGradRelTol = 1E-8

const DefaultInitStepSize = 1
const DefaultBoundedStepFloatAbsTol = 0 //
const DefaultBoundedStepFloatRelTol = 0 // Turn off step rel tol
