package gofunopter

// Each of the specific settings may not have things that make sense
// (for example, the location shouldn't have a tolerance) but it's way
// easier to just have one for each type
// Specific types can overwrite various parts if desired

// A float type which can check tolerances, initialize and save history
type OptFloat struct {
	curr     float64   // current value of the float
	init     float64   // initial value of the float
	hist     []float64 // stored history of the float
	AbsTol   float64   // Tolerance on that value
	SaveHist bool      // Do you want to save the history. Called save because maybe we want to have a save to file in the future
	RelTol   float64   // Tolerance relative to the initial value
	Name     string
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
		return Name + "absolute tolerance reached"
	}
	if c.curr/c.init < c.RelTol {
		return Name + " relative tolerance reached"
	}
	return ""
}
