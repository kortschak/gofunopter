package univariate

/*
// Step is the basic type for representing a step
// AbsTol is on the difference between the
type Step struct {
	*Float
	*AbsTol
	*RelTol

	absCurr float64
	absInit float64
}

// Disp defaults to off, init value defaults to zero
// Defaults to NaN so that we evaluate at the initial point
// unless set otherwise
// TODO: Make a Reset() function
// TODO: Add in other defaults
func NewStep() *Step {
	s := &Step{
		Float:  NewFloat("Step"),
		AbsTol: NewAbsTol(0, convergence.StepAbsTol),
		RelTol: NewRelTol(0, convergence.StepRelTol),
	}
	s.SetInit(1)
	// Disp defaults to off
	return s
}

// AddToDisplay adds the norm of the step
func (s *Step) AddToDisplay(d []*display.Struct) []*display.Struct {
	if s.disp {
		d = append(d, &display.Struct{Value: math.Abs(s.Curr()), Heading: "StepNorm"})
	}
	return d
}

// Initialize sets curr = init and sets absInit
func (s *Step) Initialize() error {
	s.Float.Initialize()
	s.absInit = math.Abs(s.init)
}

// SetCurr sets the current value and updates the value norm
func (s *Step) SetCurr(val float64) {
	s.Float.SetCurr(val)
	s.absCurr = math.Abs(val)
}

// Converged tests if either the absolute norm or the relative norm have converged
func (s *Step) Converged() convergence.C {
	// Test absolute convergence
	c := s.AbsTol.CheckConvergence(s.absCurr)
	if c != nil {
		return c
	}
	// Test relative convergence
	return s.RelTol.CheckConvergence(s.absCurr, s.absInit)
}

// SetResult sets the result at the end of the optimaziton (value found in Opt()),
// and resets the initial value to 1
func (s *Step) SetResult() {
	s.Float.SetResult(1)
}
*/
