package multivariate

import (
	"github.com/btracey/gofunopter/common/convergence"
	//"github.com/btracey/gofunopter/optimize"

	"github.com/gonum/floats"
	"math"
	"testing"

	"fmt"
	"math/rand"
	"strconv"
)

var MISO_TOLERANCE float64 = 1E-6

type Rosenbrock struct {
	nDim int
}

func (r *Rosenbrock) Eval(x []float64) (sum float64, deriv []float64, err error) {
	sum = 0
	deriv = make([]float64, len(x))
	for i := 0; i < len(x)-1; i++ {
		sum += math.Pow(1-x[i], 2) + 100*math.Pow(x[i+1]-math.Pow(x[i], 2), 2)
	}
	for i := 0; i < len(x)-1; i++ {
		deriv[i] += -1 * 2 * (1 - x[i])
		deriv[i] += 2 * 100 * (x[i+1] - math.Pow(x[i], 2)) * (-2 * x[i])
	}
	for i := 1; i < len(x); i++ {
		deriv[i] += 2 * 100 * (x[i] - math.Pow(x[i-1], 2))
	}
	return sum, deriv, nil
}

func (r *Rosenbrock) OptVal() float64 {
	return 0
}

func (r *Rosenbrock) OptLoc() []float64 {
	ans := make([]float64, r.nDim)
	floats.AddConst(1, ans)
	return ans
}

type MisoGradTestFunction interface {
	MultiGradFun
	OptVal() float64
	OptLoc() []float64
}

type MisoGradTest struct {
	MisoGradTestFunction
	InitLoc []float64
	name    string
}

func RandRosen(nDim int, low, high float64) MisoGradTest {
	rosen := MisoGradTest{
		MisoGradTestFunction: &Rosenbrock{nDim: nDim},
		name:                 "TenDRosenbrock rand, nDim = " + strconv.Itoa(nDim),
		InitLoc:              make([]float64, nDim),
	}
	for i := range rosen.InitLoc {
		rosen.InitLoc[i] = rand.Float64()*(high-low) + low
	}
	return rosen
}

func MisoGradFunctions() []MisoGradTest {
	misotest := make([]MisoGradTest, 0)
	misotest = append(misotest, RandRosen(10, -10, 10))
	misotest = append(misotest, RandRosen(50, -2, 2))
	misotest = append(misotest, RandRosen(50, -100, 100))
	return misotest
}

func MisoGradBasedTest(t *testing.T, opter MultiGradOptimizer) {
	funcs := MisoGradFunctions()
	for _, fun := range funcs {
		// Run it once until very converged

		settings := NewMultiGradSettings()
		settings.GradientAbsoluteTolerance = 1e-14
		settings.MaximumFunctionEvaluations = 1000

		//opter.Loc().SetInit(fun.InitLoc)
		//opter.Grad().SetAbsTol(1E-14)
		//opter.GetOptCommon().FunEvals().SetMax(1000)
		//opter.GetDisplay().SetValueDisplayInterval(0)
		fmt.Println("Is misograd_test, starting optimizer")
		//optVal, optLoc, c, err := opter.Optimize(fun, fun.InitLoc)
		optVal, optLoc, result, err := OptimizeGrad(fun, fun.InitLoc, settings, opter)
		c := result.Convergence
		if err != nil {
			t.Errorf("Error during optimization for function " + fun.name + ": " + err.Error())
			return
		}
		if c == nil {
			t.Errorf("Finished optimizing without error and convergence is nil")
			return
		}
		if c.Convergence() != convergence.GradAbsTol.Convergence() {
			t.Errorf("For function " + fun.name + " convergence is not GradAbsTol. It is instead " + c.Convergence())
			return
		}
		firstObjVal := optVal
		if math.Abs(firstObjVal-fun.OptVal()) > MISO_TOLERANCE {
			t.Errorf("For function "+fun.name+" optimum value not found. %v found, %v expected", firstObjVal, fun.OptVal())
			return
		}
		firstLocVal := optLoc
		if !floats.Eq(firstLocVal, fun.OptLoc(), MISO_TOLERANCE) {
			t.Errorf("For function "+fun.name+" optimum location not found. %v found, %v expected", firstLocVal, fun.OptLoc())
		}
		firstNFunEvals := result.FunctionEvaluations
		firstNIterations := result.Iterations
		// Hack to reset FunEvals

		// Run it again to test that the reset works fine
		//opter.Loc().SetInit(fun.InitLoc)
		//_, _, c, err = opter.Optimize(fun, fun.InitLoc)
		optVal, optLoc, result, err = OptimizeGrad(fun, fun.InitLoc, settings, opter)
		if err != nil {
			t.Errorf("Error while re-using optimizer: ", err.Error())
		}

		if result.FunctionEvaluations != firstNFunEvals {
			t.Errorf("For function " + fun.name + "Different number of fun evals second time")
		}

		if result.Convergence.Convergence() != convergence.GradAbsTol.Convergence() {
			t.Errorf("For function " + fun.name + " convergence is not GradAbsTol second time")
		}

		if result.Iterations != firstNIterations {
			t.Errorf("For function " + fun.name + "Different number of fun evals second time")
		}

	}
}

func TestLbfgs(t *testing.T) {
	l := NewLbfgs()
	MisoGradBasedTest(t, l)
}
