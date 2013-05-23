package gofunopter

import (
	"fmt"
	"github.com/btracey/smatrix"
	"math"
)

func DefaultLbfgs() *Lbfgs {
	l := &Lbfgs{
		loc:    DefaultLocationFloatSlice(),
		obj:    DefaultObjectiveFloat(),
		grad:   DefaultGradientFloatSlice(),
		step:   DefaultBoundedStepFloat(),
		Common: DefaultCommon(),

		wolfe:            &StrongWolfeConditions{},
		linesearchMethod: DefaultCubic(),
		nStore:           30,
	}
	l.wolfe.SetFunConst(0.0)
	l.wolfe.SetGradConst(0.5)

	// Turn off tolerances so the line search has to meet the wolfe conditions to move on
	l.step.SetAbsTol(0)
	l.step.SetRelTol(0)
	l.linesearchMethod.Grad().SetAbsTol(0)
	l.linesearchMethod.Grad().SetRelTol(0)

	// Set a max fun evals in case something gets stuck
	l.linesearchMethod.FunEvals().SetMax(100)
	SetDisplayMethods(l)
	return l
}

// IDEA: Have an Eval() and then a Fun() Grad() etc. This way, the optimizeoptimizer
// routine can do all the things for the function taking the burden off of the
// optimization routine. Also could add the Commoner part
// Other optimizers that can just use constraint evaluations (or whatever) could have
// special interfaces for their problems if need be.
// This way, the optimizer just needs to take care of its own part
// Also, add in an InitGuesser check (this also makes it really easy to do random restarts)

type Lbfgs struct {
	loc  OptFloatSlice
	obj  OptTolFloat
	grad OptTolFloatSlice
	step BoundedOptFloat
	*Common
	fun MisoGradBasedProblem

	// Tunable Parameters
	linesearchMethod SisoGradBasedOptimizer
	nStore           int // How many gradients to store
	wolfe            WolfeConditioner

	// Other needed variables
	gamma_k float64
	q       []float64
	a       []float64
	b       []float64
	counter int
	sHist   [][]float64
	yHist   [][]float64
	rhoHist []float64
}

func (lbfgs *Lbfgs) Optimize(fun MisoGradBasedProblem) {
	lbfgs.fun = fun
	Optimize(lbfgs)
}

func (lbfgs *Lbfgs) WolfeConditions() WolfeConditioner {
	return lbfgs.wolfe
}

func (lbfgs *Lbfgs) SetNumStore(i int) {
	lbfgs.nStore = i
}

func (c *Lbfgs) Loc() OptFloatSlice {
	return c.loc
}

func (c *Lbfgs) Obj() OptTolFloat {
	return c.obj
}

func (c *Lbfgs) Grad() OptTolFloatSlice {
	return c.grad
}

func (c *Lbfgs) Fun() MisoGradBasedProblem {
	return c.fun
}

func (c *Lbfgs) Step() BoundedOptFloat {
	return c.step
}

func (c *Lbfgs) SetFun(misoGradBasedProblem MisoGradBasedProblem) {
	c.fun = misoGradBasedProblem
}

func (c *Lbfgs) SetFunc(fun MisoGradBasedProblem) {
	c.fun = fun
}

func (lbfgs *Lbfgs) LinesearchMethod() SisoGradBasedOptimizer {
	return lbfgs.linesearchMethod
}

func (lbfgs *Lbfgs) SetLinesearch(linesearchMethod SisoGradBasedOptimizer) {
	lbfgs.linesearchMethod = linesearchMethod
}

func (lbfgs *Lbfgs) Initialize() error {
	iger, ok := lbfgs.fun.(InitGuesserFloatSlice)
	if ok {
		lbfgs.loc.SetInit(iger.InitGuess())
	}
	if lbfgs.loc.Init() == nil {
		return fmt.Errorf("Initial location must be provided. (Set using lbfgs.Loc().SetInit(val) ), or the function must be an InitGuesserFloatSlice")
	}

	lbfgs.Common.Initialize()
	lbfgs.loc.Initialize()
	nDim := len(lbfgs.loc.Init())

	s, ok := lbfgs.fun.(Initializer)
	if ok {
		err := s.Initialize()
		if err != nil {
			return fmt.Errorf("Error initializing user-defined function")
		}
	}

	// If the initial value isn't set, evaluate the function to get the initial value and gradient
	if math.IsNaN(lbfgs.obj.Init()) {
		f, g, err := lbfgs.fun.Eval(lbfgs.loc.Init())
		if err != nil {
			return fmt.Errorf("Error evaulating function at initial value: " + err.Error())
		}
		lbfgs.FunEvals().Add(1)
		lbfgs.obj.SetInit(f)
		lbfgs.grad.SetInit(g)
	}
	lbfgs.obj.Initialize()
	lbfgs.grad.Initialize()
	lbfgs.step.Initialize()

	lbfgs.q = make([]float64, nDim)
	lbfgs.a = make([]float64, lbfgs.nStore)
	lbfgs.b = make([]float64, lbfgs.nStore)
	lbfgs.sHist = make([][]float64, lbfgs.nStore)
	lbfgs.yHist = make([][]float64, lbfgs.nStore)
	lbfgs.rhoHist = make([]float64, lbfgs.nStore)

	for i := range lbfgs.sHist {
		lbfgs.sHist[i] = make([]float64, nDim)
		lbfgs.yHist[i] = make([]float64, nDim)
	}

	lbfgs.gamma_k = 1.0
	return nil
}

