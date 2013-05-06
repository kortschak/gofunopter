package gofunopter

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
		return "Maximum " + c.name + "reached"
	}
	return ""
}

func (c *Counter) Curr() int {
	return curr
}

type FloatMax struct {
	Max  float64
	curr float64
	Name string
}

func (f *FloatMax) Curr() float64 {
	return f.curr
}

func (f *FloatMax) Set() float64 {
	return f.curr
}

func (f *FloatMax) CheckConvergence() string {
	if curr > Max {
		return "Maximum " + f.Name + "reached"
	}
}

type FloatTol struct {
	Name string
	Tol  float64
	curr float64
}
