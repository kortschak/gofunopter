package gofunopter

import (
	"fmt"
	"github.com/btracey/gofunopter/common"
	"github.com/btracey/gofunopter/convergence"
	"github.com/btracey/gofunopter/display"
	"github.com/btracey/gofunopter/optimize"
	"testing"
)

type Template struct {
	*common.OptCommon
	*display.Display
	disp bool
}

func NewTemplate() *Template {
	return &Template{
		OptCommon: common.NewOptCommon(),
		Display:   display.NewDisplay(),
		disp:      true,
	}
}

func (t *Template) Disp() bool {
	return t.disp
}

func (t *Template) SetDisp(b bool) {
	t.disp = b
}

func (t *Template) Initialize() (err error) {
	return nil
}

func (t *Template) Iterate() (nFunEvals int, err error) {
	return 2, nil
}

func (t *Template) Converged() convergence.C {
	return nil
}

func (t *Template) SetResult() {
}

type FakeFunction struct{}

func (f *FakeFunction) Eval(x float64) float64 {
	return 0
}

// Tests that we can correctly stop the function due to
// numFunEvals
func TestFunEvalsConvergence(t *testing.T) {
	fun := &FakeFunction{}

	opt := NewTemplate()
	opt.FunEvals().SetMax(10)

	fmt.Println("Starting call to optimize")
	c, err := optimize.OptimizeOpter(opt, fun)
	if err != nil {
		t.Errorf("Error during optimization: " + err.Error())
	}
	if c == nil {
		t.Errorf("No convergence reached")
	}
	if c.Convergence() != convergence.FunEvals.Convergence() {
		t.Errorf("Convergence not from function evaluations")
	}
}

// Tests that we can correctly stop the function due to
// numFunEvals
func TestIterationsConvergence(t *testing.T) {
	fun := &FakeFunction{}

	opt := NewTemplate()
	opt.Iter().SetMax(10)
	c, err := optimize.OptimizeOpter(opt, fun)
	if err != nil {
		t.Errorf("Error during optimization: " + err.Error())
	}
	if c == nil {
		t.Errorf("No convergence reached")
	}
	if c.Convergence() != convergence.Iterations.Convergence() {
		t.Errorf("Convergence not from iterations evaluations")
	}
}
