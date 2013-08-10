package uni

import (
	"github.com/btracey/gofunopter/common/display"
	"github.com/btracey/gofunopter/common/status"
	"math"

	//"fmt"
)

// Gradient is the basic type for represeting a one-dimensional gradient
// Default is display on and a name of "Grad"
// The default is for the initial value to be NaN, which can be set
// by the user with Init.
// The gradient can converge in two ways:
// 1) If the norm of the gradient is less than absTol. Returns status.GradAbsTol
// 2) If ratio of the current norm of the gradient divided by the initial norm
// of the gradient is less than relTol. Returns status.GradRelTol.
// Default is to converge at the norm of the gradient less than 1E-6
// The optimizer should check if the initial value is NaN and call the user-defined
// function at the initial location if necessary.
type Gradient struct {
	*Float
	*status.Abs
	*status.Rel

	absCurr float64
	absInit float64
}

// Disp defaults to off, init value defaults to zero
// Defaults to NaN so that we evaluate at the initial point
// unless set otherwise
// TODO: Make a Reset() function
// TODO: Add in other defaults
func NewGradient() *Gradient {
	g := &Gradient{
		Float: NewFloat("Grad", true),
		Abs:   status.NewAbs(status.DefaultGradAbsTol, status.GradAbsTol),
		Rel:   status.NewRel(0, status.GradRelTol),
	}
	return g
}

// AddToDisplay adds the norm of the gradient
func (g *Gradient) AddToDisplay(d []*display.Struct) []*display.Struct {
	if g.disp {
		d = append(d, &display.Struct{Value: math.Abs(g.Curr()), Heading: "GradNorm"})
	}
	return d
}

// Initialize sets curr = init and sets absInit
func (g *Gradient) Initialize() error {
	g.Float.Initialize()
	g.absInit = math.Abs(g.init)
	g.absCurr = g.absInit
	return nil
}

// SetCurr sets the current value and updates the value norm
func (g *Gradient) SetCurr(val float64) {
	g.Float.SetCurr(val)
	g.absCurr = math.Abs(val)
}

// Converged tests if either the absolute norm or the relative norm have converged
func (g *Gradient) Status() status.Status {
	// Test absolute convergence
	c := g.Abs.Status(g.absCurr)
	if c != status.Continue {
		return c
	}
	// Test relative convergence
	c = g.Rel.Status(g.absCurr, g.absInit)
	if c != status.Continue {
		return c
	}
	return c
}

// SetResult sets the result at the end of the optimaziton (value found in Opt()),
// and resets the initial value to NaN
func (g *Gradient) SetResult() {
	g.Float.SetResult()
}
