package gofunopter

// Something which can display values
type Displayable interface {
	DisplayHeadings() []string
	DisplayValues() []interface{}
	//Display() *Display
}

type Displayer interface {
	Displayable
	Display() *Display
}

func SetDisplayMethods(displayer Displayer) {
	d := displayer.Display()
	d.GetHeadings = displayer.DisplayHeadings
	d.GetValues = displayer.GetHeadings
}

func SetDisplay(displayable Displayable) {
	disp := displayable.Display()
	disp.SetHeadings(displayable.DisplayHeadings())
	disp.SetValues(displayable.DisplayValues())
}

func AppendValues(values []interface{}, displayables ...Displayable) []interface{} {
	for _, displayable := range displayables {
		newValues := displayable.DisplayValues()
		for _, val := range newValues {
			values = append(currVals, val)
		}
	}
	return values
}

func AppendHeadings(headings []string, displayables ...Displayable) []string {
	for _, displayable := range displayables {
		headings := displayable.DisplayHeadings()
		for _, val := range newHeadings {
			headings := append(headings, val)
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
	CheckConvergence() string
}

func CheckConvergence(convergers ...Converger) (str string) {
	for _, converger := range convergers {
		str = converger.CheckConvergence()
		if str != "" {
			return str
		}
	}
	return ""
}

// Counts up and converges if there is a maximum
type Counter struct {
	Max  int    // Maximum allowable value of the counter
	curr int    // current value of the counter
	Name string // Name of this specific counter
}

func (c *Counter) Add(delta int) {
	c.curr += delta
}

func (c *Counter) CheckConvergence() string {
	if c.curr > c.Max {
		return "Maximum " + c.Name + "reached"
	}
	return ""
}

func (c *Counter) Curr() int {
	return c.curr
}

/*
// Store the history of a float
type HistoryFloat struct {
	Save bool
	hist []float64 // Should this be exposed? Don't really want things messing with it...
}

func DefaultHistoryFloat() *HistoryFloat {
	// Zero length to save memory
	// default save is false
	return &HistoryFloat{hist: make([]float64, 0)}
}

func (f *HistoryFloat) Add(val float) {
	if f.Save {
		f.hist = append(f.hist, val)
	}
}

func (f *FloatHistory) Get() []float64 {
	return f.hist
}

// Store the history of a float
type HistoryFloatSlice struct {
	Save bool
	hist [][]float64 // Should this be exposed? Don't really want things messing with it...
}

func DefaultHistoryFloatSlice() *HistoryFloatSlice {
	// Zero length to save memory
	// default save is false
	return &HistoryFloatSlice{hist: make([][]float64, 0)}
}

func (f *HistoryFloatSlice) Add(val []float) {
	// Copy it in case the slice changes in the future
	newSlice := make([]float64, len(val))
	copy(newSlice, val)
	if f.Save {
		f.hist = append(f.hist, newSlice)
	}
}

func (f *HistoryFloatSlice) Get() [][]float64 {
	return f.hist
}
*/
