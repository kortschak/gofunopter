package optimize

// This is in its own package because otherwise it's hard to avoid circular imports

import (
	"errors"
	//"fmt"
	"github.com/btracey/gofunopter/convergence"
	"github.com/btracey/gofunopter/display"
)

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
	common.CommonInitialize()

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
	defer common.CommonSetResult()

	// Main optimization loop:
	// Iterate until convergence, outputting the display as we go (assuming)
	// appropriate booleans are true

	converger, isConverger := fun.(convergence.Converger)
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
		c = common.CommonConverged()
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
	optDisplay.Reset()
	// SetResult will occur in the unwinding
	return c, nil
}

func DisplayOpter(optDisplay *display.Display, o Optimizer, common, displayer display.Displayer, isDisplayer bool) {
	// Display the outputs (if toggle is on)
	if o.Disp() {
		if isDisplayer {
			optDisplay.DisplayProgress(common, o, displayer)
		} else {
			optDisplay.DisplayProgress(common, o)
		}
	}
}
