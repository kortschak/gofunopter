package convergence

// CheckConvergence checks the convergence of a variadic
// number of converges and returns the first non-nil result
func CheckConvergence(cs ...Converger) C {
	for _, val := range cs {
		c := val.Converged()
		if c != nil {
			return c
		}
	}
	return nil
}

// Use type casting for varieties of convergence (grad, etc.)
// use call to convergence for specific convergence test

// A converger is a type that can test if the optimization has converged
type Converger interface {
	Converged() C
}

// C is a basic interface for expressing methods of optimizer convergence
type C interface {
	Convergence() string
}

type Basic struct{ Str string }

func (b Basic) Convergence() string {
	return b.Str
}

func (b Basic) String() string {
	return b.Str
}

// Grad is a type marking the convergence of the optimizer because of the gradient
type Grad struct{ Basic }

// GradAbsTol is a convergence because of meaning the absolute tolerance of the gradient
var GradAbsTol Grad = Grad{Basic{"convergence: gradient absolute tolerance reached"}}
var GradRelTol Grad = Grad{Basic{"convergence: gradient relative tolerance reached"}}

type Obj struct{ Basic }

var ObjAbsTol Obj = Obj{Basic{"convergence: function absolute tolerance reached"}}
var ObjRelTol Obj = Obj{Basic{"convergence: function relative tolerance reached"}}

type Step struct{ Basic }

var StepAbsTol Step = Step{Basic{"convergence: step absolute tolerance reached"}}
var StepRelTol Step = Step{Basic{"convergence: step relative tolerance reached"}}

var Iterations Basic = Basic{"convergence: maximum iterations reached"}
var FunEvals Basic = Basic{"convergence: maximum function evaluations reached"}
var Time Basic = Basic{"convergence: maximum time elapsed"}
