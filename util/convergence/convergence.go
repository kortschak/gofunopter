package convergence

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
var GradAbsTol Convergence = GradConvergence{"Gradient absolute tolerance reached"}
var GradRelTol Convergence = GradConvergence{"Gradient relative tolerance reached"}

type Obj struct{ Str string }

func (o Obj) Convergence() string {
	return o.Str
}

var ObjAbsTol Convergence = FunConvergence{"Function absolute tolerance reached"}
var ObjRelTol Convergence = FunConvergence{"Function relative tolerance reached"}

type StepConvergence struct{ Str string }

func (s StepConvergence) Convergence() string {
	return s.Str
}

var StepAbsTol Convergence = StepConvergence{"Step absolute tolerance reached"}
var StepRelTol Convergence = StepConvergence{"Step relative tolerance reached"}
