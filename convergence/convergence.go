package convergence

const DefaultGradAbsTol = 1E-6

// A converger is a type that can test if the optimization has converged
type Converger interface {
	Converged() C
}

// C is a basic interface for expressing methods of optimizer convergence
type C interface {
	Convergence() string
}

// Grad is a type marking the convergence of the optimizer because of the gradient
type Grad struct{ Str string }

// Converged returns the specific string for the convergence
func (g Grad) Convergence() string {
	return g.Str
}

// GradAbsTol is a convergence because of meaning the absolute tolerance of the gradient
var GradAbsTol Grad = Grad{"Gradient absolute tolerance reached"}
var GradRelTol Grad = Grad{"Gradient relative tolerance reached"}

type Obj struct{ Str string }

func (o Obj) Convergence() string {
	return o.Str
}

var ObjAbsTol Obj = Obj{"Function absolute tolerance reached"}
var ObjRelTol Obj = Obj{"Function relative tolerance reached"}

type Step struct{ Str string }

func (s Step) Convergence() string {
	return s.Str
}

var StepAbsTol Step = Step{"Step absolute tolerance reached"}
var StepRelTol Step = Step{"Step relative tolerance reached"}
