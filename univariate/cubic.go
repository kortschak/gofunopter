package univariate

import (
	"errors"
	"github.com/btracey/gofunopter/common/uni"
	"github.com/btracey/gofunopter/convergence"
	"github.com/btracey/gofunopter/optimize"
	"math"

	//"fmt"
)

// A cubic optimizer optimizes by making successive
// cubic approximations to the function value
type Cubic struct {
	step *uni.BoundedStep

	// Tunable parameters
	stepDecreaseMin     float64 // Minimum allowable decrease (must be a number between [0,1)) default 0.0001
	stepDecreaseMax     float64 // When decreasing what is the high
	stepIncreaseMin     float64
	stepIncreaseMax     float64
	floatingpointRelTol float64 // At what point should you stop trusting the objective falue

	// Other needed data during the run
	prevF                     float64 // The f before the current one
	currStepDirectionPositive bool
	initialGradNegative       bool
	deltaCurrent              float64
}

func NewCubic() *Cubic {
	return &Cubic{
		// Cubic has lower bound of 0 for step
		step: uni.NewBoundedStep(),

		// Default Settings
		stepDecreaseMin: 1E-4,
		stepDecreaseMax: 0.9,
		stepIncreaseMin: 1.25,
		stepIncreaseMax: 1E3,

		floatingpointRelTol: 1E-15, // dealing with numerical errors in the function value
	}
}

func (c *Cubic) Initialize(loc *uni.Location, obj *uni.Objective, grad *uni.Gradient) (err error) {

	if c.step.Init() == math.NaN() {
		c.step.SetInit(1)
	}
	// Now initialize the three to set the initial location to the current location
	err = c.step.Initialize()
	if err != nil {
		return errors.New("cubic: error initializing: " + err.Error())
	}

	// Initialize the rest of the memory
	c.initialGradNegative = (grad.Curr() < 0)
	c.currStepDirectionPositive = true
	c.deltaCurrent = 0.0 // How far is the current point from the initial point
	// Add in some checking on the Step Increase and decrease sizes
	return nil
}

func (c *Cubic) Converged() convergence.Type {
	return convergence.CheckConvergence(c.step)
}

func (c *Cubic) SetResult() {
	optimize.SetResult(c.step)
}

