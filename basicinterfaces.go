package gofunopter

type Convergence interface {
	ConvergenceType() string
}

type OptFloat interface {
	HistoryFloat
	CurrerFloat
	IniterFloat
	Displayer
	Initializer
	SetResulter
	OpterFloat
}

type OptTolFloat interface {
	Converger
	OptFloat
	AbsToler
	RelToler
}

type BoundedOptFloat interface {
	OptTolFloat
	Lb() float64
	SetLb(float64)
	Ub() float64
	SetUb(float64)
	Midpoint() float64
	WithinBounds(float64) bool
}

type OptFloatSlice interface {
	HistoryFloatSlice
	CurrerFloatSlice
	IniterFloatSlice
	Displayer
	Initializer
	SetResulter
	OpterFloatSlice
}

type OptTolFloatSlice interface {
	Converger
	OptFloatSlice
	AbsToler
	RelToler
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
	AddToHist(float64)
	Save() bool
	SetSave(bool)
}

type HistoryFloatSlice interface {
	AddToHist([]float64)
	Save() bool
	SetSave(bool)
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

// Sets the optimum from the end of the run
type SetResulter interface {
	SetResult()
}

type OpterFloat interface {
	Opt() float64
}
type OpterFloatSlice interface {
	Opt() []float64
}

func SetResults(resulters ...SetResulter) {
	for _, resulter := range resulters {
		resulter.SetResult()
	}
}

type CurrerFloat interface {
	SetCurr(float64)
	Curr() float64
}

type IniterFloat interface {
	SetInit(float64)
	Init() float64
}

type CurrerFloatSlice interface {
	SetCurr([]float64)
	Curr() []float64
}

type IniterFloatSlice interface {
	SetInit([]float64)
	Init() []float64
}

type RelToler interface {
	RelTol() float64
	SetRelTol(float64)
}

type AbsToler interface {
	AbsTol() float64
	SetAbsTol(float64)
}
