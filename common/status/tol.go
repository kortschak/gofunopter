package status

const DefaultGradAbsTol = 1E-6
const DefaultStepAbsTol = 1E-15 // any smaller and numerical issues happen

// AbsTol is a structure for representing an absolute tolerance.
// Converges if curr < tol
type Abs struct {
	tol  float64
	conv Status
}

func NewAbs(tol float64, conv Status) *Abs {
	return &Abs{tol: tol, conv: conv}
}

// AbsTol returns the value of the absolute tolerance
func (a *Abs) AbsTol() float64 {
	return a.tol
}

// SetAbsTol sets the value of the absolute tolerance
func (a *Abs) SetAbsTol(tol float64) {
	a.tol = tol
}

// CheckConvergence checks if the absolute tolerance has been reached
func (a *Abs) Status(curr float64) Status {
	if curr < a.tol {
		return a.conv
	}
	return Continue
}

// RelTol is a structure for representing an absolute tolerance.
// Converges if curr < tol * init
type Rel struct {
	tol  float64
	conv Status
}

func NewRel(tol float64, conv Status) *Rel {
	return &Rel{tol: tol, conv: conv}
}

// AbsTol returns the value of the absolute tolerance
func (r *Rel) RelTol() float64 {
	return r.tol
}

// SetAbsTol sets the value of the absolute tolerance
func (r *Rel) SetRelTol(tol float64) {
	r.tol = tol
}

// CheckConvergence checks if the relative tolerance has been reached
func (r *Rel) Status(curr, init float64) Status {
	if curr < r.tol*init {
		return r.conv
	}
	return Continue
}
