package gofunopter

import (
	"fmt"
	"github.com/btracey/smatrix"
	"math"
)

//TODO: Be more careful with resetting and error checking during optimization

var _ = fmt.Println

func AppendValues(values []interface{}, displayables ...Displayable) []interface{} {
	for _, displayable := range displayables {
		newValues := displayable.DisplayValues()
		for _, val := range newValues {
			values = append(values, val)
		}
	}
	return values
}

func AppendHeadings(headings []string, displayables ...Displayable) []string {
	for _, displayable := range displayables {
		newHeadings := displayable.DisplayHeadings()
		for _, val := range newHeadings {
			headings = append(headings, val)
		}
	}
	return headings
}

type Iterator interface {
	Iterate() error
}

func Iterate(iterators ...Iterator) (err error) {
	for _, iterator := range iterators {
		err := iterator.Iterate()
		if err != nil {
			return err
		}
	}
	return nil
}

type Converger interface {
	Converged() string
}

func Converged(convergers ...Converger) (str string) {
	for _, converger := range convergers {
		str = converger.Converged()
		if str != "" {
			return str
		}
	}
	return ""
}

type Resulter interface {
	Result()
}

func SetResults(resulters ...Resulter) {
	for _, resulter := range resulters {
		resulter.Result()
	}
}

// Counts up and converges if there is a maximum
type Counter struct {
	Max   int    // Maximum allowable value of the counter
	curr  int    // current value of the counter
	Name  string // Name of this specific counter
	Total int    // Total number at the end of the optimization run
}

func (c *Counter) Add(delta int) {
	c.curr += delta
}

func (c *Counter) Converged() string {
	if c.curr > c.Max {
		return "Maximum " + c.Name + "reached"
	}
	return ""
}

func (c *Counter) Result() {
	c.Total = c.curr
}

func (c *Counter) Curr() int {
	return c.curr
}

// Something about only major iterations?
type HistoryFloat struct {
	hist []float64
	Save bool
}

func (h *HistoryFloat) Get() []float64 {
	return h.hist
}

func (h *HistoryFloat) Set(val []float64) {
	h.hist = val
}

func (h *HistoryFloat) Add(val float64) {
	if h.Save {
		h.hist = append(h.hist, val)
	}
}

// A float type which can check tolerances, initialize and save history
type OptFloat struct {
	Curr     float64       // current value of the float
	Init     float64       // initial value of the float
	absinit  float64       // The absolute value of the initial 
	Hist     *HistoryFloat // stored history of the float
	AbsTol   float64       // Tolerance on that value
	SaveHist bool          // Do you want to save the history. Called save because maybe we want to have a save to file in the future
	RelTol   float64       // Tolerance relative to the initial value
	Name     string        // Name of the OptFloat
	Opt      float64       // Optimal value at the end of the run
	Disp     bool          // Display this output during teh run if display is on 
}

// Initializes by setting the current value to the initial value and
// appending it to the history if necessary
func (o *OptFloat) Initialize() {
	o.Curr = o.Init
	o.absinit = math.Abs(o.Init)
	if o.Hist.Get() == nil {
		o.Hist.Set(make([]float64, 0))
	}
	if o.Hist.Save {
		o.Hist.Add(o.Init)
	}

}

// Make hist return a copy? Have a CopyHist method?

func (o *OptFloat) Converged() string {
	if math.Abs(o.Curr) < o.AbsTol {
		return o.Name + " absolute tolerance reached"
	}
	if math.Abs(o.Curr)/o.absinit < o.RelTol {
		return o.Name + " relative tolerance reached"
	}
	return ""
}

func (o *OptFloat) Result() {
	o.Opt = o.Curr
}

func (o *OptFloat) DisplayHeadings() []string {
	return []string{o.Name}
}

func (o *OptFloat) DisplayValues() []interface{} {
	return []interface{}{o.Curr}
}

// Returns the default values for an input location
// Locations don't have any tolerances

type LocationFloat struct {
	*OptFloat
}

func DefaultLocationFloat() *LocationFloat {
	return &LocationFloat{
		OptFloat: &OptFloat{
			AbsTol:   math.Inf(-1),
			RelTol:   math.Inf(-1),
			Name:     "loc",
			SaveHist: false,
			Init:     0,
			Disp:     true,
			Hist:     &HistoryFloat{},
		},
	}
}

// Returns the default values for an objective
// Objectives generally don't have any tolerances
// No idea what the initial function value is

type ObjectiveFloat struct {
	*OptFloat
}

func DefaultObjectiveFloat() *ObjectiveFloat {
	o := &ObjectiveFloat{
		OptFloat: &OptFloat{
			AbsTol:   math.Inf(-1),
			RelTol:   math.Inf(-1),
			Name:     "fun",
			SaveHist: false,
			Init:     math.NaN(),
			Curr:     math.NaN(),
			Disp:     true,
			Hist:     &HistoryFloat{},
		},
	}
	return o
}

type GradientFloat struct {
	*OptFloat
}

