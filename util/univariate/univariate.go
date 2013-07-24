package univariate

import (
	"fmt"
	"gofunopter/util/display"
	"gofunopter/util/constants"
	"math"
)

// Package note: All of the fields are hidden by default. This is for two main reasons
// First of all, some of them need to be behind getters and setters. For example,
// when setting the current value of a TolFloat, the norm needs to be computed.
// Secondly, many of them are needed for satisfying interfaces. I thought it was
// better to be consistant and have all of them be behind getters and setters

// DefaultObjectiveFloat is the default to be used for a one-D objective.
// Default is that the display is on and
func DefaultObjectiveFloat() *BasicTolFloat {
	return NewBasicTolFloat("Obj", true, math.NaN(), 0, ObjAbsTol, 0, ObjRelTol)
}

func DefaultGradientFloat() *BasicTolFloat {
	return NewBasicTolFloat("Grad", true, math.NaN(), DefaultGradAbsTol, GradAbsTol, DefaultGradRelTol, GradRelTol)
}

func DefaultBoundedStepFloat() *BasicBoundsFloat {
	return NewBasicBoundsFloat("Step", false, 1.0, DefaultBoundedStepFloatAbsTol, StepBoundsAbsTol, DefaultBoundedStepFloatRelTol, StepBoundsRelTol, math.Inf(-1), math.Inf(1))
}

type Location struct {
	*Float
}

// Disp defaults to off, init value defaults to zero
// TODO: Make a Reset() function
func NewLocation() *Location {
	return Location{
		&Tol{
			name: "Loc",
			opt: math.Nan()
		},
	}
}

// Reset resets the values of the struct so that the optimizer can
// be run again (with a different initial condition for example)
func (l *Location) Reset(){
	l.init = 0
	l.hist[:0]
	l.opt = math.NaN()
}


type Objective struct {
	*TolFloat
}

// Disp defaults to off, init value defaults to zero
// Defaults to NaN so that we evaluate at the initial point
// unless set otherwise
// TODO: Make a Reset() function
func NewObjective() *Objective {
	return Location{
		&OptFloat{
			name: "Obj",
			init: math.Nan(),
			disp: true,
		},
	}
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
	return Location{
		&OptFloat{
			name: "Grad",
			init: math.Nan(),
			disp: true,
		},
	}
}

// OptFloat is a float type with the bells and whistles.
// It can be displayed,
// can store a history, can be set to an initial value
// All the normal methods minus the tols
type Float struct {
	saveHist bool
	curr     float64
	init     float64
	hist     []float64
	disp     bool
	name     string
	opt      float64
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

// Hist returns the history of the value over the course of the optimization
// It returns the pointer to the history rather than a copy of the history
// Advanced users can use this to reduce memory cost (by writing the value
// to disk and then reslicing, for example)
func (b *Float) Hist() []float64 {
	return b.hist
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
	return nil
}

// Display returns the display struct to be displayed if
// Float.Disp() is true.
func (b *Float) Display([]display.Struct) []display.Struct {
	return append([]display.Struct, &display.Struct{Value: b.curr, Heading: b.name})
}

// SetResult sets the result. This should be called by the optimizer
// at the end of optimization
func (b *Float) SetResult() {
	b.opt = b.curr
}

// Opt gets the optimimum value at the conclusion of optimization.
func (b *Float) Opt() float64 {
	return b.opt
}

// TolFloat extends Float to allow for absolute and relative
// tolerances to be set
type TolFloat struct {
	*Float
	absTol     float64
	absTolConv Convergence
	relTol     float64
	relTolConv Convergence
	absCurr    float64
	absInit    float64
}

func NewBasicTolFloat(name string, disp bool, init float64, absTol float64,
	absTolConv Convergence, relTol float64, relTolConv Convergence) *BasicTolFloat {
	return &BasicTolFloat{
		Float:      &Float{name: name, disp: disp, init: init},
		absTol:     absTol,
		absTolConv: absTolConv,
		relTol:     relTol,
		relTolConv: relTolConv,
	}
}

func (b *BasicTolFloat) SetInit(val float64) {
	b.init = val
	b.absInit = math.Abs(val)
}

func (b *BasicTolFloat) SetCurr(val float64) {
	b.curr = val
	b.absCurr = math.Abs(val)
}

func (b *BasicTolFloat) SetAbsTol(val float64) {
	b.absTol = val
}

func (b *BasicTolFloat) AbsTol() float64 {
	return b.absTol
}
func (b *BasicTolFloat) SetRelTol(val float64) {
	b.relTol = val
}

func (b *BasicTolFloat) RelTol() float64 {
	return b.relTol
}

func (b *BasicTolFloat) Initialize() error {
	err := b.Float.Initialize()
	if err != nil {
		return err
	}
	b.absCurr = math.Abs(b.curr)
	b.absInit = math.Abs(b.init)
	return nil
}

func (b *BasicTolFloat) Converged() Convergence {

	if b.absCurr < b.absTol {
		return b.absTolConv
	}
	if b.absCurr/b.absInit < b.relTol {
		return b.relTolConv
	}
	return nil
}

// TODO: Implement display bounds
type BasicBoundsFloat struct {
	*BasicTolFloat
	lb      float64
	ub      float64
	gap     float64
	initGap float64
}

func NewBasicBoundsFloat(name string, disp bool, init float64, absTol float64,
	absTolConv Convergence, relTol float64, relTolConv Convergence, lb, ub float64) *BasicBoundsFloat {
	return &BasicBoundsFloat{
		BasicTolFloat: NewBasicTolFloat(name, disp, init, absTol, absTolConv, relTol, relTolConv),
		lb:            lb,
		ub:            ub,
	}
}

func (s *BasicBoundsFloat) Initialize() error {
	s.BasicTolFloat.Initialize()
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

func (s *BasicBoundsFloat) Lb() float64 {
	return s.lb
}

func (s *BasicBoundsFloat) Ub() float64 {
	return s.ub
}

func (s *BasicBoundsFloat) SetLb(val float64) {
	s.lb = val
	s.gap = s.ub - s.lb
}

func (s *BasicBoundsFloat) SetUb(val float64) {
	s.ub = val
	s.gap = s.ub - s.lb
}

func (s *BasicBoundsFloat) AppendHeadings(strs []string) []string {
	return append(strs, s.name+"LB", s.name+"UB")
}

func (s *BasicBoundsFloat) AppendValues(vals []interface{}) []interface{} {
	return append(vals, s.lb, s.lb)
}

// Midpoint between the bounds
func (s *BasicBoundsFloat) Midpoint() float64 {
	return (s.lb + s.ub) / 2.0
}

// Is the value between the upper and lower bounds
func (s *BasicBoundsFloat) WithinBounds(val float64) bool {
	if val < s.lb {
		return false
	}
	if val > s.ub {
		return false
	}
	return true
}

func (b *BasicBoundsFloat) Converged() Convergence {
	if math.IsInf(b.ub, 0) || math.IsInf(b.lb, 0) {
		return nil
	}
	if b.gap < b.absTol {
		return b.absTolConv
	}
	if b.gap/b.initGap < b.absTol {
		return b.relTolConv
	}
	return nil
}
