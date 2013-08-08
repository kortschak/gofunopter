package multi

import (
	"github.com/btracey/gofunopter/display"

	"errors"
	"github.com/gonum/floats"
)

// Floats is the basic framework for a multi-dimensional value for an optimizer
// For users:
// A float can be set to an initial value before the optimization (float.SetInit),
// can be displayed (float.SetDisp and save its history (float.SetSaveHist)
// during the optimization, and be set to a final optimum value
// at the end of the optimization which can be retrieved with float.Opt
// For optimizer writers:
// Provides a method for getting the user-set initial value (float.Init),
// Getting and setting the current value (float.Curr, float.SetCurr), and
// an identifier for identification during display.
// Float adds to the history automatically during the call to float.Curr,
// but if there are additional function evaluations whose values are not
// used in Curr (for example during a linesearch), they can be added with
// AddToHist
// The default is not to save the history.
type Floats struct {
	name     string
	saveHist bool
	curr     []float64
	norm     float64
	normInit float64
	init     []float64
	hist     [][]float64
	disp     bool
	opt      []float64
}

// NewFloat returns a new float with the name given by the first
// input argument, and the display toggle given by the second
// curr is intentionally nil to force optimizers to call init
func NewFloat(name string, disp bool) *Floats {
	return &Floats{
		name: name,
		disp: disp,
	}
}

// AddToHist adds a copy of the variable to the history if SetHist() is true
func (f *Floats) AddToHist(val []float64) {
	if f.saveHist {
		c := make([]float64, len(val))
		copy(c, val)
		f.hist = append(f.hist, c)
	}
}

// Curr returns the current value of floats
func (f *Floats) Curr() []float64 {
	return f.curr
}

// SetCurr sets the current value of the float and appends it to the history
// if float.SaveHist() is true. Assumes that the length does not change per
// iteration. A copy is NOT made. Should only be called by the optimizer.
func (f *Floats) SetCurr(val []float64) {
	f.AddToHist(val)
	copy(f.curr, val)
	f.norm = floats.Norm(f.curr, 2)
}

// Disp returns true if the Float will be displayed during
// optimization and false otherwise (can be overridden at the
// optimizer level)
func (f *Floats) Disp() bool {
	return f.disp
}

// SetDisp() sets the display toggle. If true, the value of this
// variable will be displayed during the optimization (can be
// overridden at the optimizer level)
func (f *Floats) SetDisp(val bool) {
	f.disp = val
}

// Display returns the display struct to be displayed if
// Float.Disp() is true.
func (f *Floats) AddToDisplay(d []*display.Struct) []*display.Struct {
	if f.disp {
		d = append(d, &display.Struct{Value: f.norm, Heading: f.name})
	}
	return d
}

// Hist returns the history of the value over the course of the optimization
// It returns the pointer to the history rather than a copy of the history
// Advanced users can use this to reduce memory cost (by writing the value
// to disk and then reslicing, for example)
func (f *Floats) Hist() [][]float64 {
	return f.hist
}

// SaveHist returns true if all of the values of this variable will be
// stored during optimization
func (f *Floats) SaveHist() bool {
	return f.saveHist
}

// SetSaveHist sets the history saver toggle. If true, all of the values
// of the variable will be stored during optimization.
func (f *Floats) SetSaveHist(val bool) {
	f.saveHist = val
}

// Init returns the initial value of the variable at the start
// of the optimization. A copy is NOT made, so do not modify
// after the start of the optimization
func (f *Floats) Init() []float64 {
	return f.init
}

// Init sets the initial value of the variable to be
// used by the optimizer. (for example, initial location, initial
// function value, etc.)
func (f *Floats) SetInit(val []float64) {
	if len(f.init) > len(val) {
		f.init = f.init[:len(val)]
	} else {
		f.init = make([]float64, len(val))
	}
	copy(f.init, val)
	f.norm = floats.Norm(val, 2)
}

// Initialize initializes the Float to be ready to optimize by
// setting the history slice to have length zero, and setting
// the current value equal to the initial value
// This should be called by the optimizer at the beginning of
// the optimization
func (f *Floats) Initialize() error {
	f.hist = f.hist[:0]

	if f.init == nil {
		return errors.New("multivariate: inial slice is nil")
	}

	f.curr = make([]float64, len(f.init))
	copy(f.curr, f.init)
	f.opt = nil
	f.normInit = floats.Norm(f.curr, 2)
	return nil
}

// Name returns the identifier given to the float during the call to NewFloat
func (f *Floats) Name() string {
	return f.name
}

func (f *Floats) Norm() float64 {
	return f.norm
}

// Opt gets the optimimum value at the conclusion of optimization.
// Returns the true value (not a copy)
func (f *Floats) Opt() []float64 {
	return f.opt
}

// SetResult sets the result and resets for another optimization, setting
// the initial value to the input argument
// This should be called by the optimizer at the end of optimization
func (f *Floats) SetResult() {
	f.init = nil // can't know how big it is ahead of time
	f.opt = f.curr
	// Current is set to NaN to ensure the optimizer calls initialize
	f.curr = nil
}
