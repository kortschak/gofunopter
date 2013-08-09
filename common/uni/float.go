package uni

import (
	"errors"
	"github.com/btracey/gofunopter/common/display"
	"math"
)

// Float is the basic framework for a one-dimensional value for an optimizer
// For users:
// A float can be set to an initial value before the optimization (float.SetInit),
// can be displayed (float.SetDisp and save its history (float.SetSaveHist)
// during the optimization, and be set to a final optimum value
// at the end of the optimization which can be retrieved with float.Opt
// For optimizer writers:
// Provides a method for getting the user-set initial value (float.Init),
// Getting and setting the current value (float.Curr, float.SetCurr), and
// an identifier for identification during display.
// Float does not adds to the history automatically during the call to float.Curr,
// as some times the function may be called but the value will not be set. As a result,
// the optimizer should call AddToHist after performing a function evaluation
// The default is not to save the history.
type Float struct {
	name     string
	saveHist bool
	curr     float64
	init     float64
	hist     []float64
	disp     bool
	opt      float64
}

// NewFloat returns a new float with the name given by the first
// input argument, and the display toggle given by the second
func NewFloat(name string, disp bool) *Float {
	// Current is set to NaN to ensure the optimizer calls initialize
	return &Float{name: name,
		curr: math.NaN(),
		init: math.NaN(),
		disp: disp,
	}
}

// AddToHist adds the variable to the history if SetHist() is true
func (b *Float) AddToHist(val float64) {
	if b.saveHist {
		b.hist = append(b.hist, val)
	}
}

// Curr returns the current value of the float
func (b *Float) Curr() float64 {
	return b.curr
}

// SetCurr sets the current value of the float
func (b *Float) SetCurr(val float64) {
	b.curr = val
}

// Disp returns true if the Float will be displayed during
// optimization and false otherwise (can be overridden at the
// optimizer level)
func (b *Float) Disp() bool {
	return b.disp
}

// SetDisp() sets the display toggle. If true, the value of this
// variable will be displayed during the optimization (can be
// overridden at the optimizer level)
func (b *Float) SetDisp(val bool) {
	b.disp = val
}

// Display returns the display struct to be displayed if
// Float.Disp() is true.
func (b *Float) AddToDisplay(d []*display.Struct) []*display.Struct {
	if b.disp {
		d = append(d, &display.Struct{Value: b.curr, Heading: b.name})
	}
	return d
}

// Hist returns the history of the value over the course of the optimization
// It returns the pointer to the history rather than a copy of the history
// Advanced users can use this to reduce memory cost (by writing the value
// to disk and then reslicing, for example)
func (b *Float) Hist() []float64 {
	return b.hist
}

// SaveHist returns true if all of the values of this variable will be
// stored during optimization
func (b *Float) SaveHist() bool {
	return b.saveHist
}

// SetSaveHist sets the history saver toggle. If true, all of the values
// of the variable will be stored during optimization.
func (b *Float) SetSaveHist(val bool) {
	b.saveHist = val
}

// Init returns the initial value of the variable at the start
// of the optimization
func (b *Float) Init() float64 {
	return b.init
}

// Init sets the initial value of the variable to be
// used by the optimizer. (for example, initial location, initial
// function value, etc.)
func (b *Float) SetInit(val float64) {
	b.init = val
}

// Initialize initializes the Float to be ready to optimize by
// setting the history slice to have length zero, and setting
// the current value equal to the initial value
// This should be called by the optimizer at the beginning of
// the optimization
func (b *Float) Initialize() error {
	b.hist = b.hist[:0]
	if math.IsNaN(b.init) {
		return errors.New("float: Initial value is NaN")
	}
	b.curr = b.init
	b.opt = math.NaN()
	return nil
}

// Name returns the identifier given to the float during the call to NewFloat
func (b *Float) Name() string {
	return b.name
}

// Opt gets the optimimum value at the conclusion of optimization.
func (b *Float) Opt() float64 {
	return b.opt
}

// SetResult sets the result and resets for another optimization, setting
// the initial value to the input argument
// This should be called by the optimizer at the end of optimization
func (b *Float) SetResult() {
	b.opt = b.curr
	// Current is set to NaN to ensure the optimizer calls initialize
	b.init = math.NaN()
	b.curr = math.NaN()
}
