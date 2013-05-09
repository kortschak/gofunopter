package gofunopter

type Optimizer interface {
	Converged() string
	Initialize() error
	Result()
	Iterate() error
}

func Optimize(o Optimizer) (str string, err error) {
	// Add in some check about nil pointers and such
	err = o.Initialize()
	if err != nil {
		return "", err
	}
	// Want to return the result even if there is an error in case anything
	// gets lost (maybe a defer would be even better?)
	defer o.Result()
	// Iterate until convergence
	for {
		str := o.Converged()
		if str != "" {
			return str, nil
		}
		err = o.Iterate()
		if err != nil {
			break
		}
	}
	return "", err
}

// TODO: Add in some mechanism for "Default" optimization selection
