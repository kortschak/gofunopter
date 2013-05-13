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

func (t *SisoGBTest) Eval(val float64) error {
	t.f, t.g = t.Fun(val)
	return nil
}

func (t *SisoGBTest) Obj() float64 {
	return t.f
}

func (t *SisoGBTest) Grad() float64 {
	return t.g
}

func TestCubic(t *testing.T) {
	c := DefaultCubic()
	c.TimeInterval = 0 * time.Second
	c.HeadingInterval = 0
	fmt.Println("In test", c.Obj.Init)
	problem := &SisoGBTest{Fun: Fun1}
	c.Fun = problem
	conv, err := Optimize(c)
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println(conv)
}
