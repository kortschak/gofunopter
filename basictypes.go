package gofunopter

import (
	"math"
)

// Something which can display values
type Displayable interface {
	DisplayHeadings() []string
	DisplayValues() []interface{}
	//Display() *Display
}

type Displayer interface {
	Displayable
	GetDisplay() *Display
}

func SetDisplayMethods(displayer Displayer) {
	d := displayer.GetDisplay()
	d.GetHeadings = displayer.DisplayHeadings
	d.GetValues = displayer.DisplayValues
}

/*
func SetDisplay(displayable Displayable) {
	disp := displayable.Display()
	disp.SetHeadings(displayable.DisplayHeadings())
	disp.SetValues(displayable.DisplayValues())
}
*/

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
	Iterate()
}

func Iterate(iterators ...Iterator) {
	for _, iterator := range iterators {
		iterator.Iterate()
	}
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
	if o.Curr < o.AbsTol {
		return o.Name + " absolute tolerance reached"
	}
	if o.Curr/o.absinit < o.RelTol {
		return o.Name + " relative tolerance reached"
	}
	return ""
}

func (o *OptFloat) Result() {
	o.Opt = o.Curr
}

// Returns the default values for an input location
// Locations don't have any tolerances
func DefaultInputFloat() *OptFloat {
	return &OptFloat{
		AbsTol:   math.Inf(-1),
		RelTol:   math.Inf(-1),
		Name:     "loc",
		SaveHist: false,
		Init:     0,
	}
}

// Returns the default values for an objective
// Objectives generally don't have any tolerances
// No idea what the initial function value is
func DefaultObjectiveFloat() *OptFloat {
	return &OptFloat{
		AbsTol:   math.Inf(-1),
		RelTol:   math.Inf(-1),
		Name:     "fun",
		SaveHist: false,
		Init:     math.NaN(),
	}
}

// Returns the default values for the gradient
func DefaultGradientFloat() *OptFloat {
	return &OptFloat{
		AbsTol:   1E-6,
		RelTol:   1E-8,
		Name:     "grad",
		SaveHist: false,
		Init:     math.NaN(),
	}
}

type StepFloat struct {
	*OptFloat
	Lb float64
	Ub float64
}

// Midpoint between the bounds
func (s *StepFloat) Midpoint() float64 {
	return (s.Lb + s.Ub) / 2.0
}

// Is the value between the upper and lower bounds
func (s *StepFloat) WithinBounds(val float64) bool {
	if val < s.Lb {
		return false
	}
	if val > s.Ub {
		return false
	}
	return true
}

func (s *StepFloat) Converged() string {
	if (s.Ub - s.Lb) < s.AbsTol {
		return s.Name + " absolute tolerance reached"
	}
	return ""
}

// Returns the default values for a step size
// no default relative tolerance
func DefaultStepFloat() *StepFloat {
	return &StepFloat{
		OptFloat: &OptFloat{
			AbsTol:   1E-6,
			RelTol:   math.Inf(-1),
			Name:     "step",
			SaveHist: false,
			Init:     1,
		},
		Lb: math.Inf(-1),
		Ub: math.Inf(1),
	}
}
