package univariate

import (
	"errors"
	"github.com/btracey/gofunopter/convergence"
	"github.com/btracey/gofunopter/display"
	"math"
)

// Bounded step is the basic type for a step from the current
// location with lower and upper boundb. AbsTol tests the gap
// between the upper and lower bounds
type BoundedStep struct {
	*Float
	*convergence.Abs

	lb      float64
	ub      float64
	gap     float64
	initGap float64
}

func NewBoundedStep() *BoundedStep {
	b := &BoundedStep{
		Float: NewFloat("step", false),
		Abs:   convergence.NewAbs(convergence.DefaultStepAbsTol, convergence.StepAbsTol),
		lb:    0,
		ub:    math.Inf(1),
	}
	b.SetInit(1)
	return b
}

func (b *BoundedStep) Initialize() error {
	b.Float.Initialize()
	b.initGap = b.ub - b.lb
	b.gap = b.ub - b.lb
	if b.initGap < 0 {
		return errors.New("bounded step: lower bound is greater than upper bound")
	}
	return nil
}

func (b *BoundedStep) Lb() float64 {
	return b.lb
}

func (b *BoundedStep) Ub() float64 {
	return b.ub
}

func (b *BoundedStep) SetLb(val float64) {
	b.lb = val
	b.gap = b.ub - b.lb
}

func (b *BoundedStep) SetUb(val float64) {
	b.ub = val
	b.gap = b.ub - b.lb
}

// Gradient should display the gap and not the actual value
func (b *BoundedStep) AddToDisplay(d []*display.Struct) []*display.Struct {
	if b.disp {
		d = append(d, &display.Struct{Value: b.ub, Heading: b.name + "UB"},
			&display.Struct{Value: b.ub, Heading: b.name + "UB"})
	}
	return d
}

func (b *BoundedStep) SetResult() {
	b.Float.SetResult()
	b.Float.SetInit(1)
	//b.SetInit(1)
	b.ub = math.Inf(1)
	b.lb = 0
}

// Midpoint between the bounds
func (b *BoundedStep) Midpoint() float64 {
	return (b.lb + b.ub) / 2.0
}

// Is the value between the upper and lower bounds
func (b *BoundedStep) WithinBounds(val float64) bool {
	if val < b.lb {
		return false
	}
	if val > b.ub {
		return false
	}
	return true
}

func (b *BoundedStep) Converged() convergence.C {
	return b.Abs.CheckConvergence(b.gap)
}
