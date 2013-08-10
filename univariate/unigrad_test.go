package univariate

import (
	"github.com/btracey/gofunopter/common/optimize"
	"github.com/btracey/gofunopter/common/status"

	"fmt"
	"math"
	"testing"
)

var SISO_TOLERANCE float64 = 1E-6

type SisoGradTestFunction interface {
	optimize.UniObjGrad
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

func SisoGradBasedTest(t *testing.T, opter UniGradOptimizer) {
	funcs := OneDGradBasedFunctions()
	for _, fun := range funcs {
		// Run it once until very converged
		//opter.Loc().SetInit(fun.InitLoc)

		settings := NewUniGradSettings()
		settings.GradientAbsoluteTolerance = 1e-14
		settings.MaximumIterations = 50

		fmt.Println("Init loc", fun.InitLoc)
		fmt.Println("Settings", settings)
		fmt.Printf("opter %#v \n", opter)

		optVal, optLoc, result, err := OptimizeGrad(fun, fun.InitLoc, settings, opter)

		//optVal, optLoc, c, err := opter.Optimize(fun, fun.InitLoc)
		if err != nil {
			t.Errorf("Error during optimization for function " + fun.Name() + ": " + err.Error())
			continue
		}

		c := result.Status

		if c == status.Continue {
			fmt.Println(c)
			t.Errorf("Finished optimizing without error and status is continue")
			continue
		}
		if c != status.GradAbsTol {
			t.Errorf("For function " + fun.Name() + " status is not GradAbsTol")
			continue
		}
		if math.Abs(optVal-fun.OptVal()) > SISO_TOLERANCE {
			t.Errorf("For function "+fun.Name()+" optimum value not found. %v found, %v expected", optVal, fun.OptVal())
			continue
		}
		if math.Abs(optLoc-fun.OptLoc()) > SISO_TOLERANCE {
			t.Errorf("For function "+fun.Name()+" optimum location not found. %v found, %v expected", optLoc, fun.OptLoc())
		}
		//firstNFunEvals := opter.GetOptCommon().FunEvals().Opt()
		//firstNIterations := opter.GetOptCommon().Iter().Opt()
		firstNFunEvals := result.FunctionEvaluations
		firstNIterations := result.Iterations

		// Hack to reset FunEvals

		fmt.Println("Init loc", fun.InitLoc)
		fmt.Println("Settings", settings)
		fmt.Printf("opter %#v \n", opter)

		// Run it again to test that the reset works fine
		//opter.Loc().SetInit(fun.InitLoc)
		_, _, result, err = OptimizeGrad(fun, fun.InitLoc, settings, opter)
		if err != nil {
			t.Errorf("Error while re-using optimizer")
		}

		if result.FunctionEvaluations != firstNFunEvals {
			t.Errorf("For function " + fun.Name() + "Different number of fun evals second time")
		}

		if result.Status != status.GradAbsTol {
			t.Errorf("For function " + fun.Name() + " status is not GradAbsTol second time")
		}

		if result.Iterations != firstNIterations {
			t.Errorf("For function " + fun.Name() + "Different number of fun evals second time")
		}

	}
}

var _ = fmt.Println

type SumExpStruct struct{}

func (s SumExpStruct) ObjGrad(x float64) (f, g float64, err error) {

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
