package multivariate

import (
	"github.com/btracey/gofunopter/common/linesearch"
	"github.com/btracey/gofunopter/common/multi"
	"github.com/btracey/gofunopter/common/optimize"
	"github.com/btracey/gofunopter/common/status"
	"github.com/btracey/gofunopter/common/uni"
	"github.com/btracey/gofunopter/univariate"

	"errors"
	//"fmt"
	"github.com/gonum/floats"
)

type Lbfgs struct {
	// Basic structures for the state of the optimizer
	step *uni.BoundedStep

	// Tunable Parameters
	LinesearchMethod   linesearch.LinesearchMethod
	LinesearchSettings *univariate.UniGradSettings
	NumStore           int // How many gradients to store
	Wolfe              linesearch.WolfeConditioner

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
		step: uni.NewBoundedStep(),

		LinesearchMethod:   univariate.NewCubic(),
		LinesearchSettings: univariate.NewUniGradSettings(),
		NumStore:           30,
		Wolfe:              &linesearch.StrongWolfeConditions{},
	}
	l.Wolfe.SetFunConst(0)
	l.Wolfe.SetGradConst(0.9)
	l.LinesearchSettings.MaximumFunctionEvaluations = 100
	l.LinesearchSettings.Display = false
	l.LinesearchSettings.GradientAbsoluteTolerance = 0 // Force convergence from wolfe conditions
	return l
}

func (lbfgs *Lbfgs) UnivariateSettings() *univariate.UniGradSettings {
	return lbfgs.LinesearchSettings
}

func (lbfgs *Lbfgs) SetResult() {
	optimize.SetResult(lbfgs.step)
}

func (lbfgs *Lbfgs) Initialize(loc *multi.Location, obj *uni.Objective, grad *multi.Gradient) error {
	lbfgs.nDim = len(loc.Init())

	// Now initialize the three to set the initial location to the current location
	err := optimize.Initialize(lbfgs.step)
	if err != nil {
		return errors.New("lbfgs: error initializing: " + err.Error())
	}

	// Initialize rest of memory

	// Replace this with overwriting?
	lbfgs.q = make([]float64, lbfgs.nDim)
	lbfgs.a = make([]float64, lbfgs.NumStore)
	lbfgs.b = make([]float64, lbfgs.NumStore)
	lbfgs.sHist = make([][]float64, lbfgs.NumStore)
	lbfgs.yHist = make([][]float64, lbfgs.NumStore)
	lbfgs.rhoHist = make([]float64, lbfgs.NumStore)

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

func (lbfgs *Lbfgs) Iterate(loc *multi.Location, obj *uni.Objective, grad *multi.Gradient, fun optimize.MultiObjGrad) (status.Status, error) {
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
	for i, val := range grad.Curr() {
		q[i] = val
	}
	for i := counter - 1; i >= 0; i-- {
		a[i] = rhoHist[i] * floats.Dot(sHist[i], q)
		copy(tmp, yHist[i])
		floats.Scale(a[i], tmp)
		floats.Sub(q, tmp)
	}
	for i := lbfgs.NumStore - 1; i >= counter; i-- {
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
	for i := counter; i < lbfgs.NumStore; i++ {
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

	// Perform line search -- need to find some way to implement this, especially bookkeeping function values
	linesearchResult, err := linesearch.Linesearch(fun, lbfgs.LinesearchMethod, lbfgs.LinesearchSettings, lbfgs.Wolfe, p_k, loc.Curr(), obj.Curr(), grad.Curr())

	// In the future add a check to switch to a different linesearcher?
	if err != nil {
		return status.LinesearchFailure, err
	}
	x_kp1 := linesearchResult.Loc
	f_kp1 := linesearchResult.Obj
	g_kp1 := linesearchResult.Grad
	alpha_k := linesearchResult.Step

	// Update hessian estimate
	copy(s_k, p_k)
	floats.Scale(alpha_k, s_k)

	copy(y_k, g_kp1)
	floats.Sub(y_k, grad.Curr())
	skDotYk := floats.Dot(s_k, y_k)

	// Bookkeep the results
	stepSize := alpha_k * normP_k
	lbfgs.step.AddToHist(stepSize)
	lbfgs.step.SetCurr(stepSize)
	loc.SetCurr(x_kp1)
	//lbfgs.loc.AddToHist(x_kp1)

	//fmt.Println(lbfgs.loc.GetHist())
	obj.SetCurr(f_kp1)
	grad.SetCurr(g_kp1)

	copy(sHist[counter], s_k)
	copy(yHist[counter], y_k)
	rhoHist[counter] = 1 / skDotYk

	lbfgs.gamma_k = skDotYk / floats.Dot(y_k, y_k)

	lbfgs.counter += 1
	if lbfgs.counter == lbfgs.NumStore {
		lbfgs.counter = 0
	}
	return status.Continue, nil
}
