package gofunopter

type Convergence interface {
	ConvergenceType() string
}

type Optimizer interface {
	Converged() Convergence
	Initialize() error
	Result()
	Iterate() error
}

type SISOGradBasedOptimizer interface {
	Optimizer
	Loc() LocationFloat
	Obj() ObjectiveFloat
	Grad() GradientFloat
	Fun() MISOGradBasedProblem
}

/*
type MISOGradBasedOptimizer interface {
	Optimizer
	Loc() OptFloatSlice
	Obj() OptFloat
	Grad() OptFloatSlice
	Fun() MISOGradBasedProblem
}
*/
/*
type OptFloat interface {
	HistorySaverFloat
	CurrGetSetterFloat
	InitGetSetterFloat
	Displayer
}

type OptTolFloat interface {
	OptFloat
	AbsToler
	RelToler
}
*/
/*
type BoundedOptTolFloat interface {
	OptTolFloat
}
*/
type LocationFloat interface {
	HistorySaverFloat
	CurrGetSetterFloat
	InitGetSetterFloat
	Converger
	Displayable
	Initializer
	Resulter
}

type ObjectiveFloat interface {
	HistorySaverFloat
	CurrGetSetterFloat
	InitGetSetterFloat
	Converger
	Displayable
	AbsToler
	RelToler
	Initializer
	Resulter
}

type GradientFloat interface {
	HistorySaverFloat
	CurrGetSetterFloat
	InitGetSetterFloat
	Converger
	Displayer
	AbsToler
	Displayable
	RelToler
	Initializer
	Resulter
}

type BoundedFloat interface {
	Lb() float64
	SetLb(float64)
	Ub() float64
	SetUb(float64)
	AbsToler
	RelToler
	Midpoint() float64
	WithinBounds(float64) bool
}

type BoundedStepFloat interface {
	CurrGetSetterFloat
	InitGetSetterFloat
	Displayable
	BoundedFloat
	Initializer
	Converger
	Resulter
}

// The error is for error checking
type Initializer interface {
	Initialize() error
}

func Initialize(initializers ...Initializer) (err error) {
	for _, initializer := range initializers {
		err := initializer.Initialize()
		if err != nil {
			return err
		}
	}
	return nil
}

type HistoryFloat interface {
	Add(float64)
	Save() bool
	SetSave(bool)
}

// TODO: Replace this with an interface
type HistorySaverFloat interface {
	Hist() HistoryFloat
}

type Iterator interface {
	Iterate() error
}

func Iterate(iterators ...Iterator) (err error) {
	for _, iterator := range iterators {
		err := iterator.Iterate()
		if err != nil {
			return err
		}
	}
	return nil
}

type Converger interface {
	Converged() Convergence
}

func Converged(convergers ...Converger) (convergence Convergence) {
	for _, converger := range convergers {
		convergence = converger.Converged()
		if convergence != nil {
			return convergence
		}
	}
	return nil
}

type Resulter interface {
	SetResult() // Set the result of the optimization (usually goes through the opter interface)
}

type OpterFloat interface {
	Opt() float64
	SetOpt(float64)
}

func SetResults(resulters ...Resulter) {
	for _, resulter := range resulters {
		resulter.Result()
	}
}

type CurrGetSetterFloat interface {
	SetCurr(float64)
	Curr() float64
}

type InitGetSetterFloat interface {
	SetInit(float64)
	Init() float64
}

type CurrInitGetSetterFloat interface {
	CurrGetSetterFloat
	InitGetSetterFloat
}

type TolGetSetterFloat interface {
	SetTol(float64)
	Tol() float64
}

type AbsTol interface {
	CurrGetSetterFloat
	TolGetSetterFloat
	Converger
}

type RelTol interface {
	CurrInitGetSetterFloat
	TolGetSetterFloat
	Converger
}

type RelToler interface {
	RelTol() RelTol
}

type AbsToler interface {
	AbsTol() AbsTol
}

// A basic float type that can be used in an optimizer
// All the set methods are so users can tune the default behavior
// easily
// TODO: Decide if the setter methods should return an error

/*
type OptFloat interface {
	Optimizer
	Displayable
	Init() float64
	Curr() float64
	Opt() float64
	Hist() HistoryFloat
	RelTol() float64
	AbsTol() float64
	SetCurr() float64
	SetInit(float64)
	SetAbsTol(float64)
	SetRelTol(float64)
	SetName(string)
	SetDisp(bool) // Display this output
}

type OptFloatSlice interface {
	Optimizer
	Init() []float64
	Curr() []float64
	Opt() []float64
	Hist() HistoryFloatSlice
	SetInit([]float64)
	SetAbsTol(float64)
	SetRelTol(float64)
	SetName(string)
	SetDisp(bool) // Display this output
}

type BoundedOptFloat interface {
	OptFloat
	Lb() float64
	Ub() float64
	WithinBounds() bool
	Midpoint() float64
	SetLb(float64)
	SetUb(float64)
}
*/
