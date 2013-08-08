package univariate

import (
	"github.com/btracey/gofunopter/common"
	//"github.com/btracey/gofunopter/common/multi"
	"github.com/btracey/gofunopter/common/uni"
	"github.com/btracey/gofunopter/convergence"
	"github.com/btracey/gofunopter/display"
	"github.com/btracey/gofunopter/optimize"

	"errors"
	"math"
)

type UniGradFun interface {
	Eval(x float64) (f float64, obj float64, err error)
}

type UniGradOptimizer interface {
	Initialize(loc *uni.Location, obj *uni.Objective, grad *uni.Gradient) error
	Iterate(loc *uni.Location, obj *uni.Objective, grad *uni.Gradient, fun UniGradFun) (nFunEVals int, err error)
}

func OptimizeGrad(function UniGradFun, initialLocation float64, settings *UniGradSettings, optimizer UniGradOptimizer) (optValue float64, optLocation float64, result *UniGradResult, err error) {

	m := newUniGradStruct()
	m.fun = function
	m.settings = settings
	m.optimizer = optimizer

	m.loc.SetInit(initialLocation)
	c, err := optimize.OptimizeOpter(m, function)
	m.result.Convergence = c
	return m.obj.Opt(), m.loc.Opt(), m.result, err
}

type UniGradResult struct {
	Convergence convergence.Type
	*common.CommonResult
}

type UniGradSettings struct {
	*common.CommonSettings
	InitialObjective          float64
	InitialGradient           float64
	GradientAbsoluteTolerance float64
	Display                   bool
}

func NewUniGradSettings() *UniGradSettings {
	return &UniGradSettings{
		CommonSettings:            common.NewCommonSettings(),
		InitialObjective:          math.NaN(),
		InitialGradient:           math.NaN(),
		GradientAbsoluteTolerance: 1e-6,
		Display:                   true,
	}
}

type uniGradStruct struct {
	*common.OptCommon
	*display.Display
	disp bool

	loc  *uni.Location
	obj  *uni.Objective
	grad *uni.Gradient

	// User defined function
	fun UniGradFun

	// Optimization model
	optimizer UniGradOptimizer

	// Settings
	settings *UniGradSettings

	// result
	result *UniGradResult
}

func newUniGradStruct() *uniGradStruct {
	m := &uniGradStruct{
		OptCommon: common.NewOptCommon(),
		Display:   display.NewDisplay(),
		disp:      true,
		loc:       uni.NewLocation(),
		obj:       uni.NewObjective(),
		grad:      uni.NewGradient(),
		result:    &UniGradResult{},
	}
	// TODO: Something about settings
	return m
}

func (m *uniGradStruct) CommonSettings() *common.CommonSettings {
	return m.settings.CommonSettings
}

func (m *uniGradStruct) SetSettings() error {
	m.grad.SetInit(m.settings.InitialGradient)
	m.obj.SetInit(m.settings.InitialObjective)
	m.disp = m.settings.Display
	m.grad.SetAbsTol(m.settings.GradientAbsoluteTolerance)
	return nil
}

func (m *uniGradStruct) Disp() bool {
	return m.disp
}

func (m *uniGradStruct) SetDisp(b bool) {
	m.disp = b
}

func (m *uniGradStruct) Converged() convergence.Type {
	return convergence.CheckConvergence(m.obj, m.grad)
}

func (m *uniGradStruct) AddToDisplay(d []*display.Struct) []*display.Struct {
	if m.disp {
		d = display.AddToDisplay(d, m.loc, m.obj, m.grad)
	}
	return d
}

func (u *uniGradStruct) SetResult(c *common.CommonResult) {
	optimize.SetResult(u.loc, u.grad, u.obj)
	u.result.CommonResult = c

	setResulter, ok := u.optimizer.(optimize.SetResulter)
	if ok {
		setResulter.SetResult()
	}
}

func (u *uniGradStruct) Initialize() error {
	initLoc := u.loc.Init()
	initObj := u.obj.Init()
	initGrad := u.grad.Init()

	// The initial values need to both be NaN or both not nan
	if math.IsNaN(initObj) {
		if !math.IsNaN(initGrad) {
			return errors.New("gofunopter: cubic: initial function value and gradient must either both be set or neither set")
		}
		// Both nan, so compute the initial fuction value and gradient
		initObj, initGrad, err := u.fun.Eval(initLoc)
		if err != nil {
			return errors.New("gofunopter: cubic: error calling function during optimization")
		}
		u.obj.SetInit(initObj)
		u.grad.SetInit(initGrad)
	} else {
		if math.IsNaN(initGrad) {
			return errors.New("gofunopter: cubic: initial function value and gradient must either both be set or neither set")
		}
	}

	err := optimize.Initialize(u.loc, u.obj, u.grad)
	if err != nil {
		return err
	}
	err = u.optimizer.Initialize(u.loc, u.obj, u.grad)
	if err != nil {
		return err
	}
	return nil
}

func (u *uniGradStruct) Iterate() (nFunEvals int, err error) {
	return u.optimizer.Iterate(u.loc, u.obj, u.grad, u.fun)
}
