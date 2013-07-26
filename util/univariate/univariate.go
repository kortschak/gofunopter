package univariate

import (
	"fmt"
	"gofunopter/util/convergence"
	"gofunopter/util/display"
	"math"
)

// Package note: All of the fields are hidden by default. This is for two main reasons
// First of all, some of them need to be behind getters and setters. For example,
// when setting the current value of a TolFloat, the norm needs to be computed.
// Secondly, many of them are needed for satisfying interfaces. I thought it was
// better to be consistant and have all of them be behind getters and setters

// Exception is things that should only need to be set once, and only by the optimizer

// TODO: Reorganize order to be functions we expect the caller to use and
// functions only the optimizer should need

// Reset is something that is set during the optimization and is not a "setting". Idea is
// that optimizer should behave the same way every time you call it, but it assumes you
// don't want to do the same thing twice

//func DefaultBoundedStepFloat() *BoundedFloat {
//	return NewBoundedFloat("Step", false, 1.0, DefaultBoundedStepFloatAbsTol, StepBoundsAbsTol, DefaultBoundedStepFloatRelTol, StepBoundsRelTol, math.Inf(-1), math.Inf(1))
//}

// Needs to be able to call the reset method of float
type Location struct {
	*Float
}

// Disp defaults to off, init value defaults to zero
// TODO: Make a Reset() function
func NewLocation() *Location {
	return &Location{Float: NewFloat("Loc")}
	// Init is zero by default
	// Disp is false by default
}

// Reset resets the values of the struct so that the optimizer can
// be run again (with a different initial condition for example)
func (l *Location) Reset() {
	l.Float.Reset(0)
}

// Objective is an objective type for the optimizer
type Objective struct {
	*TolFloat
}

// Disp defaults to off, init value defaults to zero
// Defaults to NaN so that we evaluate at the initial point
// unless set otherwise
// TODO: Make a Reset() function
func NewObjective() *Objective {
	o := &Objective{
		TolFloat: NewTolFloat("Obj", convergence.ObjAbsTol, convergence.ObjRelTol),
	}
	// Set initial starting value
	o.SetInit(math.NaN())
	o.SetDisp(true)
	// Sets default tolerance behavior
	o.SetAbsTol(0)
	o.SetRelTol(0)
	return o
}

func (o *Objective) Reset() {
	o.TolFloat.Reset(math.NaN())
}

type Gradient struct {
	*TolFloat
}

// Disp defaults to off, init value defaults to zero
// Defaults to NaN so that we evaluate at the initial point
// unless set otherwise
// TODO: Make a Reset() function
// TODO: Add in other defaults
func NewGradient() *Gradient {
	g := &Gradient{
		TolFloat: NewTolFloat("Grad", convergence.GradAbsTol, convergence.GradRelTol),
	}
	g.SetInit(math.NaN())
	g.SetDisp(true)
	// Sets default tolerance behavior
	g.SetAbsTol(convergence.DefaultGradAbsTol)
	// RelTol defaults to zero
	return g
}

func (g *Gradient) Reset() {
	g.TolFloat.Reset(math.NaN())
}

// Gradient should display the norm and not the actual value
func (g *Gradient) Display(d []*display.Struct) []*display.Struct {
	return append(d, &display.Struct{Value: math.Abs(g.Curr()), Heading: "GradNorm"})
}

type Step struct {
	*TolFloat
}

// Disp defaults to off, init value defaults to zero
// Defaults to NaN so that we evaluate at the initial point
// unless set otherwise
// TODO: Make a Reset() function
// TODO: Add in other defaults
func NewStep() *Step {
	s := &Step{
		TolFloat: NewTolFloat("Step", convergence.StepAbsTol, convergence.StepRelTol),
	}
	s.SetInit(1)
	// Disp defaults to off
	// Tolerances default to zero (off)
	return s
}

func (s *Step) Reset() {
	s.TolFloat.Reset(1)
}

// Gradient should display the norm and not the actual value
func (s *Step) Display(d []*display.Struct) []*display.Struct {
	return append(d, &display.Struct{Value: math.Abs(s.Curr()), Heading: "GradNorm"})
}

type BoundedStep struct {
	*BoundedFloat
}

