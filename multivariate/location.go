package multivariate

// Location is the basic float slice type representing a multi-D location value
// It's name is "Loc", and the initial value is set to be nil (it is not possible
// to know how many dimensions a-priori)
type Location struct {
	*Floats
}

// Disp defaults to off, init value defaults to zero
func NewLocation() *Location {
	return &Location{Floats: NewFloat("Loc", false)}
	// Init is zero by default
	// Disp is false by default
}

// Sets the result at the end of the optimaziton (value found in Opt()),
// and resets the initial value to zero
func (l *Location) SetResult() {
	l.Floats.SetResult()
}
