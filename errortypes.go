package gofunopter

type FunctionError struct {
	Err error
	Loc interface{}
}

func (f *FunctionError) Error() string {
	return "Error evaluating function: " + f.Err.Error()
}

type InitializationError struct {
	Err error
}

func (i *InitializationError) Error() string {
	return "Error initializing: " + i.Err.Error()
}

type LinesearchError struct {
	Err error
}

func (l *LinesearchError) Error() string {
	return "Error in linesearch optimization: " + l.Err.Error()
}

type LinesearchConvergenceFalure struct {
	Conv     Convergence
	Loc      []float64 // Starting location for the line search
	InitStep []float64
}

func (l *LinesearchConvergenceFalure) Error() string {
	return "Linesearch did not meet Wolfe Conditions. Instead converged with " + l.Conv.ConvergenceType()
}