// Ub and Lb are assumed to be called by the optimizer, not the caller
func NewBoundedStep() {
	b := &BoundedStep{
		BoundedFloat: NewBoundedFloat("Step", convergence.StepAbsTol, convergence.StepRelTol),
	}
	b.SetInit(1)
	b.SetUb(math.Inf(1))
	b.SetLb(math.Inf(-1))
	// Disp defaults to false
	// Default both Abs and RelTol to zero
}

func (b BoundedStep) Reset() {
	b.BoundedFloat.Reset(1)
}

// OptFloat is a float type with the bells and whistles.
// It can be displayed,
// can store a history, can be set to an initial value
// All the normal methods minus the tols
// Optimizer should set Name
type Float struct {
	name string

	saveHist bool
	curr     float64
	init     float64
	hist     []float64
	disp     bool
	opt      float64
}

func NewFloat(name string) *Float {
	return &Float{name: name}
}

// AddToHist adds the variable to the history if SetHist() is true
func (b *Float) AddToHist(val float64) {
	if b.saveHist {
		b.hist = append(b.hist, val)
	}
}

// Curr returns the current value of the float
func (b *Float) Curr() float64 {
	return b.curr
}

// SetCurr sets the current value of the variable
func (b *Float) SetCurr(val float64) {
	b.curr = val
}

// Disp returns true if the Float will be displayed during
// optimization and false otherwise (assuming optimization)
// not turned off on the optimizer level)
func (b *Float) Disp() bool {
	return b.disp
}

// SetDisp() sets the display toggle. If true, the value of this
// variable will be displayed during the optimization (assuming the
// optimizer level display toggle is true)
func (b *Float) SetDisp(val bool) {
	b.disp = val
}

// Display returns the display struct to be displayed if
// Float.Disp() is true.
func (b *Float) Display(d []*display.Struct) []*display.Struct {
	return append(d, &display.Struct{Value: b.curr, Heading: b.name})
}

// Hist returns the history of the value over the course of the optimization
// It returns the pointer to the history rather than a copy of the history
// Advanced users can use this to reduce memory cost (by writing the value
// to disk and then reslicing, for example)
func (b *Float) Hist() []float64 {
	return b.hist
}

// Save returns true if all of the values of this variable will be
// stored during optimization
func (b *Float) SaveHist() bool {
	return b.saveHist
}

// SetSaveHist sets the history saver toggle. If true, all of the values
// of the variable will be stored during optimization.
func (b *Float) SetSaveHist(val bool) {
	b.saveHist = val
}

// Init returns the initial value of the variable at the start
// of the optimization
func (b *Float) Init() float64 {
	return b.init
}

// Init sets the initial value of the variable to be
// used by the optimizer. (for example, initial location, initial
// function value, etc.)
func (b *Float) SetInit(val float64) {
	b.init = val
}

// Initialize initializes the Float to be ready to optimize
// This should be called by the optimizer
func (b *Float) Initialize() error {
	if b.hist == nil {
		b.hist = make([]float64, 0)
	}
	b.hist = b.hist[:0]
	b.curr = b.init
	b.opt = math.NaN()
	return nil
}

func (b *Float) Name() string {
	return b.name
}

// Opt gets the optimimum value at the conclusion of optimization.
func (b *Float) Opt() float64 {
	return b.opt
}

// Reset causes the float to be in the condition in which it should be a
// at the initial stage of the optimization.
func (b *Float) Reset(init float64) {
	b.init = init
	b.hist = b.hist[:0]
	b.opt = math.NaN()
}

// SetResult sets the result. This should be called by the optimizer
// at the end of optimization
func (b *Float) SetResult() {
	b.opt = b.curr
}

// TolFloat extends Float to allow for absolute and relative
// tolerances to be set
type TolFloat struct {
	*Float
	absTolConv convergence.C
	relTolConv convergence.C

	absTol  float64
	relTol  float64
	absCurr float64
	absInit float64
}

func NewTolFloat(name string, absTolConv, relTolConv convergence.C) *TolFloat {
	return &TolFloat{
		Float:      NewFloat(name),
		absTolConv: absTolConv,
		relTolConv: relTolConv,
	}
}

// AbsTol returns the current value for the absolute
// tolerance of the norm
// If CurrNorm < AbsTol then the optimizer will have converged
func (b *TolFloat) AbsTol() float64 {
	return b.absTol
}

