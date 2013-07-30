package univariate

import (
	"testing"
)

func TestGradientConvergence(t *testing.T) {
	// First, test that it can converge at all
	g := NewGradient()
	g.curr = 1E-3
	g.absCurr = 1E-3
	g.SetAbsTol(1E-2)
	g.SetRelTol(0)
	g.init = 1E-3

	c := g.Converged()
	if c == nil {
		t.Errorf("Gradient incorrectly reported no convergence")
	}
	g.curr = 1E-1
	g.absCurr = 1E-1
	c = g.Converged()
	if c != nil {
		t.Errorf("Gradient incorrectly reported convergence")
	}
	g.SetCurr(1E-5)
	c = g.Converged()
	if c == nil {
		t.Errorf("After SetCurr, gradient reported no convergence ")
	}
	g.SetAbsTol(1E-8)
	c = g.Converged()
	if c != nil {
		t.Errorf("After SetAbsTol, gradient reported convergence: " + c.Convergence())
	}
	g.SetCurr(-1E-4)
	c = g.Converged()
	if c != nil {
		t.Errorf("After setting a large negative value of current, reported convergence")
	}
	g.SetCurr(-1E-9)
	c = g.Converged()
	if c == nil {
		t.Errorf("After setting a small negative value of current, reported no convergence")
	}
	g.SetCurr(1E-2)
}
