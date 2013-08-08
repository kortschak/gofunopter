package optimize

import (
	"github.com/btracey/gofunopter/common"
	"github.com/btracey/gofunopter/convergence"
	"github.com/btracey/gofunopter/display"
)

type Optimizer interface {
	Initializer
	//SetResulter
	convergence.Converger
	display.Displayer
	Disp() bool
	SetDisp(bool)
	Iterate() (int, error)
	GetOptCommon() *common.OptCommon
	GetDisplay() *display.Display
	FunEvals() *common.FunctionEvaluations
	Iter() *common.Iterations
	SetSettings() error
	CommonSettings() *common.CommonSettings
	SetResult(*common.CommonResult)
}

/*
type SisoGradOptimizer interface {
	Optimizer
	Loc() *uni.Location
	Obj() *uni.Objective
	Grad() *uni.Gradient
	Fun() SisoGrad
	Optimize(SisoGrad, float64) (float64, float64, convergence.C, error)
}

type MisoGradOptimizer interface {
	Optimizer
	Loc() *multi.Location
	Obj() *uni.Objective
	Grad() *multi.Gradient
	Fun() MisoGrad
	Optimize(MisoGrad, []float64) (float64, []float64, convergence.C, error)
}
*/