func (cubic *Cubic) Iterate(loc *uni.Location, obj *uni.Objective, grad *uni.Gradient, fun UniGradFun) (nFunEvals int, err error) {
	// This will always do one function evaluation per loop
	nFunEvals = 1

	// Initialize
	var stepMultiplier float64
	updateCurrPoint := false
	reverseDirection := false
	currG := grad.Curr()
	currF := obj.Curr()

	// Evaluate trial point
	// Step Size is from the original point
	var trialX float64

	if cubic.initialGradNegative {
		trialX = cubic.step.Curr() + loc.Init()
	} else {
		trialX = -cubic.step.Curr() + loc.Init()
	}
	/*
		fmt.Println(trialX, "trialX")
		fmt.Println("cubic current step", cubic.step.Curr())
	*/
	trialF, trialG, err := fun.Eval(trialX)
	if err != nil {
		return nFunEvals, errors.New("gofunopter: cubic: user defined function error: " + err.Error())
	}

	var newStepSize float64

	loc.AddToHist(trialX)
	obj.AddToHist(trialF)
	grad.AddToHist(trialG)

	/*
		fmt.Println()
		fmt.Println("curr step size", cubic.step.Curr())
		fmt.Println("LB", cubic.step.Lb())
		fmt.Println("UB", cubic.step.Ub())
		fmt.Println("initX", loc.Init())
		fmt.Println("currX", loc.Curr())
		fmt.Println("trialX ", trialX)
		fmt.Println("InitF \t", obj.Init())
		fmt.Println("currF \t", currF)
		fmt.Println("trialF \t", trialF)
		fmt.Println("InitG", grad.Init())
		fmt.Println("currG", currG)
		fmt.Println("trialG", trialG)
	*/

	absTrialG := math.Abs(trialG)

	// Find guess for next point
	deltaF := trialF - currF
	decreaseInValue := (deltaF <= 0)

	// See if we can trust the deltaF measurement
	var canTrustDeltaF bool
	divisor := math.Max(math.Abs(trialF), math.Abs(currF))
	/*
		fmt.Println("Divisor is ", divisor)
		fmt.Println("Delta F", deltaF)
	*/
	if divisor == 0 {
		canTrustDeltaF = true // Both are zero, so is >= 0
	} else {
		if math.Abs(deltaF) > divisor*cubic.floatingpointRelTol {
			// Change large enough to trust
			canTrustDeltaF = true
		}
		// otherwise can't trust
	}

	changeInDerivSign := (currG > 0 && trialG < 0) || (currG < 0 && trialG > 0)
	decreaseInDerivMagnitude := (absTrialG < math.Abs(currG))

	/*
		fmt.Println("Decrease in value ", decreaseInValue)
		fmt.Println("Change in deriv sign ", changeInDerivSign)
		fmt.Println("Decrease in deriv mag ", decreaseInDerivMagnitude)
	*/
	// Find coefficients of the cubic polynomial fit between the current point and the new point
	// Derived from fitting a cubic between (0, CurrF) and (1,TrialF).
	// Apply transformations later to reshift the coordinate axis

	// Need to play games with derivatives
	trialFitG := trialG
	currFitG := currG
	if cubic.initialGradNegative {
		trialFitG *= -1
		currFitG *= -1
	}
	if cubic.currStepDirectionPositive {
		trialFitG *= -1
		currFitG *= -1
	}

	var a, b, c float64
	a = trialG + currG - 2*deltaF
	b = 3*deltaF - 2*currG - trialG

	c = currG
	det := (math.Pow(b, 2) - 3*a*c)

	//fmt.Println("det", det)

	if a == 0 {
		//Perfect quadratic fit
		stepMultiplier = -c / (2 * b)
	} else if det < 0 {
		if decreaseInValue && !changeInDerivSign {
			// The trial point has lower function value
			// and steeper gradient. Set this location as a
			// lower bound for the minimum, and set the
			// next point farther in that direction.

			cubic.setBound("Lower")
			// We know we need to increase in step, but unsure how much, so make a guess
			stepMultiplier = cubic.unclearStepIncrease(loc.Curr())
			if decreaseInDerivMagnitude {
				updateCurrPoint = true
			}
		} else {
			// All other conditions we want to decrease the step size, but the
			// cubic doesn't give an estimate of how much. Just do a binary search
			cubic.setBound("Upper")
			//for i := 0; i < 10; i++ {
			//	fmt.Println("Blah")
			//}
			//fmt.Println(cubic.step.Curr)
			stepMultiplier = 0.5
		}

	} else {
		// Use the cubic projection to guess the minimum location for the line search
		minCubic := (-b + math.Sqrt(det)) / (3 * a)

		/*
		   fmt.Println("sht")
		    fmt.Println("minCubic", minCubic)
		    fmt.Println("SizeLB",cubic.SizeLB)
		    fmt.Println("SizeUB",cubic.SizeUB)
		*/

		switch {
		case changeInDerivSign:
			// There is a change in derivative sign between the current
			// point and the trial point. Make the trial point an upper
			// bound for the minimum, and set it as the current point
			// if the derivative is smaller in magnitude.
			cubic.setBound("Upper")

			stepMultiplier = cubic.sizeDecrease(minCubic)
			if decreaseInDerivMagnitude && decreaseInValue {
				updateCurrPoint = true
				reverseDirection = true
			}

		case (!changeInDerivSign && decreaseInValue) || !canTrustDeltaF:
			// No change in derivative sign, but a decrease in
			// function value. Want to move more in this direction
			// and the trial point is a new lower bound and a new
			// base for the cubic approximation

			cubic.setBound("Lower")

			updateCurrPoint = true

			if decreaseInDerivMagnitude {
				if minCubic < cubic.stepIncreaseMin {
					// Cubic gave a bad approximation (minimum is more
					// in this direction). Assume linear decrease in
					// derivative

					// Check this line
					stepMultiplier = math.Abs(currG) / (math.Abs(trialG) - math.Abs(currG))
				}
				stepMultiplier = cubic.sizeIncrease(stepMultiplier)

			} else {
				// Found a better point, but the derivative increased
				// Use cubic approximation if it gives a reasonable guess
				// otherwise just project forward
				if minCubic < cubic.stepIncreaseMin {
					stepMultiplier = cubic.unclearStepIncrease(loc.Curr())
				} else {
					stepMultiplier = cubic.sizeIncrease(minCubic)
				}

			}
		case (!changeInDerivSign && !decreaseInValue) || canTrustDeltaF:
			// Neither a decrease in value nor a change in derivative sign.
			// This means there must be a local minimum between the starting location and
			// this one. Don't update the cubic point
			cubic.setBound("Upper")
			stepMultiplier = math.Min(minCubic, 0.75)
			stepMultiplier = math.Max(stepMultiplier, 0.25)
		default:
			panic("Bad logic in cases")
		}

	}

	var deltaXTrialCurrent float64
	deltaXTrialCurrent = cubic.step.Curr() - cubic.deltaCurrent
	newDeltaXFromCurrent := deltaXTrialCurrent * stepMultiplier

	newStepSize = newDeltaXFromCurrent + cubic.deltaCurrent

	// Want to make sure that the new search location isn't pushing beyond
	// previously established bounds. If it is, just do a binary search between
	// the bounds
	if !cubic.step.WithinBounds(newStepSize) {
		newStepSize = cubic.step.Midpoint()
	}
	/*
		}

		else {
			fmt.Println("Can't trust deltaF")
			// Can't trust delta F, so just do a binary search
			updateCurrPoint = true
			switch {
			case changeInDerivSign && decreaseInDerivMagnitude:
				cubic.setBound("Upper") // this trial is the new upper bound
				updateCurrPoint = true
			case changeInDerivSign && !decreaseInDerivMagnitude:
				cubic.setBound("Upper") // this trial is the new upper bound
				updateCurrPoint = false
			case !changeInDerivSign && decreaseInDerivMagnitude:
				cubic.setBound("Lower")
				updateCurrPoint = true

			case !changeInDerivSign && !decreaseInDerivMagnitude:
				return nFunEvals, errors.New("Gradient tolerance too high for the precision of the gradient")
			}
			newStepSize = cubic.step.Midpoint()
		}
	*/

	if updateCurrPoint {

		loc.SetCurr(trialX)
		cubic.prevF = obj.Curr()
		obj.SetCurr(trialF)
		grad.SetCurr(trialG)
		cubic.deltaCurrent = trialX - loc.Init()
		if cubic.initialGradNegative {
			cubic.deltaCurrent *= -1
		}
		if reverseDirection {
			cubic.currStepDirectionPositive = !cubic.currStepDirectionPositive
		}
	}
	cubic.step.SetCurr(newStepSize)
	return nFunEvals, nil
}

