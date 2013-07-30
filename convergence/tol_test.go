package convergence

import "testing"

func TestAbsSet(t *testing.T) {
	a := NewAbs(1E-6, Basic{})
	newTol := 1E-7
	a.SetAbsTol(newTol)
	if a.tol != newTol {
		t.Errorf("Setting abs tol failed")
	}
}

func TestAbsTolGet(t *testing.T) {
	tol := 1E-6
	a := NewAbs(tol, Basic{})
	get := a.AbsTol()
	if get != tol {
		t.Errorf("Getting abs tol failed")
	}
}

func TestAbsTolCheckConv(t *testing.T) {
	tol := 1E-6
	a := NewAbs(tol, Basic{})
	c := a.CheckConvergence(1)
	if c != nil {
		t.Errorf("Passed convergence when curr > tol")
	}
	c = a.CheckConvergence(1E-8)
	if c == nil {
		t.Errorf("Failed convergence when curr < tol")
	}
}

func TestRelTolSet(t *testing.T) {
	a := NewRel(1E-6, Basic{})
	newTol := 1E-7
	a.SetRelTol(newTol)
	if a.tol != newTol {
		t.Errorf("Setting abs tol failed")
	}
}

func TestRelTolGet(t *testing.T) {
	tol := 1E-6
	a := NewRel(tol, Basic{})
	get := a.RelTol()
	if get != tol {
		t.Errorf("Getting abs tol failed")
	}
}

func TestRelTolCheckConv(t *testing.T) {
	tol := 1E-6
	a := NewRel(tol, Basic{})
	c := a.CheckConvergence(1, 3)
	if c != nil {
		t.Errorf("Passed convergence when curr > tol*init")
	}
	c = a.CheckConvergence(1E-6, 3)
	if c == nil {
		t.Errorf("Failed convergence when curr < tol*init")
	}
}
