package gofunopter

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"
)

var _ = fmt.Println

func Rosenbrock(x []float64) (sum float64, deriv []float64) {
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
	return sum, deriv
}

type MisoGBTest struct {
	Fun func([]float64) (float64, []float64)
}

func (t *MisoGBTest) Eval(val []float64) (float64, []float64, error) {
	f, g := t.Fun(val)
	return f, g, nil
}

func TestLbfgs(t *testing.T) {
	nDim := 10
	scale := 10.0
	c := DefaultLbfgs()
	//fmt.Println("Init loc", c.Loc().Init())
	c.TimeInterval = 0 * time.Second
	c.HeadingInterval = 0
	problem := &MisoGBTest{Fun: Rosenbrock}
	c.SetFun(problem)
	initX := make([]float64, nDim)
	for i := range initX {
		initX[i] = rand.Float64() * scale
	}
	c.Loc().SetInit(initX)
	conv, err := Optimize(c)
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println("Num fun evals", c.FunEvals().Opt())
	fmt.Println(conv)
}
