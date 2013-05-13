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

func (c *Counter) Total() int {
	return c.total
}

/*
func (c *Counter) SetTotal(val int) {
	c.total = val
}
*/

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

func (c *Counter) Result() {
	c.total = c.curr
}

func (c *Counter) Curr() int {
	return c.curr
}

func (c *Counter) AppendHeadings(strs []string) []string {
	return append(strs, c.name)
}

func (c *Counter) AppendValues(vals []interface{}) []interface{} {
	return append(vals, c.curr)
}

type BasicHistory struct {
	save bool
}

func (b *BasicHistory) Save() bool {
	return b.save
}

func (b *BasicHistory) SetSave(val bool) {
	b.save = val
}

func DefaultBasicHistory() *BasicHistory {
	return &BasicHistory{save: false}
	//return b
}

// Something about only major iterations?
type BasicHistoryFloat struct {
	hist []float64
	*BasicHistory
}

func (h *BasicHistoryFloat) Set(val []float64) {
	h.hist = val
}

func (h *BasicHistoryFloat) Get() []float64 {
	return h.hist
}

func (h *BasicHistoryFloat) Add(val float64) {
	if h.Save() {
		h.hist = append(h.hist, val)
	}
}

func DefaultHistoryFloat() *BasicHistoryFloat {
	return &BasicHistoryFloat{hist: make([]float64, 0), BasicHistory: DefaultBasicHistory()}
}

type HistorySaverFloatStruct struct {
	hist HistoryFloat
}

func (h *HistorySaverFloatStruct) Hist() HistoryFloat {
	return h.hist
}

func DefaultHistorySaverFloat() *HistorySaverFloatStruct {
	return &HistorySaverFloatStruct{hist: DefaultHistoryFloat()}
}

// Float which implements the curr set getter interface
type CurrFloatStruct struct {
	curr float64
}

func NewCurrFloat() *CurrFloatStruct {
	return &CurrFloatStruct{}
}

func (b *CurrFloatStruct) Curr() float64 {
	return b.curr
}

func (b *CurrFloatStruct) SetCurr(val float64) {
	b.curr = val
}

type InitFloatStruct struct {
	init float64
}

func NewInitFloat(val float64) *InitFloatStruct {
	return &InitFloatStruct{init: val}
}

func (i *InitFloatStruct) Init() float64 {
	return i.init
}

func (i *InitFloatStruct) SetInit(val float64) {
	i.init = val
}

type CurrInitFloatStruct struct {
	InitGetSetterFloat
	CurrGetSetterFloat
}

func NewCurrInitFloat(init float64) *CurrInitFloatStruct {
	return &CurrInitFloatStruct{
		InitGetSetterFloat: NewInitFloat(init),
		CurrGetSetterFloat: NewCurrFloat(),
	}
}

type AbsTolStruct struct {
	tol float64
	CurrGetSetterFloat
	Convergence
}

func (b *AbsTolStruct) Tol() float64 {
	return b.tol
}

func (b *AbsTolStruct) SetTol(val float64) {
	b.tol = val
}

func (b *AbsTolStruct) Converged() Convergence {
	if b.Curr() < b.tol {
		return b.Convergence
	}
	return nil
}

func NewAbsTol(tol float64, c Convergence) *AbsTolStruct {
	return &AbsTolStruct{
		tol:         tol,
		Convergence: c,
	}
}

type RelTolStruct struct {
	tol float64
	CurrInitGetSetterFloat
	Convergence
}

func NewRelTol(tol float64, c Convergence) *RelTolStruct {
	return &RelTolStruct{
		tol:         tol,
		Convergence: c,
	}
}

func (b *RelTolStruct) Tol() float64 {
	return b.tol
}

func (b *RelTolStruct) SetTol(val float64) {
	b.tol = val
}

func (b *RelTolStruct) Converged() Convergence {
	if b.Curr()/b.Init() < b.tol {
		return b.Convergence
	}
	return nil
}

var LocAbsTol Convergence = LocConvergence{"Location absolute tolerance reached"}
var LocRelTol Convergence = LocConvergence{"Location relative tolerance reached"}
var ObjAbsTol Convergence = FunConvergence{"Function absolute tolerance reached"}
var ObjRelTol Convergence = FunConvergence{"Function relative tolerance reached"}
var GradAbsTol Convergence = GradConvergence{"Gradient absolute tolerance reached"}
var GradRelTol Convergence = GradConvergence{"Gradient relative tolerance reached"}

