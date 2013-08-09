package optimize

// This is in its own package because otherwise it's hard to avoid circular imports

import (
	"errors"
	"github.com/btracey/gofunopter/common"
	"github.com/btracey/gofunopter/common/convergence"
	"github.com/btracey/gofunopter/common/display"
)

type Optimizer interface {
	Initializer
	//SetResulter
	convergence.Converger
	display.Displayer
	Iterate() error
	GetOptCommon() *common.OptCommon
	SetSettings() error
	CommonSettings() *common.CommonSettings
	SetResult(*common.CommonResult)
}

// OptimizeOpter is the basic method for using optimizers. Not intended to
// be called by the user
func OptimizeOpter(o Optimizer, fun interface{}) (convergence.Type, error) {
	var err error
	// Set all the settings
	commonSettings := o.CommonSettings()

	// Initialize the common value
	common := o.GetOptCommon()
	common.SetSettings(commonSettings)

	o.SetSettings()

	// Initialize the caller's function if it is an initializer
	initer, ok := fun.(Initializer)
	if ok {
		err = initer.Initialize()
		if err != nil {
			return nil, errors.New("opt: error during user defined function initialization: " + err.Error())
		}
	}
	common.CommonInitialize()

	// Initialize the optimizer
	err = o.Initialize()
	if err != nil {
		return nil, errors.New("opt: error during optimizer initialization, " + err.Error())
	}

	// Get the displayer
	optDisplay := common.Display

	// Defer call to set result
	// Want to return the result even if there is an error (don't want to waste function
	// evaluations if the caller can handle the error)
	// Want to defer call to user-defined set result first to it can unwind after the
	// optimizer does

	defer SetOptResults(o, common, fun)
	//defer o.SetResult()
	//defer common.CommonSetResult()

	// Main optimization loop:
	// Iterate until convergence, outputting the display as we go (assuming)
	// appropriate booleans are true

	converger, isConverger := fun.(convergence.Converger)
	displayer, isDisplayer := fun.(display.Displayer)

	var c convergence.Type
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
		err := o.Iterate()
		common.Iter.Add(1)

		// Function evaluations and history dealt with locally

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
	if optDisplay.Disp {
		if isDisplayer {
			optDisplay.DisplayProgress(common, o, displayer)
		} else {
			optDisplay.DisplayProgress(common, o)
		}
	}
}

func SetOptResults(o Optimizer, c *common.OptCommon, fun interface{}) {
	commonResult := c.CommonResult()
	o.SetResult(commonResult)
	setResulter, ok := fun.(SetResulter)
	if ok {
		setResulter.SetResult()
	}
}
