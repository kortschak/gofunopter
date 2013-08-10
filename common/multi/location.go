package multi

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

type LocationSettings struct {
	DisplayLocation     bool
	KeepLocationHistory bool
}

func NewLocationSettings() *LocationSettings {
	return &LocationSettings{
		DisplayLocation: false,
	}
}

type LocationResult struct {
	LocationHistory [][]float64
	Location        []float64
}

func (o *Location) SetSettings(s *LocationSettings) {
	o.SetDisp(s.DisplayLocation)
	o.SetSaveHist(s.KeepLocationHistory)
}

// SetResult sets the optimum value, and resets the initial value to NaN
func (o *Location) Result() *LocationResult {
	return &LocationResult{
		LocationHistory: o.Floats.Hist(),
		Location:        o.Floats.Opt(),
	}
}
