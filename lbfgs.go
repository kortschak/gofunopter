package gofunopter

// Change common so that it doesn't implement the optimizer stuff (force optimizer to write it)

import (
	"github.com/btracey/gofunopter/common"
	"github.com/btracey/gofunopter/convergence"
	"github.com/btracey/gofunopter/display"
	"github.com/btracey/gofunopter/linesearch"
	"github.com/btracey/gofunopter/multivariate"
	"github.com/btracey/gofunopter/optimize"
	"github.com/btracey/gofunopter/univariate"

	"errors"
	"fmt"
	"github.com/gonum/floats"
	"math"
)

type Lbfgs struct {
	// Needed for general optimizer
	*common.OptCommon
	*display.Display
	disp bool

	// Basic structures for the state of the optimizer
	loc  *multivariate.Location
	obj  *univariate.Objective
	grad *multivariate.Gradient
	step *univariate.BoundedStep

	// User defined function
	fun optimize.MisoGrad

	// Tunable Parameters
	linesearchMethod   linesearch.Method
	nStore             int // How many gradients to store
	wolfe              linesearch.WolfeConditioner
	maxLineSearchEvals int

	// Other needed variables
	gamma_k float64
	q       []float64
	a       []float64
	b       []float64
	counter int
	sHist   [][]float64
	yHist   [][]float64
	rhoHist []float64
	nDim    int
	p_k     []float64
	s_k     []float64
	y_k     []float64
	z       []float64
	tmp     []float64 // For storing intermediate values
}

func NewLbfgs() *Lbfgs {
	l := &Lbfgs{
		OptCommon: common.NewOptCommon(),
		Display:   display.NewDisplay(),
		disp:      true,

		loc:  multivariate.NewLocation(),
		obj:  univariate.NewObjective(),
		grad: multivariate.NewGradient(),
		step: univariate.NewBoundedStep(),

		linesearchMethod:   NewCubic(),
		nStore:             30,
		wolfe:              &linesearch.StrongWolfeConditions{},
		maxLineSearchEvals: 100,
	}
	l.wolfe.SetFunConst(0)
	l.wolfe.SetGradConst(0.9)
	l.linesearchMethod.SetDisp(false)
	l.linesearchMethod.FunEvals().SetMax(l.maxLineSearchEvals)
	l.linesearchMethod.Grad().SetAbsTol(0) // Turn it off so that it must converge with Wolfe
	l.linesearchMethod.Grad().SetRelTol(0) // Turn it off so that it must converge with Wolfe
	return l
}

func (lbfgs *Lbfgs) Optimize(fun optimize.MisoGrad) (convergence.C, error) {
	lbfgs.fun = fun
	return optimize.OptimizeOpter(lbfgs, fun)
}

func (lbfgs *Lbfgs) Wolfe() linesearch.WolfeConditioner {
	return lbfgs.wolfe
}

func (lbfgs *Lbfgs) SetWolfe(w linesearch.WolfeConditioner) {
	lbfgs.wolfe = w
}

func (lbfgs *Lbfgs) SetLinesearchMethod(method linesearch.Method) {
	lbfgs.linesearchMethod = method
}

func (lbfgs *Lbfgs) SetStoredIterations(k int) {
	lbfgs.nStore = k
}

func (lbfgs *Lbfgs) Disp() bool {
	return lbfgs.disp
}

func (lbfgs *Lbfgs) SetDisp(b bool) {
	lbfgs.disp = b
}

func (lbfgs *Lbfgs) Loc() *multivariate.Location {
	return lbfgs.loc
}

func (lbfgs *Lbfgs) Obj() *univariate.Objective {
	return lbfgs.obj
}

func (lbfgs *Lbfgs) Grad() *multivariate.Gradient {
	return lbfgs.grad
}

func (lbfgs *Lbfgs) Fun() optimize.MisoGrad {
	return lbfgs.fun
}

func (lbfgs *Lbfgs) Converged() convergence.C {
	return convergence.CheckConvergence(lbfgs.obj, lbfgs.grad)
	//if c != nil {
	//	fmt.Println("Lbfgs converged")
	//}
	//return c
}

func (lbfgs *Lbfgs) AddToDisplay(d []*display.Struct) []*display.Struct {
	if lbfgs.disp {
		d = display.AddToDisplay(d, lbfgs.loc, lbfgs.obj, lbfgs.grad, lbfgs.step)
	}
	return d
}

func (lbfgs *Lbfgs) SetResult() {
	optimize.SetResult(lbfgs.loc, lbfgs.obj, lbfgs.grad)
}

