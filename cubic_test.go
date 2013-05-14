package gofunopter

import (
	"fmt"
	"math"
	"testing"
	"time"
)

var _ = fmt.Println

func Fun1(x float64) (ans, deriv float64) {

	c1 := 0.3
	c2 := 3.0
	ans = c1*math.Exp(-c2*(x-1)) + math.Exp((x - 1))
	deriv = -c1*c2*math.Exp(-c2*(x-1)) + math.Exp((x - 1))
	return ans, deriv
}

type SisoGBTest struct {
	f   float64
	g   float64
	Fun func(float64) (float64, float64)
}

func (t *SisoGBTest) Eval(val float64) (float64, float64, error) {
	f, g := t.Fun(val)
	return f, g, nil
}

func TestCubic(t *testing.T) {
	c := DefaultCubic()
	c.TimeInterval = 0 * time.Second
	c.HeadingInterval = 0
	problem := &SisoGBTest{Fun: Fun1}
	c.SetFun(problem)
	conv, err := Optimize(c)
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println(conv)
}
