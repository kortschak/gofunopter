package gofunopter

func Optimize(o Optimizer) (c Convergence, err error) {
	// Add in some check about nil pointers and such
	err = o.Initialize()
	if err != nil {
		return nil, &InitializationError{Err: err}
	}
	// Want to return the result even if there is an error in case anything
	// gets lost (maybe a defer would be even better?)
	defer o.SetResult()
	// Iterate until convergence
	for {
		c := o.Converged()
		if c != nil {
			return c, nil
		}
		err = o.Iterate()
		if err != nil {
			break
		}
	}
	return nil, err
}

type Optimizer interface {
	Converger
	Initializer
	SetResulter
	Iterator
}

type CommonOpter interface {
	FunEvals() *FunctionEvaluations
	Iter() *Iterations
	SetDisp(bool)
}

type SisoGradBasedOptimizer interface {
	Optimizer
	CommonOpter
	Loc() OptFloat
	Obj() OptTolFloat
	Grad() OptTolFloat
	Fun() SisoGradBasedProblem
	SetFun(SisoGradBasedProblem)
}

type SisoStepOptimizer interface {
	SisoGradBasedOptimizer
	Step() BoundedOptFloat
}

type MisoGradBasedOptimizer interface {
	Optimizer
	CommonOpter
	Loc() OptFloatSlice
	Obj() OptTolFloat
	Grad() OptTolFloatSlice
	Fun() MisoGradBasedProblem
}

// TODO: Add in some mechanism for "Default" optimization selection
