package univariate

import (
	"github.com/btracey/gofunopter/common"
	//"github.com/btracey/gofunopter/common/multi"
	"github.com/btracey/gofunopter/common/display"
	"github.com/btracey/gofunopter/common/optimize"
	"github.com/btracey/gofunopter/common/status"
	"github.com/btracey/gofunopter/common/uni"

	"errors"
	"math"
)

type moddedFun struct {
	uni      optimize.UniObjGrad
	loc      *uni.Location
	obj      *uni.Objective
	grad     *uni.Gradient
	funEvals *common.FunctionEvaluations
}

func newModdedFun(fun optimize.UniObjGrad, loc *uni.Location, obj *uni.Objective, grad *uni.Gradient, funEvals *common.FunctionEvaluations) *moddedFun {
	return &moddedFun{
		uni:      fun,
		loc:      loc,
		obj:      obj,
		grad:     grad,
		funEvals: funEvals,
	}
}

func (m *moddedFun) ObjGrad(x float64) (obj float64, grad float64, err error) {
	obj, grad, err = m.uni.ObjGrad(x)
	m.loc.AddToHist(x)
	m.obj.AddToHist(obj)
	m.grad.AddToHist(grad)
	m.funEvals.Add(1)
	return
}

type UniGradOptimizer interface {
	Initialize(loc *uni.Location, obj *uni.Objective, grad *uni.Gradient) error
	Iterate(loc *uni.Location, obj *uni.Objective, grad *uni.Gradient, fun optimize.UniObjGrad) (status.Status, error)
}

func OptimizeGrad(function optimize.UniObjGrad, initialLocation float64, settings *UniGradSettings, optimizer UniGradOptimizer) (optValue float64, optLocation float64, result *UniGradResult, err error) {

	if settings == nil {
		settings = NewUniGradSettings()
	}

	if optimizer == nil {
		optimizer = NewCubic()
	}

	m := newUniGradStruct()

	m.fun = newModdedFun(function, m.loc, m.obj, m.grad, m.FunEvals)
	m.settings = settings
	m.optimizer = optimizer

	m.loc.SetInit(initialLocation)
	err = optimize.OptimizeOpter(m, function)
	//m.result.Status = c
	return m.obj.Opt(), m.loc.Opt(), m.Result(), err
}

type UniGradResult struct {
	Status status.Status
	*common.CommonResult
	*uni.ObjectiveResult
	*uni.GradientResult
}

type UniGradSettings struct {
	*common.CommonSettings
	*uni.GradientSettings
	*uni.ObjectiveSettings
}

func NewUniGradSettings() *UniGradSettings {
	return &UniGradSettings{
		CommonSettings:    common.NewCommonSettings(),
		GradientSettings:  uni.NewGradientSettings(),
		ObjectiveSettings: uni.NewObjectiveSettings(),
	}
}

type uniGradStruct struct {
	*common.OptCommon

	loc  *uni.Location
	obj  *uni.Objective
	grad *uni.Gradient

	// User defined function
	fun optimize.UniObjGrad

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
	m.obj.SetSettings(m.settings.ObjectiveSettings)
	m.grad.SetSettings(m.settings.GradientSettings)
	//m.grad.SetInit(m.settings.InitialGradient)
	//m.obj.SetInit(m.settings.InitialObjective)
	//m.grad.SetAbsTol(m.settings.GradientAbsoluteTolerance)
	return nil
}

func (m *uniGradStruct) Status() status.Status {
	return status.CheckStatus(m.obj, m.grad)
}

func (m *uniGradStruct) AddToDisplay(d []*display.Struct) []*display.Struct {
	d = display.AddToDisplay(d, m.loc, m.obj, m.grad)
	return d
}

func (u *uniGradStruct) Result() *UniGradResult {
	return &UniGradResult{
		CommonResult:    u.OptCommon.CommonResult(),
		ObjectiveResult: u.obj.Result(),
		GradientResult:  u.grad.Result(),
	}
}

func (u *uniGradStruct) SetResult() {
	optimize.SetResult(u.loc, u.grad, u.obj)
	//u.result.CommonResult = c

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
		initObj, initGrad, err := u.fun.ObjGrad(initLoc)
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

func (u *uniGradStruct) Iterate() (stat status.Status, err error) {
	return u.optimizer.Iterate(u.loc, u.obj, u.grad, u.fun)
}