// SetAbsTol sets the value for the absolute
// tolerance of the norm
// If CurrNorm < AbsTol then the optimizer will have converged
func (b *TolFloat) SetAbsTol(val float64) {
	b.absTol = val
}

// Converged checks if the norm has converged. Returns
// a non-nil value if the norm has fallen below the
// absolute tolerance or the relative tolerance
func (b *TolFloat) Converged() convergence.C {
	if b.absCurr < b.absTol {
		return b.absTolConv
	}
	if b.absCurr < b.relTol*b.absInit {
		return b.relTolConv
	}
	return nil
}

// Gets Curr from Float

// SetCurr sets the current value of TolFloat
func (b *TolFloat) SetCurr(val float64) {
	b.curr = val
	b.absCurr = math.Abs(val)
}

// Initializes. Only the optimizer should need
func (b *TolFloat) Initialize() error {
	err := b.Float.Initialize()
	if err != nil {
		return err
	}
	b.absCurr = math.Abs(b.curr)
	b.absInit = math.Abs(b.init)
	return nil
}

// Gets Init from Float

// SetInit sets the initial value for the optimizer
func (b *TolFloat) SetInit(val float64) {
	b.init = val
	b.absInit = math.Abs(val)
}

// RelTol returns the current value for the relative tolerance
// for the optimizer. If CurrNorm/InitNorm is less than RelTol the
// optimizer will have converged
func (b *TolFloat) RelTol() float64 {
	return b.relTol
}

// RelTol sets the value for the relative tolerance
// for the optimizer. If CurrNorm/InitNorm is less than RelTol the
// optimizer will have converged
func (b *TolFloat) SetRelTol(val float64) {
	b.relTol = val
}

// TODO: Implement display bounds
type BoundedFloat struct {
	*TolFloat
	lb      float64
	ub      float64
	gap     float64
	initGap float64
}

func NewBoundedFloat(name string, absTolConv convergence.C, relTolConv convergence.C) *BoundedFloat {
	return &BoundedFloat{
		TolFloat: NewTolFloat(name, absTolConv, relTolConv),
	}

}

func (s *BoundedFloat) Initialize() error {
	s.TolFloat.Initialize()
	s.initGap = s.ub - s.lb
	s.gap = s.ub - s.lb
	if s.initGap < 0 {
		return fmt.Errorf("Lower bound is greater than upper bound")
	}
	if s.curr <= 0 {
		return fmt.Errorf("Initial step size must be positive")
	}
	return nil
}

func (s *BoundedFloat) Lb() float64 {
	return s.lb
}

func (s *BoundedFloat) Ub() float64 {
	return s.ub
}

func (s *BoundedFloat) SetLb(val float64) {
	s.lb = val
	s.gap = s.ub - s.lb
}

func (s *BoundedFloat) SetUb(val float64) {
	s.ub = val
	s.gap = s.ub - s.lb
}

// Gradient should display the norm and not the actual value
func (b *BoundedFloat) Display(d []*display.Struct) []*display.Struct {
	return append(d, &display.Struct{Value: b.ub, Heading: b.name + "UB"},
		&display.Struct{Value: b.ub, Heading: b.name + "UB"})
}

/*
func (s *BoundedFloat) AppendHeadings(strs []string) []string {
	return append(strs, s.name+"LB", s.name+"UB")
}

func (s *BoundedFloat) AppendValues(vals []interface{}) []interface{} {
	return append(vals, s.lb, s.lb)
}
*/

func (b *BoundedFloat) Reset(init float64) {
	b.TolFloat.Reset(init)
	b.ub = math.Inf(1)
	b.lb = math.Inf(-1)
}

// Midpoint between the bounds
func (s *BoundedFloat) Midpoint() float64 {
	return (s.lb + s.ub) / 2.0
}

// Is the value between the upper and lower bounds
func (s *BoundedFloat) WithinBounds(val float64) bool {
	if val < s.lb {
		return false
	}
	if val > s.ub {
		return false
	}
	return true
}

func (b *BoundedFloat) Converged() convergence.C {
	if math.IsInf(b.ub, 0) || math.IsInf(b.lb, 0) {
		return nil
	}
	if b.gap < b.absTol {
		return b.absTolConv
	}
	if b.gap < b.absTol*b.initGap {
		return b.relTolConv
	}
	return nil
}
