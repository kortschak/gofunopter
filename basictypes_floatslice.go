package gofunopter

import (
	"fmt"
	"github.com/btracey/smatrix"
)

func DefaultLocationFloatSlice() *BasicOptFloatSlice {
	return NewBasicOptFloatSlice("Loc", false, nil)
}

// TODO: Turning display off makes it think it's converged. This is BAD
func DefaultGradientFloatSlice() *BasicTolFloatSlice {
	return NewBasicTolFloatSlice("Grad", false, nil, DefaultGradAbsTol, GradAbsTol, DefaultGradRelTol, GradRelTol)
}

// All the normal methods minus the tols
type BasicOptFloatSlice struct {
	save bool
	curr []float64
	init []float64
	hist [][]float64
	disp bool
	name string
	opt  []float64
}

func NewBasicOptFloatSlice(name string, disp bool, init []float64) *BasicOptFloatSlice {
	return &BasicOptFloatSlice{
		name: name,
		disp: disp,
		init: init,
	}
}

func (b *BasicOptFloatSlice) Hist() [][]float64 {
	return b.hist
}

func (b *BasicOptFloatSlice) Disp() bool {
	return b.disp
}

func (b *BasicOptFloatSlice) SetDisp(val bool) {
	b.disp = val
}

func (b *BasicOptFloatSlice) Save() bool {
	return b.save
}

func (b *BasicOptFloatSlice) SetSave(val bool) {
	b.save = val
}

func (b *BasicOptFloatSlice) AddToHist(val []float64) {
	// Make a copy so the pointer can change later
	newSlice := make([]float64, len(val))
	copy(newSlice, val)
	b.hist = append(b.hist, newSlice)
}

func (b *BasicOptFloatSlice) Curr() []float64 {
	return b.curr
}

func (b *BasicOptFloatSlice) SetCurr(val []float64) {
	b.curr = val
}

func (b *BasicOptFloatSlice) Init() []float64 {
	return b.init
}

func (b *BasicOptFloatSlice) SetInit(val []float64) {
	b.init = val
}

func (b *BasicOptFloatSlice) Initialize() error {
	b.hist = make([][]float64, 0)
	b.curr = b.init
	return nil
}

func (b *BasicOptFloatSlice) AppendHeadings(headings []string) []string {
	headings = append(headings, b.name)
	return headings
}

func (b *BasicOptFloatSlice) AppendValues(vals []interface{}) []interface{} {
	vals = append(vals, b.curr)
	return vals
}

func (b *BasicOptFloatSlice) SetResult() {
	b.opt = b.curr
}

func (b *BasicOptFloatSlice) Opt() []float64 {
	return b.opt
}

type BasicTolFloatSlice struct {
	*BasicOptFloatSlice
	absTol      float64
	absTolConv  Convergence
	relTol      float64
	relTolConv  Convergence
	avgNormCurr float64 // Two norm divided by the number of elements
	avgNormInit float64 // Two norm divided by the number of elements
}

func NewBasicTolFloatSlice(name string, disp bool, init []float64, absTol float64,
	absTolConv Convergence, relTol float64, relTolConv Convergence) *BasicTolFloatSlice {
	return &BasicTolFloatSlice{
		BasicOptFloatSlice: &BasicOptFloatSlice{name: name, disp: disp, init: init},
		absTol:             absTol,
		absTolConv:         absTolConv,
		relTol:             relTol,
		relTolConv:         relTolConv,
	}
}

func (b *BasicTolFloatSlice) SetInit(val []float64) {
	b.init = val
	b.avgNormInit = smatrix.VectorTwoNorm(val) / (float64(len(val)))
}

func (b *BasicTolFloatSlice) SetCurr(val []float64) {
	b.curr = val
	b.avgNormCurr = smatrix.VectorTwoNorm(val) / (float64(len(val)))
}

func (b *BasicTolFloatSlice) SetAbsTol(val float64) {
	b.absTol = val
}

func (b *BasicTolFloatSlice) AbsTol() float64 {
	return b.absTol
}
func (b *BasicTolFloatSlice) SetRelTol(val float64) {
	b.relTol = val
}

func (b *BasicTolFloatSlice) RelTol() float64 {
	return b.relTol
}

func (b *BasicTolFloatSlice) Initialize() error {
	err := b.BasicOptFloatSlice.Initialize()
	if err != nil {
		return err
	}
	b.avgNormInit = smatrix.VectorTwoNorm(b.init) / (float64(len(b.init)))
	b.avgNormCurr = b.avgNormInit
	return nil
}

func (b *BasicTolFloatSlice) Converged() Convergence {
	fmt.Println("In tol float slice converged")
	if b.avgNormCurr < b.absTol {
		fmt.Println("abs tol converged")
		return b.absTolConv
	}
	if b.avgNormCurr/b.avgNormInit < b.relTol {
		fmt.Println("rel tol converged")
		return b.relTolConv
	}
	fmt.Println("None converged")
	return nil
}