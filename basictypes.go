package gofunopter

import (
	"fmt"
	//"github.com/btracey/smatrix"
	"math"
)

//TODO: Be more careful with resetting and error checking during optimization

var _ = fmt.Println

type BasicConvergence struct {
	Str string
}

func (b BasicConvergence) ConvergenceType() string {
	return b.Str
}

type GradConvergence struct{ Str string }

func (b GradConvergence) ConvergenceType() string {
	return b.Str
}

type LocConvergence struct{ Str string }

func (b LocConvergence) ConvergenceType() string {
	return b.Str
}

type FunConvergence struct{ Str string }

func (b FunConvergence) ConvergenceType() string {
	return b.Str
}

type StepConvergence struct{ Str string }

func (b StepConvergence) ConvergenceType() string {
	return b.Str
}

// Counts up and converges if there is a maximum
type Counter struct {
	max   int // Maximum allowable value of the counter
	curr  int // current value of the counter
	total int // Total number at the end of the optimization run
	conv  Convergence
	name  string
}

func NewCounter(name string, max int, conv Convergence) *Counter {
	return &Counter{max: max, conv: conv}
}

func (c *Counter) Max() int {
	return c.max
}

func (c *Counter) SetMax(val int) {
	c.max = val
}

func (c *Counter) Add(delta int) {
	c.curr += delta
}

//var MaxCounter = &BasicConvergence{"Max counter reached"}

func (c *Counter) Converged() Convergence {
	// returns a bool because we want to force implementers to make
	// a better convergence type for whatever they're using the counter for
	if c.curr > c.max {
		return c.conv
	}
	return nil
}

func (c *Counter) SetOpt() {
	c.total = c.curr
}

func (c *Counter) Curr() int {
	return c.curr
}

func (c *Counter) Opt() int {
	return c.total
}

func (c *Counter) AppendHeadings(strs []string) []string {
	return append(strs, c.name)
}

func (c *Counter) AppendValues(vals []interface{}) []interface{} {
	return append(vals, c.curr)
}

// All the normal methods minus the tols
type BasicOptFloat struct {
	save bool
	curr float64
	init float64
	hist []float64
	disp bool
	name string
	opt  float64
}

func NewBasicOptFloat(name string, disp bool, init float64) *BasicOptFloat {
	return &BasicOptFloat{
		name: name,
		disp: disp,
		init: init,
	}
}

func (b *BasicOptFloat) Disp() bool {
	return b.disp
}

func (b *BasicOptFloat) SetDisp(val bool) {
	b.disp = val
}

func (b *BasicOptFloat) Save() float64 {
	return b.save
}

func (b *BasicOptFloat) SetSave(val bool) {
	b.save = val
}

func (b *BasicOptFloat) AddToHist(val float64) {
	b.hist = append(b.hist, val)
}

func (b *BasicOptFloat) Curr() float64 {
	return b.curr
}

func (b *BasicOptFloat) SetCurr(val float64) float64 {
	b.curr = val
}

func (b *BasicOptFloat) Init() float64 {
	return b.init
}

func (b *BasicOptFloat) SetInit(val float64) {
	b.init = val
}

func (b *BasicOptFloat) Initialize() error {
	b.save = make([]float64, 0)
	b.curr = b.init
}
func (b *BasicOptFloat) Converged() Convergence {
	return nil
}

func (b *BasicOptFloat) AppendHeadings(headings []string) []string {
	headings = append(headings, b.name)
}

func (b *BasicOptFloat) AppendValues(vals []interface{}) []interface{} {
	vals = append(vals, b.curr)
}

func (b *BasicOptFloat) SetResult() {
	b.opt = b.curr
}

func (b *BasicOptFloat) Opt() float64 {
	b.opt = b.curr
}

type BasicTolerFloat struct {
	*BasicOptFloat
	absTol     float64
	absTolConv Convergence
	relTol     float64
	relTolConv Convergence
	absCurr    float64
	absInit    float64
}

func NewBasicTolerFloat(name string, disp bool, init float64, absTol float64,
	absTolConv Convergence, relTol float64, relTolConv Convergence) *BasicTolerFloat {
	return &BasicTolerFloat{
		BasicOptFloat: &BasicOptFloat{name: name, disp: disp, init: init},
		absTol:        absTol,
		absTolCov:     absTolConv,
		relTol:        relTol,
		relTolConv:    relTolConv,
	}
}

func (b *BasicTolerFloat) SetInit(val float64) {
	b.init = val
	b.absInit = math.Abs(val)
}

func (b *BasicTolerFloat) SetCurr(val float64) {
	b.curr = val
	b.absCurr = math.Abs(val)
}

func (b *BasicTolerFloat) SetAbsTol(val float64) {
	b.absTol = val
}

func (b *BasicTolerFloat) AbsTol() float64 {
	return b.absTol
}
func (b *BasicTolerFloat) SetRelTol(val float64) {
	b.relTol = val
}

func (b *BasicTolerFloat) RelTol() float64 {
	return b.relTol
}

func (b *BasicTolerFloat) Converged() Convergence {
	if b.absCurr < b.absTol {
		return b.absTolConv
	}
	if b.absCurr/b.absInit < b.absTol {
		return b.relTolConv
	}
	return nil
}

