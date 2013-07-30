package gofunopter

import (
	"github.com/btracey/gofunopter/convergence"
	"github.com/btracey/gofunopter/optimize"

	"fmt"
	"math"
	"testing"
)

var SISO_TOLERANCE float64 = 1E-6

type SisoGradTestFunction interface {
	optimize.SisoGrad
	Name() string
	OptVal() float64
	OptLoc() float64
}

type SisoGradTest struct {
	SisoGradTestFunction
	InitLoc float64
}

func OneDGradBasedFunctions() []SisoGradTest {
	return []SisoGradTest{SisoGradTest{SumExpStruct{}, 2}}
}

func SisoGradBasedTest(t *testing.T, opter optimize.SisoGradOptimizer) {
	funcs := OneDGradBasedFunctions()
	for _, fun := range funcs {
		// Run it once until very converged
		opter.Loc().SetInit(fun.InitLoc)
		opter.Grad().SetAbsTol(1E-14)
		opter.GetOptCommon().FunEvals().SetMax(50)
		opter.GetDisplay().SetValueDisplayInterval(0)
		c, err := opter.Optimize(fun)
		if err != nil {
			t.Errorf("Error during optimization for function " + fun.Name() + ": " + err.Error())
			continue
		}
		if c == nil {
			t.Errorf("Finished optimizing without error and convergence is nil")
			continue
		}
		if c.Convergence() != convergence.GradAbsTol.Convergence() {
			t.Errorf("For function " + fun.Name() + " convergence is not GradAbsTol. It is instead " + c.Convergence())
			continue
		}
		firstObjVal := opter.Obj().Opt()
		if math.Abs(firstObjVal-fun.OptVal()) > SISO_TOLERANCE {
			t.Errorf("For function "+fun.Name()+" optimum value not found. %v found, %v expected", firstObjVal, fun.OptVal())
			continue
		}
		firstLocVal := opter.Loc().Opt()
		if math.Abs(firstLocVal-fun.OptLoc()) > SISO_TOLERANCE {
			t.Errorf("For function "+fun.Name()+" optimum location not found. %v found, %v expected", firstLocVal, fun.OptLoc())
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
			t.Errorf("For function " + fun.Name() + "Different number of fun evals second time")
		}

		if c.Convergence() != convergence.GradAbsTol.Convergence() {
			t.Errorf("For function " + fun.Name() + " convergence is not GradAbsTol second time")
		}

		if opter.GetOptCommon().Iter().Opt() != firstNIterations {
			t.Errorf("For function " + fun.Name() + "Different number of fun evals second time")
		}

	}
}

var _ = fmt.Println

type SumExpStruct struct{}

func (s SumExpStruct) Eval(x float64) (f, g float64, err error) {

	// http://www.wolframalpha.com/input/?i=0.3+*+exp%28+-+3+%28x-1%29%29+%2B+exp%28x-1%29
	c1 := 0.3
	c2 := 3.0
	f = c1*math.Exp(-c2*(x-1)) + math.Exp((x - 1))
	g = -c1*c2*math.Exp(-c2*(x-1)) + math.Exp((x - 1))
	return f, g, nil
}

func (s SumExpStruct) Name() string {
	return "SumExp"
}

func (s SumExpStruct) OptVal() float64 {
	return 1.298671661900395685896941595
}

func (s SumExpStruct) OptLoc() float64 {
	return 0.9736598710855434246931247548
}

func TestCubic(t *testing.T) {
	c := NewCubic()
	SisoGradBasedTest(t, c)
}
