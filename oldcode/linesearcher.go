package gofunopter

import (
	//"fmt"
	"github.com/btracey/smatrix"
	"math"
)

// TODO: Add in error checking for positive initial gradient? Maybe should be a panic
// because it shouldn't ever occur
// TODO: Make AddToHist thread safe so multiple linesearches could be called simulaneously

type WolfeConditioner interface {
	Converged() Convergence
	FunConst() float64
	GradConst() float64
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

func (s *StrongWolfeConditions) FunConst() float64 {
	return s.funConst
}

func (s *StrongWolfeConditions) GradConst() float64 {
	return s.gradConst
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

// TODO: Change to make it okay to use any SISO method, not just gradient based
type LinesearchFun struct {
	Linesearch       Linesearcher
	Direction        []float64
	InitSearchVector []float64
	InitLoc          []float64
	CurrLoc          []float64
	CurrGrad         []float64
}

func (l *LinesearchFun) Eval(step float64) (float64, float64, error) {
	loc := make([]float64, len(l.InitLoc))

	for i, val := range l.InitSearchVector {
		loc[i] = val*step + l.InitLoc[i]
	}
	l.CurrLoc = loc
	l.Linesearch.Loc().AddToHist(loc)
	f, g, err := l.Linesearch.Fun().Eval(loc)
	l.CurrGrad = g
	directionalG := smatrix.DotVector(l.Direction, g)
	l.Linesearch.WolfeConditions().SetCurr(f, directionalG, step)
	return f, directionalG, err
}

func (l *LinesearchFun) Converged() Convergence {
	// Set the function and gradient values for the line searcher
	return l.Linesearch.WolfeConditions().Converged()
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
func Linesearch(linesearcher Linesearcher, initSearchVector []float64, initLoc []float64, initObj float64, initGrad []float64) (*LinesearchResult, error) {
	//func (s *SeqLinesearch) Linesearch(linesearcher Linsearchable, direction []float64, initLoc []float64, initObj float64, initGrad []float64) {
	//newX := make([]float64, len(x0.Curr()))
	//fmt.Println(initGradProjection)

	// Need to add a Reset method to the linesearcher so the bulk can be reused (right now step size is being saved)
	sisoGradBased := linesearcher.LinesearchMethod()
	sisoGradBased.Loc().SetInit(0)
	sisoGradBased.Obj().SetInit(initObj)

	//sisoGradBased.FunEvals().SetMax(100)

	// This is wrong, shouldn't be renormalizing I believe.
	//direction := smatrix.UnitVector(initSearchVector)
	//initGradProjection := smatrix.DotVector(direction, initGrad)
	direction := initSearchVector
	initGradProjection := smatrix.DotVector(direction, initGrad)

	//fmt.Println("Init Grad", initGrad)
	//fmt.Println("Init Grad norm", smatrix.VectorTwoNorm(initGrad))
	//fmt.Println("Direction", direction)
	//fmt.Println("Direction norm", smatrix.VectorTwoNorm(direction))
	//fmt.Println("InitGradProjection", initGradProjection)

	// Set wolfe constants
	linesearcher.WolfeConditions().SetInit(initObj, initGradProjection)
	linesearcher.WolfeConditions().SetCurr(initObj, initGradProjection, 1.0)

	sisoGradBased.Grad().SetInit(initGradProjection)
	fun := &LinesearchFun{
		Linesearch:       linesearcher,
		Direction:        direction,
		InitSearchVector: initSearchVector,
		InitLoc:          initLoc,
	}

	// Maybe it isn't resetting any of the other counters and such.
	// Upon calling result it should all reset itself so can use again

	sisoGradBased.SetFun(fun)
	sisoGradBased.SetDisp(false)
	convergence, err := Optimize(sisoGradBased)

	r := &LinesearchResult{
		Loc:       fun.CurrLoc,
		Obj:       sisoGradBased.Obj().Opt(),
		Grad:      fun.CurrGrad,
		Step:      sisoGradBased.Loc().Opt(),
		NFunEvals: sisoGradBased.FunEvals().Opt(),
	}

	if err != nil {
		return r, &LinesearchError{Err: err}
	}
	// Need to pass on the strings
	_, ok := convergence.(WolfeConvergence)
	if !ok {
		// Check if the wolfe conditions are met anyway
		c := fun.Converged()
		if c != nil {
			// Conditions met, no problem
			return r, nil
		}
		return r, &LinesearchConvergenceFalure{Conv: convergence, Loc: initLoc, InitStep: initSearchVector}
	}
	//fmt.Println("Finished linesearch")
	//fmt.Println()
	return r, nil
}