// Returns the default values for the gradient
func DefaultGradientFloat() *GradientFloat {
	return &GradientFloat{
		OptFloat: &OptFloat{
			AbsTol:   1E-6,
			RelTol:   1E-8,
			Name:     "grad",
			SaveHist: false,
			Init:     math.NaN(),
			Curr:     math.NaN(),
			Disp:     true,
			Hist:     &HistoryFloat{},
		},
	}
}

type BoundedFloat struct {
	*OptFloat
	Lb         float64
	Ub         float64
	DispBounds bool
}

func (s *BoundedFloat) DisplayHeadings() []string {
	strs := s.OptFloat.DisplayHeadings()
	if s.DispBounds {
		strs = append(strs, s.Name+"LB", s.Name+"UB")
	}
	return strs
}

func (s *BoundedFloat) DisplayValues() []interface{} {
	vals := s.OptFloat.DisplayValues()
	if s.DispBounds {
		vals = append(vals, s.Lb, s.Ub)
	}
	return vals
}

// Midpoint between the bounds
func (s *BoundedFloat) Midpoint() float64 {
	return (s.Lb + s.Ub) / 2.0
}

// Is the value between the upper and lower bounds
func (s *BoundedFloat) WithinBounds(val float64) bool {
	if val < s.Lb {
		return false
	}
	if val > s.Ub {
		return false
	}
	return true
}

func (s *BoundedFloat) Converged() string {
	if (s.Ub - s.Lb) < s.AbsTol {
		return s.Name + " absolute tolerance reached"
	}
	return ""
}

//TODO: Change this to have a real default bounded float rather than just the step

// Returns the default values for a step size
// no default relative tolerance
func DefaultStepFloat() *BoundedFloat {
	return &BoundedFloat{
		OptFloat: &OptFloat{
			AbsTol:   1E-6,
			RelTol:   math.Inf(-1),
			Name:     "step",
			SaveHist: false,
			Init:     1,
			Disp:     true,
			Hist:     &HistoryFloat{},
		},
		Lb:         math.Inf(-1),
		Ub:         math.Inf(1),
		DispBounds: false,
	}
}

// Something about only major iterations?
type HistoryFloatSlice struct {
	hist [][]float64
	Save bool
}

func (h *HistoryFloatSlice) Get() [][]float64 {
	return h.hist
}

func (h *HistoryFloatSlice) Set(val [][]float64) {
	h.hist = val
}

func (h *HistoryFloatSlice) Add(val []float64) {
	// Make a copy in case the pointer changes
	if h.Save {
		newVal := make([]float64, len(val))
		copy(newVal, val)
		h.hist = append(h.hist, newVal)
	}
}

// A float type which can check tolerances, initialize and save history
type OptFloatSlice struct {
	Curr     []float64          // current value of the float
	Init     []float64          // initial value of the float
	norminit float64            // The norm of the initial point
	Hist     *HistoryFloatSlice // stored history of the float
	AbsTol   float64            // Tolerance on the norm of the value
	SaveHist bool               // Do you want to save the history. Called save because maybe we want to have a save to file in the future
	RelTol   float64            // Tolerance relative to the norm of the initial value
	Name     string             // Name of the OptFloat
	Opt      []float64          // Optimal value at the end of the run
	Disp     bool               // Display this output during teh run if display is on 
}

// Initializes by setting the current value to the initial value and
// appending it to the history if necessary
func (o *OptFloatSlice) Initialize() {
	newInit := make([]float64, len(o.Curr))
	o.Curr = newInit
	o.norminit = smatrix.VectorTwoNorm(o.Init)
	if o.Hist.Get() == nil {
		o.Hist.Set(make([][]float64, 0))
	}
	if o.Hist.Save {
		o.Hist.Add(newInit)
	}

}

// Make hist return a copy? Have a CopyHist method?

func (o *OptFloatSlice) Converged() string {
	norm := smatrix.VectorTwoNorm(o.Curr)
	if norm < o.AbsTol {
		return o.Name + " absolute tolerance reached"
	}
	if norm/o.norminit < o.RelTol {
		return o.Name + " relative tolerance reached"
	}
	return ""
}

func (o *OptFloatSlice) Result() {
	o.Opt = o.Curr
}

func (o *OptFloatSlice) DisplayHeadings() []string {
	return []string{o.Name}
}

func (o *OptFloatSlice) DisplayValues() []interface{} {
	return []interface{}{smatrix.VectorTwoNorm(o.Curr)}
}

type LocationFloatSlice struct {
	*OptFloatSlice
}

func DefaultLocationFloatSlice() *LocationFloatSlice {
	return &LocationFloatSlice{
		OptFloatSlice: &OptFloatSlice{
			AbsTol:   math.Inf(-1),
			RelTol:   math.Inf(-1),
			Name:     "loc",
			SaveHist: false,
			Disp:     true,
			Hist:     &HistoryFloatSlice{},
		},
	}
}

type GradientFloatSlice struct {
	*OptFloatSlice
}

// Returns the default values for the gradient
func DefaultGradientFloatSlice() *GradientFloatSlice {
	return &GradientFloatSlice{
		OptFloatSlice: &OptFloatSlice{
			AbsTol:   1E-6,
			RelTol:   1E-8,
			Name:     "grad",
			SaveHist: false,
			Disp:     true,
			Hist:     &HistoryFloatSlice{},
		},
	}
}
