package multivariate

import (
	"github.com/btracey/gofunopter/common"
	"github.com/btracey/gofunopter/common/multi"
	"github.com/btracey/gofunopter/common/uni"
	"github.com/btracey/gofunopter/convergence"
	"github.com/btracey/gofunopter/display"
	"github.com/btracey/gofunopter/optimize"
)

func OptimizeGrad(function multiGradFun, initialLocation []float64, settings MultiGradSettings, optimizer MultiGradOptimizer) (optValue float64, optLocation []float64, result MultiGradResult, err error) {

	m := NewMisoGradStruct()
	m.fun = function
	m.settings = settings
	m.optimizer = optimizer

	m.loc.SetInit(initialLocation)
	convergence, err = optimize.OptimizeOpter(m, optimizer, fun)

}

type MultiGradResult struct {
	Convergence convergence.C
}

type MultiGradSettings struct {
	InitialObjective          float64
	InitialGradient           []float64
	GradientAbsoluteTolerance float64
	Display                   bool
}

type MisoGradStruct struct {
	*common.OptCommon
	*display.Display
	disp bool

	loc  *multi.Location
	obj  *uni.Objective
	grad *multi.Gradient

	// User defined function
	fun optimize.MisoGrad

	// Optimization model
	optimizer MultiGradOptimizer

	// Settings
	settings MultiGradSettings

	// result
	result MultiGradResult
}

func NewMisoGradStruct() *MisoGradStruct {
	m := &MisoGradStruct{
		OptCommon: common.NewOptCommon(),
		Display:   display.NewDisplay(),
		disp:      true,
		loc:       multi.NewLocation(),
		obj:       uni.NewObjective(),
		grad:      multi.NewGradient(),
	}
	// TODO: Something about settings
}

func (m *MisoGradStruct) SetSettings() {
	m.grad.SetInit(m.settings.InitialGradient)
	m.obj.SetInit(m.settings.InitialObjective)
	m.disp = m.settings.Display
	m.grad.SetAbsTol(m.settings.GradientAbsoluteTolerance)
}

func (m *MisoGradStruct) Disp() bool {
	return m.disp
}

func (m *MisoGradStruct) SetDisp(b bool) {
	m.disp = b
}

func (m *MisoGradStruct) Converged() convergence.C {
	return convergence.CheckConvergence(m.obj, m.grad)
}

func (m *MisoGradStruct) AddToDisplay(d []*display.Struct) []*display.Struct {
	if m.disp {
		d = display.AddToDisplay(d, lbfgs.loc, lbfgs.obj, lbfgs.grad, lbfgs.step)
	}
}

func (m *Miso) SetResult(c convergence.C) {
	m.result.Convergence = c
}
