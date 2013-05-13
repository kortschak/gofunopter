package gofunopter

func Optimize(o Optimizer) (c Convergence, err error) {
	// Add in some check about nil pointers and such
	err = o.Initialize()
	if err != nil {
		return InitializationError(err), InitializationError(err)
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
	return InitializationError(err), InitializationError(err)
}

type Optimizer interface {
	Converger
	Initializer
	SetResulter
	Iterator
}

type SisoGradBasedOptimizer interface {
	Optimizer
	Location() OptFloat
	Objective() OptTolFloat
	Gradient() OptTolFloat
	Function() MisoGradBasedProblem
}

// TODO: Add in some mechanism for "Default" optimization selection
