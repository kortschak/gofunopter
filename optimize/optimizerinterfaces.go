package optimize

import (
	"github.com/btracey/gofunopter/common"
	"github.com/btracey/gofunopter/convergence"
	"github.com/btracey/gofunopter/display"
	"github.com/btracey/gofunopter/multivariate"
	"github.com/btracey/gofunopter/univariate"
)

type Optimizer interface {
	Initializer
	SetResulter
	convergence.Converger
	display.Displayer
	Disp() bool
	SetDisp(bool)
	Iterate() (int, error)
	GetOptCommon() *common.OptCommon
	GetDisplay() *display.Display
	FunEvals() *common.FunctionEvaluations
	Iter() *common.Iterations
}

type SisoGradOptimizer interface {
	Optimizer
	Loc() *univariate.Location
	Obj() *univariate.Objective
	Grad() *univariate.Gradient
	Fun() SisoGrad
	Optimize(SisoGrad) (convergence.C, error)
}

type MisoGradOptimizer interface {
	Optimizer
	Loc() *multivariate.Location
	Obj() *univariate.Objective
	Grad() *multivariate.Gradient
	Fun() MisoGrad
	Optimize(MisoGrad) (convergence.C, error)
}
