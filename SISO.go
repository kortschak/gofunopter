package gofunopter

// Single input single output

// Not sure if we want to do this. Makes individual implementations
// much easier, but also harder to read. 

/*
type SISOCommon struct {
	Loc *OptFloat // Location
	Obj *OptFloat // Function Value
	*Common
}

func DefaultSISOCommon() *SISOCommon {
	s := &SISOCommon{
		Loc:    DefaultInputFloat(),
		Obj:    DefaultObjectiveFloat(),
		Common: DefaultCommon(),
	}
	SetDisplayMethods(s)
	return s
}

type SISOOptimizer struct {
	SISOCommon
	Fun SISOProblem
}

func DefaultSISOOptimizer() *SISOOptimizer {
	s := &SISOOptimizer{SISOCommon: DefaultSISOCommon()}
	return s
}

type SISOGradBasedOptimizer struct {
	*SISOCommon
	Grad *OptFloat // Gradient value
	Fun  *SISOProblem
}

func DefaultSISOGradBasedOptimizer() *SISOGradBasedOptimizer {
	s := &SISOGradBasedOptimizer{
		SISOCommon: DefaultSISOCommon(),
		Grad:       DefaultGradientFloat(),
	}
	SetDisplayMethods(s)
}

func (s *SISOGradBasedOptimizer) Initialize(Fun *SISOProblem) {
	// Initialize takes all of these in so function evaluations can be saved if 
	// the information is already there
	c.Loc.Initialize()
	if math.IsNaN(f.Curr()) {
		// Initial function value hasn't been set, so do it.
		err = fun.Eval(x.Curr())
		if err != nil {
			return fmt.Errorf("Error evaluating the function at the set initial value %v", x.Curr())
		}
		c.Obj.Init = fun.Obj()
		c.Grad.Init = fun.Grad()
	}
	c.Obj.Initialize()
	c.Grad.Initialize()
}
*/
