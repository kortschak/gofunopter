package status

// CheckConvergence checks the convergence of a variadic
// number of converges and returns the first non-nil result
func CheckStatus(cs ...Statuser) Status {
	for _, val := range cs {
		c := val.Status()
		if c != Continue {
			return c
		}
	}
	return Continue
}

// Use type casting for varieties of convergence (grad, etc.)
// use call to convergence for specific convergence test

// A converger is a type that can test if the optimization has converged
type Statuser interface {
	Status() Status
}

// Status is a type for expressing if the optimizer has finished or not
// Zero signifies no convergence or error so the optimizer should continue.
//  Positive values indicate successful convergence
// negative values express failure for some way
type Status int

const (
	Continue Status = iota
	GradAbsTol
	GradRelTol
	ObjAbsTol
	ObjRelTol
	StepAbsTol
	StepRelTol
	WolfeConditionsMet
)

const (
	_                        = iota
	UserFunctionError Status = -1 * iota
	OptimizerError
	Infeasible
	MaximumIterations
	MaximumFunctionEvaluations
	MaximumRuntime
	LinesearchFailure
)
