package gofunopter

import (
	"errors"
	"gofunopter/common"
	"gofunopter/convergence"
	"gofunopter/display"
)

type Optimizer interface {
	Initializer
	SetResulter
	Converger
	display.Displayer

	Iterate() (int, error)
	GetOptCommon() *common.OptCommon
	GetDisplay() *display.Display
}

type Initializer interface {
	Initialize() error
}

type SetResulter interface {
	SetResult()
}

type Converger interface {
	Converged() convergence.C
}

// OptimizeOpter is the basic method for using optimizers. Not intended to
// be called by the user
func OptimizeOpter(o Optimizer, fun interface{}) (convergence.C, error) {
	var err error
	// Initialize the caller's function if it is an initializer
	initer, ok := fun.(Initializer)
	if ok {
		err = initer.Initialize()
		if err != nil {
			return nil, errors.New("opt: errer during user defined function initialization")
		}
	}

	// Initialize the common value
	common := o.GetOptCommon()
	common.Initialize()

	// Get the displayer
	optDisplay := o.GetDisplay()

	// Initialize the optimizer
	err = o.Initialize()
	if err != nil {
		return nil, errors.New("opt: error during optimizer initialization, " + err.Error())
	}

	// Defer call to set result
	// Want to return the result even if there is an error (don't want to waste function
	// evaluations if the caller can handle the error)
	// Want to defer call to user-defined set result first to it can unwind after the
	// optimizer does
	setResulter, ok := fun.(SetResulter)
	if ok {
		setResulter.SetResult()
	}
	defer o.SetResult()

	// Main optimization loop:
	// Iterate until convergence, outputting the display as we go (assuming)
	// appropriate booleans are true

	converger, isConverger := fun.(Converger)
	displayer, isDisplayer := fun.(display.Displayer)

	var c convergence.C
	for {
		// Check if the optimizer has converged
		c = o.Converged()
		if c != nil {
			break
		}

		// Check if the user-defined function has converged
		if isConverger {
			c = converger.Converged()
			if c != nil {
				break
			}
		}

		// Check if common has converged (iterations, funevals, etc.)
		c = common.Converged()
		if c != nil {
			break
		}

		// Display the outputs (if toggle is on)
		DisplayOpter(optDisplay, o, common, displayer, isDisplayer)

		// If the optimizer has not converged, take an iteration
		// in the optimizer
		nFunEvals, err := o.Iterate()
		common.Iter().Add(1)
		common.FunEvals().Add(nFunEvals)
		if err != nil {
			return nil, errors.New("opt: Error during optimizer iteration, " + err.Error())
		}
	}
	// Display at end of optimization
	optDisplay.IncreaseValueTime()
	DisplayOpter(optDisplay, o, common, displayer, isDisplayer)
	// SetResult will occur in the unwinding
	return c, nil
}

func DisplayOpter(optDisplay *display.Display, o, common, displayer display.Displayer, isDisplayer bool) {
	// Display the outputs (if toggle is on)
	if o.Disp() {
		if isDisplayer {
			optDisplay.DisplayProgress(o, common, displayer)
		} else {
			optDisplay.DisplayProgress(o, common)
		}
	}
}
