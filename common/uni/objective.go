package uni

import (
	"github.com/btracey/gofunopter/common/status"
	"math"
)

// Objective is the default type for a univariate objective value
// For users:
// It defaults to having a display of On, and a name of "Obj"
// It has two different ways it can converge.
// 1) If the function value
// is below absTol (defaulted to negative infinity), it will converge
// with status.ObjAbsTol.
// 2) If the last change in the function value divided by the first
// change in the function value is less than relTol, it will converge
// with status.ObjRelTol. Default is relTol = 0
// The defaults are so that the default case is to drive the gradient
// to convergence
// For optimizers:
// The initial value is defaulted to NaN. The optimizer should test
// if the initial value has been set by the user, otherwise the
// optimizer should evaluate the function at the initial location.
type Objective struct {
	*Float
	*status.Abs
	*status.Rel

	delta     float64
	initDelta float64 // initial change off of which the delta is based
	relconv   status.Status
	absconv   status.Status
}

// NewObjective returns the default objective structure
func NewObjective() *Objective {
	o := &Objective{
		delta: math.NaN(),
		Float: NewFloat("Obj", true),
		Abs:   status.NewAbs(math.Inf(-1), status.ObjAbsTol),
		Rel:   status.NewRel(0, status.ObjRelTol),
	}
	return o
}

// SetResult sets the optimum value, and resets the initial value to NaN
func (o *Objective) SetResult() {
	o.Float.SetResult()
}

// SetCurr sets a new value  for the current location and updates the
// delta from the last value
func (o *Objective) SetCurr(f float64) {
	// Find the current delta in the values
	o.delta = math.Abs(o.Float.curr - f)
	// Set the initial delta if necessary
	if math.IsNaN(o.initDelta) {
		o.initDelta = o.delta
	}
	// Set the new current value
	o.Float.SetCurr(f)
}

// Converged tests if either AbsTol or RelTol have converged
func (o *Objective) Status() status.Status {
	// Test absolute convergence
	c := o.Abs.Status(o.curr)
	if c != status.Continue {
		return c
	}
	// Test relative convergence
	return o.Rel.Status(o.delta, o.initDelta)
}