func NewAbsTol(tol float64, c Convergence) *AbsTolStruct {
	return &AbsTolStruct{
		tol:         tol,
		Convergence: c,
	}
}

var LocAbsTol Convergence = LocConvergence{"Location absolute tolerance reached"}
var LocRelTol Convergence = LocConvergence{"Location relative tolerance reached"}
var ObjAbsTol Convergence = FunConvergence{"Function absolute tolerance reached"}
var ObjRelTol Convergence = FunConvergence{"Function relative tolerance reached"}
var GradAbsTol Convergence = GradConvergence{"Gradient absolute tolerance reached"}
var GradRelTol Convergence = GradConvergence{"Gradient relative tolerance reached"}

const DefaultGradAbsTol = 1E-6
const DefaultGradRelTol = 1E-8

func DefaultLocationFloat() *BasicOptFloat {
	return NewBasicOptFloat("Loc", false, 0)
}

func DefaultObjectiveFloat() *BasicTolerFloat {
	return NewBasicTolerFloat("Obj", true, math.NaN(), 0, ObjAbsTol, 0, ObjRelTol)
}

func DefaultGradientFloat() *BasicTolerFloat {
	return NewBasicTolerFloat("Grad", true, math.NaN(), DefaultGradAbsTol, GradAbsTol, DefaultGradRelTol, GradRelTol)
}

// TODO: Implement display bounds
type BoundsFloatStruct struct {
	lb float64
	ub float64
	CurrInitGetSetterFloat
	abstol AbsTol
	reltol RelTol
	Name   string
	Displayer
}

func NewBoundsFloat(name string, lb, ub, abstol float64, absconv Convergence, reltol float64, relconv Convergence) *BoundsFloatStruct {
	b := &BoundsFloatStruct{}
	b.lb = math.Inf(-1)
	b.ub = math.Inf(1)
	b.CurrInitGetSetterFloat, b.abstol, b.reltol = NewAbsRelTolStruct(b.ub-b.lb, abstol, absconv, reltol, relconv)
	b.Name = name
	b.Displayer = NewDisplay(false)
	return b
}

func (b *BoundsFloatStruct) AbsTol() AbsTol {
	return b.abstol
}

func (b *BoundsFloatStruct) RelTol() RelTol {
	return b.reltol
}

func (s *BoundsFloatStruct) Lb() float64 {
	return s.lb
}

func (s *BoundsFloatStruct) Ub() float64 {
	return s.ub
}

func (s *BoundsFloatStruct) SetLb(val float64) {
	s.lb = val
	s.SetCurr(s.ub - s.lb)
}

func (s *BoundsFloatStruct) SetUb(val float64) {
	s.ub = val
	s.SetCurr(s.ub - s.lb)
}

func (s *BoundsFloatStruct) AppendHeadings(strs []string) []string {
	return append(strs, s.Name+"LB", s.Name+"UB")
}

func (s *BoundsFloatStruct) AppendValues(vals []interface{}) []interface{} {
	return append(vals, s.lb, s.lb)
}

// Midpoint between the bounds
func (s *BoundsFloatStruct) Midpoint() float64 {
	return (s.lb + s.ub) / 2.0
}

// Is the value between the upper and lower bounds
func (s *BoundsFloatStruct) WithinBounds(val float64) bool {
	if val < s.lb {
		return false
	}
	if val > s.ub {
		return false
	}
	return true
}

//var BoundgapAbsTol Convergence = BasicConvergence{"Bound gap absolute tolerance reached"}

func (s *BoundsFloatStruct) Converged() Convergence {
	return Converged(s.AbsTol(), s.RelTol())
}

var StepAbsTol Convergence = StepConvergence{"Step absolute tolerance reached"}
var StepRelTol Convergence = StepConvergence{"Step relative tolerance reached"}
var StepBoundsAbsTol Convergence = StepConvergence{"Step bounds absolute tolerance reached"}
var StepBoundsRelTol Convergence = StepConvergence{"Step bounds absolute tolerance reached"}

type BoundedStepFloatStruct struct {
	CurrInitGetSetterFloat
	HistorySaverFloat
	Displayer
	*BoundsFloatStruct
}

func (s *BoundedStepFloatStruct) Initialize() error {
	return nil
}

const DefaultInitStepSize = 1
const DefaultBoundedStepFloatAbsTol = 0 //
const DefaultBoundedStepFloatRelTol = 0 // Turn off step rel tol
// Returns the default values for a step size
// no default relative tolerance
func DefaultBoundedStepFloat() *BoundedStepFloatStruct {
	return &BoundedStepFloatStruct{
		CurrInitGetSetterFloat: NewCurrInitFloat(DefaultInitStepSize),
		HistorySaverFloat:      DefaultHistorySaverFloat(),
		Displayer:              NewDisplay(false),
		BoundsFloatStruct:      NewBoundsFloat("Step", math.Inf(-1), math.Inf(1), DefaultBoundedStepFloatAbsTol, StepBoundsAbsTol, DefaultBoundedStepFloatRelTol, StepBoundsRelTol),
	}
}

//func (b *BoundedStepFloatStruct) Converged() Convergence {
//	return Converged(c.BoundsFloatStruct)
//}
