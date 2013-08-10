package multivariate

import (
	"github.com/btracey/gofunopter/common"
	//"github.com/btracey/gofunopter/common/multi"
	"github.com/btracey/gofunopter/common/display"
	"github.com/btracey/gofunopter/common/multi"
	"github.com/btracey/gofunopter/common/optimize"
	"github.com/btracey/gofunopter/common/status"
	"github.com/btracey/gofunopter/common/uni"

	"errors"
	"math"

	//"fmt"
)

type moddedFun struct {
	fun      optimize.MultiObjGrad
	loc      *multi.Location
	obj      *uni.Objective
	grad     *multi.Gradient
	funEvals *common.FunctionEvaluations
}

func newModdedFun(fun optimize.MultiObjGrad, loc *multi.Location, obj *uni.Objective, grad *multi.Gradient, funEvals *common.FunctionEvaluations) *moddedFun {
	return &moddedFun{
		fun:      fun,
		loc:      loc,
		obj:      obj,
		grad:     grad,
		funEvals: funEvals,
	}
}

func (m *moddedFun) ObjGrad(x []float64) (obj float64, grad []float64, err error) {
	obj, grad, err = m.fun.ObjGrad(x)
	m.loc.AddToHist(x)
	m.obj.AddToHist(obj)
	m.grad.AddToHist(grad)
	m.funEvals.Add(1)
	return
}

type MultiGradOptimizer interface {
	Initialize(loc *multi.Location, obj *uni.Objective, grad *multi.Gradient) error
	Iterate(loc *multi.Location, obj *uni.Objective, grad *multi.Gradient, fun optimize.MultiObjGrad) (status.Status, error)
}

func OptimizeGrad(function optimize.MultiObjGrad, initialLocation []float64, settings *MultiGradSettings, optimizer MultiGradOptimizer) (optValue float64, optLocation []float64, result *MultiGradResult, err error) {

	if settings == nil {
		settings = NewMultiGradSettings()
	}

	if optimizer == nil {
		optimizer = NewLbfgs()
	}

	m := newMultiGradStruct()
	//m.fun = function
	m.fun = newModdedFun(function, m.loc, m.obj, m.grad, m.FunEvals)
	m.settings = settings
	m.optimizer = optimizer

	m.loc.SetInit(initialLocation)
	c, err := optimize.OptimizeOpter(m, function)
	m.result.Status = c
	return m.obj.Opt(), m.loc.Opt(), m.result, err
}

type MultiGradResult struct {
	Status status.Status
	*common.CommonResult
}

type MultiGradSettings struct {
	*common.CommonSettings
	InitialObjective          float64
	InitialGradient           []float64
	GradientAbsoluteTolerance float64
}

func NewMultiGradSettings() *MultiGradSettings {
	return &MultiGradSettings{
		CommonSettings:            common.NewCommonSettings(),
		InitialObjective:          math.NaN(),
		InitialGradient:           nil,
		GradientAbsoluteTolerance: 1e-6,
	}
}

type multiGradStruct struct {
	*common.OptCommon

	loc  *multi.Location
	obj  *uni.Objective
	grad *multi.Gradient

	// User defined function
	fun optimize.MultiObjGrad

	// Optimization model
	optimizer MultiGradOptimizer

	// Settings
	settings *MultiGradSettings

	// result
	result *MultiGradResult
}

func newMultiGradStruct() *multiGradStruct {
	m := &multiGradStruct{
		OptCommon: common.NewOptCommon(),
		loc:       multi.NewLocation(),
		obj:       uni.NewObjective(),
		grad:      multi.NewGradient(),
		result:    &MultiGradResult{},
	}
	// TODO: Something about settings
	return m
}

func (m *multiGradStruct) CommonSettings() *common.CommonSettings {
	return m.settings.CommonSettings
}

func (m *multiGradStruct) SetSettings() error {
	m.grad.SetInit(m.settings.InitialGradient)
	m.obj.SetInit(m.settings.InitialObjective)
	m.grad.SetAbsTol(m.settings.GradientAbsoluteTolerance)
	return nil
}

func (m *multiGradStruct) Status() status.Status {
	return status.CheckStatus(m.obj, m.grad)
}

func (m *multiGradStruct) AddToDisplay(d []*display.Struct) []*display.Struct {
	//fmt.Println("In multi add to display")
	return display.AddToDisplay(d, m.loc, m.obj, m.grad)
}

func (m *multiGradStruct) SetResult(c *common.CommonResult) {
	optimize.SetResult(m.loc, m.grad, m.obj)
	m.result.CommonResult = c

	setResulter, ok := m.optimizer.(optimize.SetResulter)
	if ok {
		setResulter.SetResult()
	}
}

func (m *multiGradStruct) Initialize() error {
	initLoc := m.loc.Init()
	initObj := m.obj.Init()
	initGrad := m.grad.Init()

	// The initial values need to both be NaN or both not nan
	if math.IsNaN(initObj) {
		if len(initGrad) != 0 {
			return errors.New("initial function value and gradient must either both be set or neither set")
		}
		// Both nan, so compute the initial fuction value and gradient
		initObj, initGrad, err := m.fun.ObjGrad(initLoc)
		if err != nil {
			return errors.New("error calling function during optimization: \n" + err.Error())
		}
		m.obj.SetInit(initObj)
		m.grad.SetInit(initGrad)
	} else {
		if len(initGrad) == 0 {
			return errors.New("initial function value and gradient must either both be set or neither set")
		}
	}

	err := optimize.Initialize(m.loc, m.obj, m.grad)
	if err != nil {
		return err
	}
	err = m.optimizer.Initialize(m.loc, m.obj, m.grad)
	if err != nil {
		return err
	}
	return nil
}

func (m *multiGradStruct) Iterate() (status.Status, error) {
	return m.optimizer.Iterate(m.loc, m.obj, m.grad, m.fun)
}
