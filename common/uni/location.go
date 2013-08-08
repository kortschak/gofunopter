package uni

// Location is the basic float type representing a one-D location value
// It's name is "Loc", and the initial value is set to be 0 (the starting
// location for the optimizer unless the user sets a alternative starting
// location). The default is to not display
type Location struct {
	*Float
}

// Disp defaults to off, init value defaults to zero
func NewLocation() *Location {
	return &Location{Float: NewFloat("Loc", false)}
	// Init is zero by default
	// Disp is false by default
}

// Sets the result at the end of the optimaziton (value found in Opt()),
// and resets the initial value to zero
func (l *Location) SetResult() {
	l.Float.SetResult()
}
