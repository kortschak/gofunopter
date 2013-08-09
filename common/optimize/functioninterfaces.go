package optimize

/*

// This is in its own package because otherwise it's hard to avoid circular imports

// Siso is an interface for a function input of a
//gradient-free unconstrained optimizer
type Siso interface {
	Eval(x float64) (f float64)
}

// SisoGrad is an interface for a function which is single
// input single output and has gradient information
type SisoGrad interface {
	Eval(x float64) (f, g float64, err error)
}

/*
// SisoGradIndividual is an interface for a function which
// is a SisoGrad, but the function and the gradient can be
// called separately
type SisoGradIndividual interface {
	SisoGrad
	Function(x float64) (f float64)
	Gradient(x float64) (g float64)
}

type Miso interface {
	Eval(x []float64) (f float64)
}

type MisoGrad interface {
	Eval(x []float64) (f float64, g []float64, err error)
}
*/
