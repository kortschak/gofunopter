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

type SisoProblem interface {
	Eval(loc float64) (obj float64, err error)
}

// Gradient based SISO
type SisoGradBasedProblem interface {
	Eval(loc float64) (obj float64, grad float64, err error)
}

type MisoProblem interface {
	Eval(loc []float64) (obj float64, err error)
}

type MisoGradBasedProblem interface {
	Eval(loc []float64) (obj float64, grad []float64, err error)
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