func (lbfgs *Lbfgs) Converged() Convergence {
	conv := Converged(lbfgs.obj, lbfgs.grad, lbfgs.step)
	if conv != nil {
		return conv
	}
	s, ok := lbfgs.fun.(Converger)
	if ok {
		conv = s.Converged()
		if conv != nil {
			return conv
		}
	}
	return nil
}

func (c *Lbfgs) AppendHeadings(headings []string) []string {
	headings = AppendHeadings(headings, c.Common, c.loc, c.obj, c.grad, c.step)
	s, ok := c.fun.(Displayer)
	if ok {
		headings = AppendHeadings(headings, s)
	}
	return headings
}

func (c *Lbfgs) AppendValues(values []interface{}) []interface{} {
	values = AppendValues(values, c.Common, c.loc, c.obj, c.grad, c.step)
	s, ok := c.fun.(Displayer)
	if ok {
		values = AppendValues(values, s)
	}
	return values
}

func (lbfgs *Lbfgs) Iterate() error {
	// TODO: Should there be an iterate for loc et al?
	err := lbfgs.Common.Iterate()
	if err != nil {
		fmt.Println("error")
		return err
	}

	counter := lbfgs.counter
	q := lbfgs.q
	a := lbfgs.a
	b := lbfgs.b
	rhoHist := lbfgs.rhoHist
	sHist := lbfgs.sHist
	yHist := lbfgs.yHist
	gamma_k := lbfgs.gamma_k

	// Calculate search direction
	for i, val := range lbfgs.grad.Curr() {
		q[i] = val
	}
	for i := counter - 1; i >= 0; i-- {
		a[i] = rhoHist[i] * smatrix.DotVector(sHist[i], q)
		smatrix.SubtractVectorInPlace(q, smatrix.ScaleVector(yHist[i], a[i]))
	}
	for i := lbfgs.nStore - 1; i >= counter; i-- {
		a[i] = rhoHist[i] * smatrix.DotVector(sHist[i], q)
		smatrix.SubtractVectorInPlace(q, smatrix.ScaleVector(yHist[i], a[i]))
	}

	// Assume H_0 is the identity times gamma_k
	z := smatrix.ScaleVector(q, gamma_k)
	// Second loop for update, going oldest to newest
	for i := counter; i < lbfgs.nStore; i++ {
		b[i] = rhoHist[i] * smatrix.DotVector(yHist[i], z)
		smatrix.AddVectorInPlace(z, smatrix.ScaleVector(sHist[i], a[i]-b[i]))
	}
	for i := 0; i < counter; i++ {
		b[i] = rhoHist[i] * smatrix.DotVector(yHist[i], z)
		smatrix.AddVectorInPlace(z, smatrix.ScaleVector(sHist[i], a[i]-b[i]))
	}

	lbfgs.a = a
	lbfgs.b = b

	p_k := smatrix.ScaleVector(z, -1)
	normP_k := smatrix.VectorTwoNorm(p_k)

	// Perform line search -- need to find some way to implement this, especially bookkeeping function values

	//x_kp1, f_kp1, g_kp1, alpha_k, nFunEval, err := lbfgs.line.Linesearch(lbfgs.fun, lbfgs.x, lbfgs.f, lbfgs.g, p_k)
	linesearchResult, err := Linesearch(lbfgs, p_k, lbfgs.loc.Curr(), lbfgs.obj.Curr(), lbfgs.grad.Curr())
	// In the future add a check to switch to a different linesearcher?
	if err != nil {
		return err
	}

	lbfgs.FunEvals().Add(linesearchResult.NFunEvals)
	x_kp1 := linesearchResult.Loc
	f_kp1 := linesearchResult.Obj
	g_kp1 := linesearchResult.Grad
	alpha_k := linesearchResult.Step

	// Update hessian estimate

	sk := smatrix.ScaleVector(p_k, alpha_k)
	yk := smatrix.SubtractVector(g_kp1, lbfgs.grad.Curr())
	skDotYk := smatrix.DotVector(sk, yk)

	// Bookkeep the results
	stepSize := alpha_k * normP_k
	lbfgs.step.AddToHist(stepSize)
	lbfgs.step.SetCurr(stepSize)
	lbfgs.loc.SetCurr(x_kp1)
	//lbfgs.loc.AddToHist(x_kp1)

	//fmt.Println(lbfgs.loc.GetHist())
	lbfgs.obj.SetCurr(f_kp1)
	lbfgs.grad.SetCurr(g_kp1)

	sHist[counter] = sk
	yHist[counter] = yk
	rhoHist[counter] = 1 / skDotYk

	lbfgs.gamma_k = skDotYk / smatrix.DotVector(yk, yk)

	lbfgs.counter += 1
	if lbfgs.counter == lbfgs.nStore {
		lbfgs.counter = 0
	}
	return nil
}