func (c *Cubic) unclearStepIncrease(currLoc float64) float64 {
	// Increase the step. If there is an upper bound, do a binary
	// search between the upper and lower bound. If there is no
	// upper bound, just double the step size.

	// Clean up this code!
	var stepMultiplier float64
	if c.currStepDirectionPositive {
		if math.IsInf(c.step.Ub(), 1) {
			stepMultiplier = (2*c.step.Curr() - c.deltaCurrent) / (c.step.Curr() - c.deltaCurrent)
		} else {
			stepMultiplier = (c.step.Midpoint() - c.deltaCurrent) / (c.step.Curr() - c.deltaCurrent)
		}
	} else {
		if math.IsInf(c.step.Lb(), -1) {
			stepMultiplier = (2*c.step.Curr() - currLoc) / (c.step.Curr() - currLoc)
		} else {
			stepMultiplier = (c.step.Midpoint() - currLoc) / (c.step.Curr() - currLoc)
		}
	}
	return stepMultiplier
}

func (c *Cubic) setBound(dir string) {
	if dir == "Lower" {
		// Want to go farther in this direction
		if c.currStepDirectionPositive == true {
			c.step.SetLb(c.step.Curr())
		} else {
			c.step.SetUb(c.step.Curr())
		}
	} else if dir == "Upper" {
		if c.currStepDirectionPositive == true {
			c.step.SetUb(c.step.Curr())
		} else {
			c.step.SetLb(c.step.Curr())
		}
	}
	return
}

func (c *Cubic) sizeDecrease(minCubic float64) float64 {

	stepMultiplier := math.Max(minCubic, c.stepDecreaseMin)
	stepMultiplier = math.Min(stepMultiplier, c.stepDecreaseMax)
	return stepMultiplier
}

func (c *Cubic) sizeIncrease(minCubic float64) float64 {
	stepMultiplier := math.Max(minCubic, c.stepIncreaseMin)
	stepMultiplier = math.Min(stepMultiplier, c.stepIncreaseMax)
	return stepMultiplier
}
