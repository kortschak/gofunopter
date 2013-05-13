package gofunopter

// There is a problem is trying to define all of the different kinds of optimizers.
// This solution decouples the input from the output, so we can build in all the 
// permutations. (Single input, multiple input, boolean input, etc., plus the
// same combinations of outputs and constraints.). The API for this may be large, but
// it's important to save function evaluations whenever possible. Saving function evaluations
// is more important than maximizing ease of use. We can also provide easy translators
// for some of the most common use cases

// Basic problem types

// General pattern is to have Eval() which evaluates everything at a point, and the type should
// cache the result. Obj() then returns that cache result. I'm not sure if we should have
// specific functions for certain parts, such as EvalObj(), EvalConstraints(), and EvalGrad()
// not added now, but possibly in the future when we have a better idea of all the optimizers

// InputFloat is a type that can evaluate the objective at for a single real objective
type InputFloat interface {
	Eval(float64) error
}

// Single input is a type that can evaluate the objective for a multiple input vector
type InputFloatSlice interface {
	Eval([]float64)
}

// Returns the cached single objective value
type OutputFloat interface {
	Obj() float64
}

// Returns the cached multiple objective value
type OutputFloatSlice interface {
	Obj() []float64
}

// Returns the cached gradient value
type GradientFloat interface {
	Grad() float64
}

// Returns the cached gradient value
type GradientFloatSlice interface {
	Grad() []float64
}

type ConstraintFloat interface {
	Constr() float64 // Return the constraint violation
}

type ConstraintFloatSlice interface {
	Constr() []float64 // Returns a vector of all the constraint violations
}

type SISOProblem interface {
	InputFloat
	OutputFloat
}

// Gradient based SISO
type SISOGradBasedProblem interface {
	SISOProblem
	GradientFloat
}

type MISOProblem interface {
	InputFloatSlice
	OutputFloat
}

type MISOGradBasedProblem interface {
	MISOProblem
	GradientFloatSlice
}

/*
// Single input single output
type SISOProblem interface {
	InputFloat
	OutputFloat
}

// Gradient based SISO
type SISOGradBasedProblem interface {
	SISOProblem
	GradientFloat
}
*/
