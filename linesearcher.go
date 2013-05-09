package gofunopter

import "math"

// TODO: Add in error checking for positive initial gradient? Maybe should be a panic
// because it shouldn't ever occur

type WolfeConditioner interface {
	WolfeConditionsMet(obj, directionalderivative, step float64) bool
	SetInit(initObj, initGrad float64)
}

type WeakWolfeConditions struct {
	FunConst  float64
	GradConst float64
	initObj   float64
	initGrad  float64
}

func (s *WeakWolfeConditions) SetInit(obj, grad float64) {
	s.initObj = obj
	s.initGrad = grad
}

func (s *WeakWolfeConditions) WolfeConditionsMet(obj, directionalderivative, step float64) bool {
	if obj >= s.initObj+s.FunConst*step*directionalderivative {
		return false
	}
	if directionalderivative <= s.GradConst*s.initGrad {
		return false
	}
	return true
}

type StrongWolfeConditions struct {
	FunConst    float64
	GradConst   float64
	initObj     float64
	initGrad    float64
	absInitGrad float64
}

func (s *StrongWolfeConditions) SetInit(obj, grad float64) {
	s.initObj = obj
	s.initGrad = grad
	s.absInitGrad = math.Abs(initGrad)
}

func (s *StrongWolfeConditions) WolfeConditionsMet(obj, directionalderivative, step float64) bool {
	if obj >= s.initObj+s.FunConst*step*directionalderivative {
		return false
	}
	if math.Abs(directionalderivative) >= s.GradConst*s.absInitGrad {
		return false
	}
	return true
}

// Maybe everything should be through interfaces to make everything easier
// to set. Make OptFloat an interface. Also probably makes it easier to customize.
// Harder to save possibly, but not hard to just save the float values

type Linesearcher interface {
	Optimizer
	Loc() *LocationFloat
	Obj() *ObjectiveFloat
	Grad() *GradientFloat
}

type Linesearchable interface {
	Linesearcher() *Linesearcher
	//Loc() *OptFloatSlice
	//Opt() *OptFloat
	//Grad() *OptFloatSlice
	WolfeFunConst() float64
	WolfeGradConst() float64
}

type CubicLinesearch struct {
	*Cubic
	StrongWolfeConditions
}

type LinesearchProblem struct {
	// Add in MISO problem stuff here
}

func Linesearch(l Linesearchable, direction []float64, Loc *LocationFloatSlice, Opt *ObjectiveFloatSlice, Grad *GradientFloatSlice) {
	newX := make([]float64, len(x0.Curr()))
	linesearcher = l.Linesearcher
	linesearcher.Loc.Init = 0
	linesearcher.Opt.Init = l.Opt.Curr

	stepDirection := smatrix.UnitVector(initialSearchVector)
	initGradProjection := smatrix.DotVector(stepDirection, g0.Curr())

	linesearcher.Grad.Init = initGradProjection

}
