package gofunopter

// An interface for displaying values 
type Displayer interface {
	DisplayHeadings() []string
	DisplayValues() []interface{}
	Display() *Display
}

func SetDisplay(displayer Displayer) {
	disp := displayer.Display()
	disp.SetHeadings(displayer.DisplayHeadings())
	disp.SetValues(displayer.DisplayValues())
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

type FloatTol struct {
	Name string
	Tol  float64
	curr float64
}
