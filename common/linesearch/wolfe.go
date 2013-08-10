package linesearch

import (
	"github.com/btracey/gofunopter/common/status"
	"math"
)

// WolfeConditioner is an iterface for wolfe conditions (strong or weak)
type WolfeConditioner interface {
	Status() status.Status
	FunConst() float64
	GradConst() float64
	SetFunConst(funConst float64)
	SetGradConst(gradConst float64)
	SetInitState(initObj, initGrad float64)
	SetCurrState(currObj, currGrad, step float64)
}

type WeakWolfeConditions struct {
	funConst  float64
	gradConst float64
	currObj   float64
	currGrad  float64
	initObj   float64
	initGrad  float64
	step      float64
}

func (w *WeakWolfeConditions) WolfeConditions() WolfeConditioner {
	return w
}

func (w *WeakWolfeConditions) SetInitState(initObj, initGrad float64) {
	w.initObj = initObj
	w.initGrad = initGrad
	w.step = math.Inf(1)
}

func (w *WeakWolfeConditions) SetCurrState(currObj, currGrad, currStep float64) {
	w.currObj = currObj
	w.currGrad = currGrad
	w.step = currStep
}

//func (s *WeakWolfeConditions) WolfeConditionsMet(obj, directionalderivative, step float64) bool {
func (w *WeakWolfeConditions) Status() status.Status {
	if w.currObj >= w.initObj+w.funConst*w.step*w.currGrad {
		return status.Continue
	}
	if w.currGrad <= w.gradConst*w.initGrad {
		return status.Continue
	}
	return status.WolfeConditionsMet
}

func (w *WeakWolfeConditions) SetFunConst(val float64) {
	w.funConst = val
}

func (w *WeakWolfeConditions) SetGradConst(val float64) {
	w.gradConst = val
}

func (w *WeakWolfeConditions) FunConst() float64 {
	return w.funConst
}

func (w *WeakWolfeConditions) GradConst() float64 {
	return w.gradConst
}

type StrongWolfeConditions struct {
	funConst  float64
	gradConst float64
	currObj   float64
	currGrad  float64
	initObj   float64
	initGrad  float64
	step      float64
}

func (s *StrongWolfeConditions) WolfeConditions() WolfeConditioner {
	return s
}

func (s *StrongWolfeConditions) SetInitState(initObj, initGrad float64) {
	s.initObj = initObj
	s.initGrad = initGrad
	s.step = math.Inf(1)
}

func (s *StrongWolfeConditions) SetCurrState(currObj, currGrad, currStep float64) {
	s.currObj = currObj
	s.currGrad = currGrad
	s.step = currStep
}

func (s *StrongWolfeConditions) SetFunConst(val float64) {
	s.funConst = val
}

func (s *StrongWolfeConditions) SetGradConst(val float64) {
	s.gradConst = val
}

func (s *StrongWolfeConditions) FunConst() float64 {
	return s.funConst
}

func (s *StrongWolfeConditions) GradConst() float64 {
	return s.gradConst
}

func (s *StrongWolfeConditions) Status() status.Status {
	if s.currObj >= s.initObj+s.funConst*s.step*s.currGrad {
		return status.Continue
	}
	if math.Abs(s.currGrad) >= s.gradConst*math.Abs(s.initGrad) {
		return status.Continue
	}
	return status.WolfeConditionsMet
}
