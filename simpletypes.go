package gofunopter

import "math"

// Each of the specific settings may not have things that make sense
// (for example, the location shouldn't have a tolerance) but it's way
// easier to just have one for each type
// Specific types can overwrite various parts if desired

// A float type which can check tolerances, initialize and save history
type OptFloat struct {
	curr     float64   // current value of the float
	Init     float64   // initial value of the float
	absinit  float64   // The absolute value of the initial 
	hist     []float64 // stored history of the float
	AbsTol   float64   // Tolerance on that value
	SaveHist bool      // Do you want to save the history. Called save because maybe we want to have a save to file in the future
	RelTol   float64   // Tolerance relative to the initial value
	Name     string
}

// Initializes by setting the current value to the initial value and
// appending it to the history if necessary
func (o *OptFloat) Initialize() {
	o.curr = o.Init
	o.absinit = math.Abs(o.Init)
	if SaveHist {
		o.hist = make([]float64, 1)
		o.hist = append(o.hist, val)
	}

}

// Currently every time this is set it adds to the history, is there
// a desire to change that?
func (o *OptFloat) Set(val float64) {
	o.curr = val
	if SaveHist {
		o.hist = append(o.hist, val)
	}
}

// Get the current value of OptFloat
func (o *OptFloat) Curr() float64 {
	return o.curr
}

// Get the initial value of OptFloat
func (o *OptFloat) Init() float64 {
	return o.init
}

// Returns the hist slice (not a copy of the hist slice)
func (o *OptFloat) Hist() []float64 {
	return o.hist
}

// Make hist return a copy? Have a CopyHist method?

func (o *OptFloat) CheckConvergence() string {
	if c.curr < c.AbsTol {
		return name + " absolute tolerance reached"
	}
	if c.curr/c.absinit < c.RelTol {
		return name + " relative tolerance reached"
	}
	return ""
}

// Returns the default values for an input location
// Locations don't have any tolerances
func DefaultInputFloat() *OptFloat {
	return &OptFloat{
		AbsTol:   math.Inf(-1),
		RelTol:   math.Inf(-1),
		Name:     "loc",
		SaveHist: false,
		Init:     0,
	}
}

// Returns the default values for an objective
// Objectives generally don't have any tolerances
// No idea what the initial function value is
func DefaultObjectiveFloat() *OptFloat {
	return &OptFloat{
		AbsTol:   math.Inf(-1),
		RelTol:   math.Inf(-1),
		Name:     "fun",
		SaveHist: false,
		Init:     math.NaN(),
	}
}

// Returns the default values for the gradient
func DefaultGradientFloat() *OptFloat {
	return &OptFloat{
		AbsTol:   1E-6,
		RelTol:   1E-8,
		Name:     "grad",
		SaveHist: false,
		Init:     math.NaN(),
	}
}

// Returns the default values for a step size
// no default relative tolerance
func DefaultStepFloat() *OptFloat {
	return &OptFloat{
		AbsTol:   1E-6,
		RelTol:   math.Inf(-1),
		Name:     "step",
		SaveHist: false,
		Init:     1,
	}
}
