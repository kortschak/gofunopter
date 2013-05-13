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

type OptimizerFailure struct {
	Str string
	Err error
}

func (o OptimizerFailure) ConvergenceType() string {
	return o.Str + o.Err.Error()
}

func (o OptimizerFailure) Error() string {
	return o.Str + o.Err.Error()
}

func InitializationError(err error) OptimizerFailure {
	return OptimizerFailure{
		Str: "Optimizer failed to initialize, ",
		Err: err,
	}
}

// Counts up and converges if there is a maximum
type Counter struct {
	max   int // Maximum allowable value of the counter
	curr  int // current value of the counter
	total int // Total number at the end of the optimization run
	conv  Convergence
	name  string
	disp  bool
}

// No counter iterate method because we don't know how many we want to add on each iteration

func NewCounter(name string, max int, conv Convergence, disp bool) *Counter {
	return &Counter{
		name: name,
		max:  max,
		conv: conv,
		disp: disp,
	}
}

func (c *Counter) Disp() bool {
	return c.disp
}

func (c *Counter) SetDisp(val bool) {
	c.disp = val
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
func (c *Counter) Converged() Convergence {
	if c.curr > c.max {
		return c.conv
	}
	return nil
}

func (c *Counter) SetResult() {
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

func (b *BasicOptFloat) Save() bool {
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

func (b *BasicOptFloat) SetCurr(val float64) {
	b.curr = val
}

func (b *BasicOptFloat) Init() float64 {
	return b.init
}

func (b *BasicOptFloat) SetInit(val float64) {
	b.init = val
}

func (b *BasicOptFloat) Initialize() error {
	b.hist = make([]float64, 0)
	b.curr = b.init
	return nil
}

func (b *BasicOptFloat) AppendHeadings(headings []string) []string {
	headings = append(headings, b.name)
	return headings
}

func (b *BasicOptFloat) AppendValues(vals []interface{}) []interface{} {
	vals = append(vals, b.curr)
	return vals
}

func (b *BasicOptFloat) SetResult() {
	b.opt = b.curr
}

func (b *BasicOptFloat) Opt() float64 {
	return b.opt
}

type BasicTolFloat struct {
	*BasicOptFloat
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
		BasicOptFloat: &BasicOptFloat{name: name, disp: disp, init: init},
		absTol:        absTol,
		absTolConv:    absTolConv,
		relTol:        relTol,
		relTolConv:    relTolConv,
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
	err := b.BasicOptFloat.Initialize()
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
	if b.absCurr/b.absInit < b.absTol {
		return b.relTolConv
	}
	return nil
}

var LocAbsTol Convergence = LocConvergence{"Location absolute tolerance reached"}
var LocRelTol Convergence = LocConvergence{"Location relative tolerance reached"}
var ObjAbsTol Convergence = FunConvergence{"Function absolute tolerance reached"}
var ObjRelTol Convergence = FunConvergence{"Function relative tolerance reached"}
var GradAbsTol Convergence = GradConvergence{"Gradient absolute tolerance reached"}
var GradRelTol Convergence = GradConvergence{"Gradient relative tolerance reached"}
var StepAbsTol Convergence = StepConvergence{"Step absolute tolerance reached"}
var StepRelTol Convergence = StepConvergence{"Step relative tolerance reached"}
var StepBoundsAbsTol Convergence = StepConvergence{"Step bounds absolute tolerance reached"}
var StepBoundsRelTol Convergence = StepConvergence{"Step bounds absolute tolerance reached"}

const DefaultGradAbsTol = 1E-6
const DefaultGradRelTol = 1E-8

func DefaultLocationFloat() *BasicOptFloat {
	return NewBasicOptFloat("Loc", false, 0)
}

func DefaultObjectiveFloat() *BasicTolFloat {
	return NewBasicTolFloat("Obj", true, math.NaN(), 0, ObjAbsTol, 0, ObjRelTol)
}

func DefaultGradientFloat() *BasicTolFloat {
	return NewBasicTolFloat("Grad", true, math.NaN(), DefaultGradAbsTol, GradAbsTol, DefaultGradRelTol, GradRelTol)
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
	s.curr = s.init
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

const DefaultInitStepSize = 1
const DefaultBoundedStepFloatAbsTol = 0 //
const DefaultBoundedStepFloatRelTol = 0 // Turn off step rel tol

func DefaultBoundedStepFloat() *BasicBoundsFloat {
	return NewBasicBoundsFloat("Step", false, 1, DefaultBoundedStepFloatAbsTol, StepBoundsAbsTol, DefaultBoundedStepFloatRelTol, StepBoundsRelTol, math.Inf(-1), math.Inf(1))
}
