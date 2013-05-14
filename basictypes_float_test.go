package gofunopter

import (
	//"fmt"
	"testing"
)

func BasicOptEq(b1, b2 *BasicOptFloat) string {
	// Checks if they are equal
	str := ""
	if b1.save != b2.save {
		str += " Save doesn't match "
	}
	if b1.curr != b2.curr {
		str += " Curr doesn't match "
	}
	if b1.init != b2.init {
		str += " Init doesn't match "
	}
	if b1.disp != b2.disp {
		str += " Disp doesn't match "
	}
	if b1.name != b2.name {
		str += " Name doesn't match "
	}
	if b1.opt != b2.opt {
		str += " Opt doesn't match "
	}
	if len(b1.hist) != len(b2.hist) {
		str += " Hist length doesn't match"
	} else {
		for i, val := range b2.hist {
			if b1.hist[i] != val {
				str += " Hist doesn't match "
			}
		}
	}
	return str
}

// Test script to make sure all the methods are working the way we think they are

func TestBasicOptFloat(t *testing.T) {
	save := true
	curr := 1.0
	init := 0.5
	disp := true
	name := "Test"
	opt := 0.0
	hist := make([]float64, 0)

	b := &BasicOptFloat{
		save: save,
		curr: curr,
		init: init,
		hist: hist,
		disp: disp,
		name: name,
		opt:  opt,
	}

	b2 := &BasicOptFloat{
		save: save,
		curr: curr,
		init: init,
		hist: hist,
		disp: disp,
		name: name,
		opt:  opt,
	}

	str := BasicOptEq(b, b2)
	if str != "" {
		t.Errorf("Something wrong with BasicOptEq code")
	}
	newCurr := 5.0
	b.SetCurr(newCurr)
	b2.curr = newCurr
	str = BasicOptEq(b, b2)
	if str != "" {
		t.Errorf("Curr setter doesn't behave as expected: " + str)
	}
	if b.Curr() != newCurr {
		t.Errorf("Curr getter and setter don't match")
	}
	str = BasicOptEq(b, b2)
	if str != "" {
		t.Errorf("Curr getter doesn't behave as expected: " + str)
	}
	newInit := 6.0
	b.SetInit(newInit)
	b2.init = newInit
	str = BasicOptEq(b, b2)
	if str != "" {
		t.Errorf("Init setter doesn't behave as expected: " + str)
	}
	if b.Init() != newInit {
		t.Errorf("Init getter and setter don't match")
	}
	str = BasicOptEq(b, b2)
	if str != "" {
		t.Errorf("Init getter doesn't behave as expected: " + str)
	}
	b.SetDisp(!disp)
	b2.disp = !disp
	str = BasicOptEq(b, b2)
	if str != "" {
		t.Errorf("Disp setter doesn't behave as expected: " + str)
	}
	if b.Disp() != !disp {
		t.Errorf("Disp getter and setter don't match")
	}
	str = BasicOptEq(b, b2)
	if str != "" {
		t.Errorf("Disp getter doesn't behave as expected: " + str)
	}
	b.AddToHist(5)
	b.AddToHist(6)
	b.AddToHist(7)
	b2.hist = []float64{5, 6, 7}
	str = BasicOptEq(b, b2)
	if str != "" {
		t.Errorf("Add to hist doesn't behave as expected: " + str)
	}
	b.Initialize()
	b2.curr = b2.init
	b2.hist = make([]float64, 0)
	str = BasicOptEq(b, b2)
	if str != "" {
		t.Errorf("Initialize doesn't behave as expected: " + str)
	}

	b.SetResult()
	b2.opt = b2.curr
	str = BasicOptEq(b, b2)
	if str != "" {
		t.Errorf("Initialize doesn't behave as expected: " + str)
	}
	if b.Opt() != b.opt {
		t.Errorf("Opt not set properly after result")
	}
	str = BasicOptEq(b, b2)
	if str != "" {
		t.Errorf("Opt getter doesn't behave as expected: " + str)
	}
}
