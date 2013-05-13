package gofunopter

import "math"

// TODO: Add in error checking for positive initial gradient? Maybe should be a panic
// because it shouldn't ever occur
// TODO: Make AddToHist thread safe so multiple linesearches could be called simulaneously

type WolfeConditioner interface {
	IsConverged(initObj, initGrad, currObj, currGrad, step float64) bool
	SetFunConst(funConst float64)
	SetGradConst(gradConst float64)
	SetInit(initObj, initGrad float64)
	SetCurr(currObj, currGrad, step float64)
}

type WolfeConvergence struct {
	Str string
}

func (w WolfeConvergence) ConvergenceType() string {
	return w.Str
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

func (w *WeakWolfeConditions) SetInit(initObj, initGrad float64) {
	w.initObj = initObj
	w.initGrad = initGrad
	w.step = math.Inf(1)
}

func (w *WeakWolfeConditions) SetCurr(currObj, currGrad, currStep float64) {
	w.currObj = currObj
	w.currGrad = currGrad
	w.step = currStep
}

//func (s *WeakWolfeConditions) WolfeConditionsMet(obj, directionalderivative, step float64) bool {
func (w *WeakWolfeConditions) Converged() Convergence {
	if w.currObj >= w.initObj+w.funConst*w.step*w.currGrad {
		return nil
	}
	if w.currGrad <= w.gradConst*w.initGrad {
		return nil
	}
	return WolfeConvergence{"Weak Wolfe conditions met"}
}

func (w *WeakWolfeConditions) SetFunConst(val float64) {
	w.funConst = val
}

func (w *WeakWolfeConditions) SetGradConst(val float64) {
	w.gradConst = val
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

func (s *StrongWolfeConditions) SetInit(initObj, initGrad float64) {
	s.initObj = initObj
	s.initGrad = initGrad
	s.step = math.Inf(1)
}

func (s *StrongWolfeConditions) SetCurr(currObj, currGrad, currStep float64) {
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

func (s *StrongWolfeConditions) Converged() Convergence {
	if s.currObj >= s.initObj+s.funConst*s.step*s.currGrad {
		return nil
	}
	if math.Abs(s.currGrad) >= s.gradConst*math.Abs(s.initGrad) {
		return nil
	}
	return WolfeConvergence{"Strong Wolfe conditions met"}
}

// Maybe everything should be through interfaces to make everything easier
// to set. Make OptFloat an interface. Also probably makes it easier to customize.
// Harder to save possibly, but not hard to just save the float values

type Linesearcher interface {
	MisoGradBasedOptimizer
	LinesearchMethod() SisoGradBasedOptimizer
	WolfeConditions() WolfeConditioner
}

/*
type Linesearcher interface {
	Method() SISOGradBasedOptimizer
	SetMethod(SISOGradBasedOptimizer)
	Wolfe() WolfeConditioner
	SetWolfe(WolfeConditioner)
	Linesearch(linesearcher Linsearchable, direction []float64, initLoc []float64, initObj float64, initGrad []float64) (Convergence, err)
}
*/

type LinesearchResult struct {
	Loc       []float64
	Obj       float64
	Grad      []float64
	Step      float64
	NFunEvals int
}

// Should this be changed to be an interface, so you could, for example, evaluate several points
// in the line search in parallel?
// Could also define a linesearch interface, and then have sequential, parallel, etc.

type LinesearchFun struct {
	Linesearch Linesearcher
	//MisoProb   MisoGradBasedProblem
	Direction []float64
	Loc       []float64
	InitLoc   []float64
}

func (l *LinesearchFun) Eval(step float64) error {
	for i, val := range l.Direction {
		l.Loc[i] = val*step + l.InitLoc[i]
	}
	l.Linesearch.Location().AddToHist(l.Loc)
	err := l.Linesearch.Function().Eval(l.Loc)
	return err
}

func (l *LinesearchFun) Obj() float64 {
	o := l.Linesearch.Function().Obj()
	l.Linesearch.Objective().AddToHist(o)
	return o
}

func (l *LinesearchFun) Grad() float64 {
	g := l.MisoProb.Function().Grad()
	l.Linesearch.Grad().AddToHist(g)
	return smatrix.DotVector(l.Direction, g)
}

func (l *LinesearchFun) Converged() string {
	// Set the function and gradient values for the line searcher
	l.Linesearch.WolfeConditions().Set(l.Linesearch.Fun().Obj(), l.Linesearch.Fun().Grad())
	return l.Linesearch.Wolfe().Converged()
}

type LineSearchSuccess BasicConvergence
type LineSearchFailure BasicConvergence

// Move the SISO into here
/*
type SeqLinesearch struct {
	Siso  *SISOGradBasedOptimizer
	Wolfe *WolfeConditioner
}

func DefaultSequentialLinesearch() *SeqLinesearch {
	return &SeqLinesearch{
		Siso: DefaultCubic(),
		Wolfe: &StrongWolfeConditions{
			FunConst:  1.0,
			GradConst: 0.0,
		},
	}
}
*/
func Linesearch(MisoGradBasedOptimizer, direction []float64, initLoc []float64, initObj float64, initGrad []float64) {
	//func (s *SeqLinesearch) Linesearch(linesearcher Linsearchable, direction []float64, initLoc []float64, initObj float64, initGrad []float64) {
	newX := make([]float64, len(x0.Curr()))

	sisoGradBased = linesearcher.LinesearchMethod()
	sisoGradBased.Loc().SetInit(0)
	sisoGradBased.Opt().SetInit(initObj)
	stepDirection := smatrix.UnitVector(direction)
	initGradProjection := smatrix.DotVector(stepDirection, initGrad)

	sisoGradBased.Grad().SetInit(initGradProjection)
	fun := &LinesearchFun{
		Miso:      linesearcher,
		Direction: direction,
		InitLoc:   initLoc,
	}
	sisoGradBased.SetFun(fun)
	convergence, err := Optimize(sisoGradBased)

	if err != nil {
		return OptimizerError, err
	}
	// Need to pass on the strings
	_, ok := convergence.(WolfeConvergence)
	if !ok {
		return &LinesearchFailure{}, nil
	}
	return &LinesearchSuccess{}, nil
}
