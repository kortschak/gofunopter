package gofunopter

import (
	//"fmt"
	"testing"
	"time"
)

type FakeDisplayer struct {
	*Display
}

func (f *FakeDisplayer) DisplayHeadings() []string {
	return []string{"Iter", "FunEvals", "Grad"}
}

func (f *FakeDisplayer) DisplayValues() []interface{} {
	return []interface{}{10, 15, 2.0}
}

func TestDisplay(t *testing.T) {
	f := &FakeDisplayer{}
	f.Display = DefaultDisplay()
	SetDisplayMethods(f)
	f.Display.HeadingInterval = 3
	f.Display.TimeInterval = time.Millisecond
	for i := 0; i < 10; i++ {
		time.Sleep(300 * time.Microsecond)
		f.Iterate()
	}

	//TODO: Somehow make a test that's non-visual
}