// Returns the default values for an input location
// Locations don't have any tolerances
type LocationFloatStruct struct {
	CurrInitGetSetterFloat
	HistorySaverFloat
	Displayer
}

/*
func (c *LocationFloatStruct) Initialize() error {
	c.SetCurr(c.Init())
	return nil
}
*/

func (b *LocationFloatStruct) Converged() Convergence {
	return nil
}

func (c *LocationFloatStruct) Initialize() error {
	c.SetCurr(c.Init())
	return nil
}

func (l *LocationFloatStruct) AppendHeadings(vals []string) []string {
	return append(vals, "Loc")
}

func (l *LocationFloatStruct) AppendValues(vals []interface{}) []interface{} {
	return append(vals, l.Curr())
}

// Gets AppendValues from CurrFloat

func DefaultLocationFloat() *LocationFloatStruct {
	return &LocationFloatStruct{
		Displayer:         NewDisplay(false),
		HistorySaverFloat: DefaultHistorySaverFloat(),
	}
}

func NewAbsRelTolStruct(init float64, abstol float64, absConv Convergence, reltol float64, relConv Convergence) (CurrInitGetSetterFloat, AbsTol, RelTol) {
	absTol := NewAbsTol(abstol, absConv)
	relTol := NewRelTol(reltol, relConv)
	currInit := NewCurrInitFloat(init)
	absTol.CurrGetSetterFloat = currInit
	relTol.CurrInitGetSetterFloat = currInit
	return currInit, absTol, relTol
}

type ObjectiveFloatStruct struct {
	CurrInitGetSetterFloat
	abstol AbsTol
	reltol RelTol
	HistorySaverFloat
	Displayer
}

func (c *ObjectiveFloatStruct) Initialize() error {
	c.SetCurr(c.Init())
	return nil
}

func DefaultObjectiveFloat() *ObjectiveFloatStruct {
	o := &ObjectiveFloatStruct{}
	o.CurrInitGetSetterFloat, o.abstol, o.reltol = NewAbsRelTolStruct(math.NaN(), 0, ObjAbsTol, 0, ObjRelTol)
	o.HistorySaverFloat = DefaultHistorySaverFloat()
	return o
}

func (b *ObjectiveFloatStruct) AbsTol() AbsTol {
	return b.abstol
}

func (b *ObjectiveFloatStruct) RelTol() RelTol {
	return b.reltol
}

func (l *ObjectiveFloatStruct) AppendHeadings(vals []string) []string {
	return append(vals, "Obj")
}

func (l *ObjectiveFloatStruct) AppendValues(vals []interface{}) []interface{} {
	return append(vals, l.Curr())
}

func (o *ObjectiveFloatStruct) Converged() Convergence {
	return Converged(o.AbsTol(), o.RelTol())
}

type GradientFloatStruct struct {
	CurrInitGetSetterFloat
	abstol AbsTol
	reltol RelTol
	HistorySaverFloat
	Displayer
}

func (g *GradientFloatStruct) Converged() Convergence {
	return Converged(g.abstol, g.reltol)
}

const DefaultGradAbsTol = 1E-6
const DefaultGradRelTol = 1E-8

func DefaultGradientFloat() *GradientFloatStruct {
	o := &GradientFloatStruct{}
	o.CurrInitGetSetterFloat, o.abstol, o.reltol = NewAbsRelTolStruct(math.NaN(), DefaultGradAbsTol, GradAbsTol, DefaultGradRelTol, GradRelTol)
	o.HistorySaverFloat = DefaultHistorySaverFloat()
	return o
}

func (c *GradientFloatStruct) Initialize() error {
	c.SetCurr(c.Init())
	return nil
}

func (b *GradientFloatStruct) AbsTol() AbsTol {
	return b.abstol
}

func (b *GradientFloatStruct) RelTol() RelTol {
	return b.reltol
}

func (l *GradientFloatStruct) AppendHeadings(vals []string) []string {
	return append(vals, "Grad")
}

func (l *GradientFloatStruct) AppendValues(vals []interface{}) []interface{} {
	return append(vals, math.Abs(l.Curr()))
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