func (lbfgs *Lbfgs) Initialize() (err error) {
	// Get the initial function value
	initLoc := lbfgs.loc.Init()
	if initLoc == nil {
		return errors.New("lbfgs: initial location is nil")
	}

	lbfgs.nDim = len(initLoc)

	initObj := lbfgs.obj.Init()
	initGrad := lbfgs.grad.Init()

	// The initial values need to both be NaN/nil or both not nan/nil
	if math.IsNaN(initObj) {
		if !(initGrad == nil) {
			return errors.New("lbfgs: initial function value and gradient must either both be set or neither set")
		}
		// Both nan, so compute the initial fuction value and gradient
		initObj, initGrad, err := lbfgs.fun.Eval(initLoc)
		if err != nil {
			return errors.New("lbfgs: error calling function during optimization")
		}
		lbfgs.obj.SetInit(initObj)
		lbfgs.grad.SetInit(initGrad)
	} else {
		if initGrad == nil {
			return errors.New("lbfgs: initial function value and gradient must either both be set or neither set")
		}
	}

	// Now initialize the three to set the initial location to the current location
	err = optimize.Initialize(lbfgs.loc, lbfgs.obj, lbfgs.grad, lbfgs.step)
	if err != nil {
		return errors.New("lbfgs: error initializing: " + err.Error())
	}

	// Initialize rest of memory

	// Replace this with overwriting?
	lbfgs.q = make([]float64, lbfgs.nDim)
	lbfgs.a = make([]float64, lbfgs.nStore)
	lbfgs.b = make([]float64, lbfgs.nStore)
	lbfgs.sHist = make([][]float64, lbfgs.nStore)
	lbfgs.yHist = make([][]float64, lbfgs.nStore)
	lbfgs.rhoHist = make([]float64, lbfgs.nStore)

	for i := range lbfgs.sHist {
		lbfgs.sHist[i] = make([]float64, lbfgs.nDim)
		lbfgs.yHist[i] = make([]float64, lbfgs.nDim)
	}

	lbfgs.gamma_k = 1.0

	lbfgs.tmp = make([]float64, lbfgs.nDim)
	lbfgs.p_k = make([]float64, lbfgs.nDim)
	lbfgs.s_k = make([]float64, lbfgs.nDim)
	lbfgs.y_k = make([]float64, lbfgs.nDim)
	lbfgs.z = make([]float64, lbfgs.nDim)
	return nil
}

func (lbfgs *Lbfgs) Iterate() (int, error) {

	nFunEvals := 0
	counter := lbfgs.counter
	q := lbfgs.q
	a := lbfgs.a
	b := lbfgs.b
	rhoHist := lbfgs.rhoHist
	sHist := lbfgs.sHist
	yHist := lbfgs.yHist
	gamma_k := lbfgs.gamma_k
	tmp := lbfgs.tmp
	p_k := lbfgs.p_k
	s_k := lbfgs.s_k
	y_k := lbfgs.y_k
	z := lbfgs.z

	// Calculate search direction
	for i, val := range lbfgs.grad.Curr() {
		q[i] = val
	}
	for i := counter - 1; i >= 0; i-- {
		a[i] = rhoHist[i] * floats.Dot(sHist[i], q)
		copy(tmp, yHist[i])
		floats.Scale(a[i], tmp)
		floats.Sub(q, tmp)
	}
	for i := lbfgs.nStore - 1; i >= counter; i-- {
		a[i] = rhoHist[i] * floats.Dot(sHist[i], q)
		copy(tmp, yHist[i])
		floats.Scale(a[i], tmp)
		//fmt.Println(q)
		//fmt.Println(tmp)
		floats.Sub(q, tmp)
	}

	// Assume H_0 is the identity times gamma_k
	copy(z, q)
	floats.Scale(gamma_k, z)
	// Second loop for update, going oldest to newest
	for i := counter; i < lbfgs.nStore; i++ {
		b[i] = rhoHist[i] * floats.Dot(yHist[i], z)
		copy(tmp, sHist[i])
		floats.Scale(a[i]-b[i], tmp)
		floats.Add(z, tmp)
	}
	for i := 0; i < counter; i++ {
		b[i] = rhoHist[i] * floats.Dot(yHist[i], z)
		copy(tmp, sHist[i])
		floats.Scale(a[i]-b[i], tmp)
		floats.Add(z, tmp)
	}

	lbfgs.a = a
	lbfgs.b = b

	copy(p_k, z)
	floats.Scale(-1, p_k)
	normP_k := floats.Norm(p_k, 2)

	if lbfgs.linesearchMethod.Disp() {
		fmt.Println("Starting linesearch")
	}
	// Perform line search -- need to find some way to implement this, especially bookkeeping function values
	linesearchResult, err := linesearch.Linesearch(lbfgs, lbfgs.linesearchMethod, lbfgs.wolfe, p_k, lbfgs.loc.Curr(), lbfgs.obj.Curr(), lbfgs.grad.Curr())
	if lbfgs.linesearchMethod.Disp() {
		fmt.Println("Done linesearch")
	}
	// In the future add a check to switch to a different linesearcher?
	nFunEvals += linesearchResult.NFunEvals
	if err != nil {
		return nFunEvals, err
	}

	x_kp1 := linesearchResult.Loc
	f_kp1 := linesearchResult.Obj
	g_kp1 := linesearchResult.Grad
	alpha_k := linesearchResult.Step

	// Update hessian estimate
	copy(s_k, p_k)
	floats.Scale(alpha_k, s_k)

	copy(y_k, g_kp1)
	floats.Sub(y_k, lbfgs.grad.Curr())
	skDotYk := floats.Dot(s_k, y_k)

	// Bookkeep the results
	stepSize := alpha_k * normP_k
	lbfgs.step.AddToHist(stepSize)
	lbfgs.step.SetCurr(stepSize)
	lbfgs.loc.SetCurr(x_kp1)
	//lbfgs.loc.AddToHist(x_kp1)

	//fmt.Println(lbfgs.loc.GetHist())
	lbfgs.obj.SetCurr(f_kp1)
	lbfgs.grad.SetCurr(g_kp1)

	copy(sHist[counter], s_k)
	copy(yHist[counter], y_k)
	rhoHist[counter] = 1 / skDotYk

	lbfgs.gamma_k = skDotYk / floats.Dot(y_k, y_k)

	lbfgs.counter += 1
	if lbfgs.counter == lbfgs.nStore {
		lbfgs.counter = 0
	}
	return nFunEvals, nil
}
