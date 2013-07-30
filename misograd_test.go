package gofunopter

import (
	"gofunopter/convergence"
	"gofunopter/optimize"

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
	optimize.MisoGrad
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

func MisoGradBasedTest(t *testing.T, opter optimize.MisoGradOptimizer) {
	funcs := MisoGradFunctions()
	for _, fun := range funcs {
		// Run it once until very converged
		opter.Loc().SetInit(fun.InitLoc)
		opter.Grad().SetAbsTol(1E-14)
		opter.GetOptCommon().FunEvals().SetMax(1000)
		opter.GetDisplay().SetValueDisplayInterval(0)
		fmt.Println("Is misograd_test, starting optimizer")
		c, err := opter.Optimize(fun)
		if err != nil {
			t.Errorf("Error during optimization for function " + fun.name + ": " + err.Error())
			continue
		}
		if c == nil {
			t.Errorf("Finished optimizing without error and convergence is nil")
			continue
		}
		if c.Convergence() != convergence.GradAbsTol.Convergence() {
			t.Errorf("For function " + fun.name + " convergence is not GradAbsTol. It is instead " + c.Convergence())
			continue
		}
		firstObjVal := opter.Obj().Opt()
		if math.Abs(firstObjVal-fun.OptVal()) > MISO_TOLERANCE {
			t.Errorf("For function "+fun.name+" optimum value not found. %v found, %v expected", firstObjVal, fun.OptVal())
			continue
		}
		firstLocVal := opter.Loc().Opt()
		if !floats.Eq(firstLocVal, fun.OptLoc(), MISO_TOLERANCE) {
			t.Errorf("For function "+fun.name+" optimum location not found. %v found, %v expected", firstLocVal, fun.OptLoc())
		}
		firstNFunEvals := opter.GetOptCommon().FunEvals().Opt()
		firstNIterations := opter.GetOptCommon().Iter().Opt()
		// Hack to reset FunEvals

		// Run it again to test that the reset works fine
		opter.Loc().SetInit(fun.InitLoc)
		c, err = opter.Optimize(fun)
		if err != nil {
			t.Errorf("Error while re-using optimizer")
		}

		if opter.GetOptCommon().FunEvals().Opt() != firstNFunEvals {
			t.Errorf("For function " + fun.name + "Different number of fun evals second time")
		}

		if c.Convergence() != convergence.GradAbsTol.Convergence() {
			t.Errorf("For function " + fun.name + " convergence is not GradAbsTol second time")
		}

		if opter.GetOptCommon().Iter().Opt() != firstNIterations {
			t.Errorf("For function " + fun.name + "Different number of fun evals second time")
		}

	}
}

func TestLbfgs(t *testing.T) {
	l := NewLbfgs()
	MisoGradBasedTest(t, l)
}