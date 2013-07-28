package common

import (
	"testing"
	"time"
)

// Test to see if all of the pieces test for their convergence properly
func TestIterationsConvergence(t *testing.T) {
	// Test fun evals convergence
	iter := NewIterations()
	iter.Initialize()
	c := iter.Converged()
	if c != nil {
		t.Errorf("Iter converged after initialization")
	}
	iter.SetMax(10)
	iter.Add(5)
	c = iter.Converged()
	if c != nil {
		t.Errorf("Iter converged before max reached")
	}
	iter.Add(6)
	c = iter.Converged()
	if c == nil {
		t.Errorf("Iter did not converge when max reached")
	}
}

func TestFunctionEvaluationsConvergence(t *testing.T) {
	funevals := NewFunctionEvaluations()
	funevals.Initialize()
	c := funevals.Converged()
	if c != nil {
		t.Errorf("FunEvals converged after initialization")
	}
	funevals.SetMax(10)
	funevals.Add(5)
	c = funevals.Converged()
	if c != nil {
		t.Errorf("FunEvals converged before max reached")
	}
	funevals.Add(6)
	c = funevals.Converged()
	if c == nil {
		t.Errorf("FunEvals did not converge when max reached")
	}
}

func TestTimeConvergence(t *testing.T) {
	runtime := NewTime()
	runtime.Initialize()
	c := runtime.Converged()
	if c != nil {
		t.Errorf("Time converged after initialization")
	}
	runtime.SetMax(100 * time.Millisecond)
	time.Sleep(50 * time.Millisecond)
	c = runtime.Converged()
	if c != nil {
		t.Errorf("Time converged before max reached")
	}
	time.Sleep(60 * time.Millisecond)
	c = runtime.Converged()
	if c == nil {
		t.Errorf("Time did not converge when max reached")
	}
}

func TestCommonConvergence(t *testing.T) {
	common := NewOptCommon()
	common.Initialize()
	c := common.Converged()
	if c != nil {
		t.Errorf("Common converged after initialization")
	}
	// Test that common converges when iter converges
	common.Iter().SetMax(10)
	common.Iter().Add(5)
	c = common.Converged()
	if c != nil {
		t.Errorf("Common converged before it should")
	}
	common.Iter().Add(6)
	c = common.Converged()
	if c == nil {
		t.Errorf("Common did not converge when iter converged")
	}

	common = NewOptCommon()
	common.Initialize()
	common.FunEvals().SetMax(10)
	common.FunEvals().Add(11)
	c = common.Converged()
	if c == nil {
		t.Errorf("Common did not converge when FunEvals converged")
	}

	common = NewOptCommon()
	common.Initialize()
	common.Time().SetMax(100 * time.Millisecond)
	time.Sleep(120 * time.Millisecond)
	c = common.Converged()
	if c == nil {
		t.Errorf("Common did not converge when Time converged")
	}

}
