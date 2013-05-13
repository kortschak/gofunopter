package gofunopter

import (
	"fmt"
	"math"
)

// Maybe all of these should go in their own subpackages?
// SISOGradBased, etc. That way have SISO.Optimizable?

type Cubic struct {
	// Basic values (should these be interfaces?)
	Loc  OptFloat        // Location
	Obj  OptTolFloat     // Function Value
	Grad OptTolFloat     // Gradient value
	Step BoundedOptFloat // Step size
	*Common
	Fun SISOGradBasedProblem

	// Tunable parameters
	StepDecreaseMin float64 // Minimum allowable decrease (must be a number between [0,1)) default 0.0001
	StepDecreaseMax float64 // When decreasing what is the high
	StepIncreaseMin float64
	StepIncreaseMax float64
	InitStepSize    float64

	// Other needed data during the run
	currStepDirectionPositive bool
	initialGradNegative       bool
	deltaCurrent              float64
}

func DefaultCubic() *Cubic {
	c := &Cubic{
		Loc:             DefaultLocationFloat(),
		Obj:             DefaultObjectiveFloat(),
		Grad:            DefaultGradientFloat(),
		Step:            DefaultBoundedStepFloat(),
		StepDecreaseMin: 1E-4,
		StepDecreaseMax: 0.9,
		StepIncreaseMin: 1.25,
		StepIncreaseMax: 1E3,
		Common:          DefaultCommon(),
	}
	SetDisplayMethods(c)
	return c
}

// Should we add error checking to the evaluations?
func (c *Cubic) Initialize() (err error) {
	// Initialize takes all of these in so function evaluations can be saved if 
	// the information is already there
	c.Common.Initialize()
	c.Loc.Initialize()
	fmt.Println("initialize")
	fmt.Println(c.Obj.Init())
	if math.IsNaN(c.Obj.Init()) || math.IsNaN(c.Grad.Curr()) {
		fmt.Println("In isnan")
		// Initial function value hasn't been set, so do it.
		err = c.Fun.Eval(c.Loc.Init())
		if err != nil {
			return fmt.Errorf("Error evaluating the function at the set initial value %v", c.Loc.Curr)
		}
		c.FunEvals.Add(1)
		c.Obj.SetInit(c.Fun.Obj())
		c.Grad.SetInit(c.Fun.Grad())
	}
	c.Obj.Initialize()
	c.Grad.Initialize()
	if c.Step.Init() <= 0 {
		return fmt.Errorf("Initial step must be positive")
	}
	c.initialGradNegative = (c.Grad.Curr() < 0)
	c.currStepDirectionPositive = true
	c.deltaCurrent = 0.0 // How far is the current point from the initial point

	// Add in some checking on the Step Increase and decrease sizes
	return nil
}

func (c *Cubic) Converged() Convergence {
	conv := Converged(c.Obj, c.Grad, c.Step)
	if conv != nil {
		return conv
	}
	s, ok := c.Fun.(Converger)
	if ok {
		conv = s.Converged()
		if conv != nil {
			return conv
		}
	}
	return nil
}

func (c *Cubic) AppendHeadings(headings []string) []string {
	headings = AppendHeadings(headings, c.Common, c.Loc, c.Obj, c.Grad, c.Step)
	s, ok := c.Fun.(Displayer)
	if ok {
		headings = AppendHeadings(headings, s)
	}
	return headings
}

func (c *Cubic) AppendValues(values []interface{}) []interface{} {
	values = AppendValues(values, c.Common)
	values = append(values, c.Grad.Curr(), c.Step.Curr)
	s, ok := c.Fun.(Displayer)
	if ok {
		values = AppendValues(values, s)
	}
	return values
}

func (c *Cubic) Result() {
	SetResults(c.Common, c.Loc, c.Obj, c.Grad, c.Step)
}

type CubicResult struct {
	StepHist []float64
}

func (cubic *Cubic) Iterate() (err error) {

	// Initialize
	var stepMultiplier float64
	updateCurrPoint := false
	reverseDirection := false
	currG := cubic.Grad.Curr()
	currF := cubic.Obj.Curr()

	// Evaluate trial point
	// Step Size is from the original point
	var trialX float64

	if cubic.initialGradNegative {
		trialX = cubic.Step.Curr() + cubic.Loc.Init()
	} else {
		trialX = -cubic.Step.Curr() + cubic.Loc.Init()
	}

	err = cubic.Fun.Eval(trialX)
	trialF := cubic.Fun.Obj()
	trialG := cubic.Fun.Grad()
	// Should this be embedded into Fun so every time eval is called
	// the count is updated?
	cubic.FunEvals.Add(1)
	cubic.Loc.AddToHist(trialX)
	cubic.Loc.AddToHist(trialF)
	cubic.Loc.AddToHist(trialG)

	/*
		fmt.Println("curr step size", cubic.Step.Curr)
		fmt.Println("LB", cubic.Step.Lb)
		fmt.Println("UB", cubic.Step.Ub)
		fmt.Println("initX", cubic.Loc.Init)
		fmt.Println("currX", cubic.Loc.Curr)
		fmt.Println("trialX", trialX)
		fmt.Println("InitF", cubic.Obj.Init())
		fmt.Println("currF", currF)
		fmt.Println("trialF", trialF)
		fmt.Println("InitG", cubic.Grad.Init)
		fmt.Println("currG", currG)
		fmt.Println("trialG", trialG)
	*/

	//cubic.AddToHist(trialX, trialF,trialG)
	absTrialG := math.Abs(trialG)

	// Find guess for next point
	deltaF := trialF - currF
	decreaseInValue := (deltaF < 0)
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
	if a == 0 {
		//Perfect quadratic fit
		stepMultiplier = -c / (2 * b)
	} else if det < 0 {
		if decreaseInValue && !changeInDerivSign {
			// The trial point has lower function value
			// and steeper gradient. Set this location as a
			// lower bound for the minimum, and set the
			// next point farther in that direction.

			cubic.SetBound("Lower")
			// We know we need to increase in step, but unsure how much, so make a guess
			stepMultiplier = cubic.UnclearStepIncrease()
			if decreaseInDerivMagnitude {
				updateCurrPoint = true
			}
		} else {
			// All other conditions we want to decrease the step size, but the
			// cubic doesn't give an estimate of how much. Just do a binary search
			cubic.SetBound("Upper")
			for i := 0; i < 10; i++ {
				fmt.Println("Blah")
			}
			fmt.Println(cubic.Step.Curr)
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
			cubic.SetBound("Upper")

			stepMultiplier = cubic.SizeDecrease(minCubic)
			if decreaseInDerivMagnitude && decreaseInValue {
				updateCurrPoint = true
				reverseDirection = true
			}

		case !changeInDerivSign && decreaseInValue:
			// No change in derivative sign, but a decrease in
			// function value. Want to move more in this direction
			// and the trial point is a new lower bound and a new
			// base for the cubic approximation

			cubic.SetBound("Lower")

			updateCurrPoint = true

			if decreaseInDerivMagnitude {
				if minCubic < cubic.StepIncreaseMin {
					// Cubic gave a bad approximation (minimum is more
					// in this direction). Assume linear decrease in
					// derivative

					// Check this line
					stepMultiplier = math.Abs(currG) / (math.Abs(trialG) - math.Abs(currG))
				}
				stepMultiplier = cubic.SizeIncrease(stepMultiplier)

			} else {
				// Found a better point, but the derivative increased
				// Use cubic approximation if it gives a reasonable guess
				// otherwise just project forward
				if minCubic < cubic.StepIncreaseMin {
					stepMultiplier = cubic.UnclearStepIncrease()
				} else {
					stepMultiplier = cubic.SizeIncrease(minCubic)
				}

			}
		case !changeInDerivSign && !decreaseInValue:
			// Neither a decrease in value nor a change in derivative sign.
			// This means there must be a local minimum between the starting location and
			// this one. Don't update the cubic point
			cubic.SetBound("Upper")
			stepMultiplier = math.Min(minCubic, 0.75)
			stepMultiplier = math.Max(stepMultiplier, 0.25)
		}

	}

	var deltaXTrialCurrent float64
	deltaXTrialCurrent = cubic.Step.Curr() - cubic.deltaCurrent
	newDeltaXFromCurrent := deltaXTrialCurrent * stepMultiplier

	var newStepSize float64
	newStepSize = newDeltaXFromCurrent + cubic.deltaCurrent

	// Want to make sure that the new search location isn't pushing beyond
	// previously established bounds. If it is, just do a binary search between
	// the bounds
	if !cubic.Step.WithinBounds(newStepSize) {
		newStepSize = cubic.Step.Midpoint()
	}

	if updateCurrPoint {

		cubic.Loc.SetCurr(trialX)
		cubic.Obj.SetCurr(trialF)
		cubic.Grad.SetCurr(trialG)
		cubic.deltaCurrent = trialX - cubic.Loc.Init()
		if cubic.initialGradNegative {
			cubic.deltaCurrent *= -1
		}
		if reverseDirection {
			cubic.currStepDirectionPositive = !cubic.currStepDirectionPositive
		}
	}
	cubic.Step.SetCurr(newStepSize)
	return nil
}

func (c *Cubic) UnclearStepIncrease() float64 {
	// Increase the step. If there is an upper bound, do a binary
	// search between the upper and lower bound. If there is no
	// upper bound, just double the step size.

	// Clean up this code!
	var stepMultiplier float64
	if c.currStepDirectionPositive {
		if math.IsInf(c.Step.Ub(), 1) {
			stepMultiplier = (2*c.Step.Curr() - c.deltaCurrent) / (c.Step.Curr() - c.deltaCurrent)
		} else {
			stepMultiplier = (c.Step.Midpoint() - c.deltaCurrent) / (c.Step.Curr() - c.deltaCurrent)
		}
	} else {
		if math.IsInf(c.Step.Lb(), -1) {
			stepMultiplier = (2*c.Step.Curr() - c.Loc.Curr()) / (c.Step.Curr() - c.Loc.Curr())
		} else {
			stepMultiplier = (c.Step.Midpoint() - c.Loc.Curr()) / (c.Step.Curr() - c.Loc.Curr())
		}
	}
	return stepMultiplier
}

func (c *Cubic) SetBound(dir string) {
	if dir == "Lower" {
		// Want to go farther in this direction
		if c.currStepDirectionPositive == true {
			c.Step.SetLb(c.Step.Curr())
		} else {
			c.Step.SetUb(c.Step.Curr())
		}
	} else if dir == "Upper" {
		if c.currStepDirectionPositive == true {
			c.Step.SetUb(c.Step.Curr())
		} else {
			c.Step.SetLb(c.Step.Curr())
		}
	}
	return
}

func (c *Cubic) SizeDecrease(minCubic float64) float64 {

	stepMultiplier := math.Max(minCubic, c.StepDecreaseMin)
	stepMultiplier = math.Min(stepMultiplier, c.StepDecreaseMax)
	return stepMultiplier
}

func (c *Cubic) SizeIncrease(minCubic float64) float64 {
	stepMultiplier := math.Max(minCubic, c.StepIncreaseMin)
	stepMultiplier = math.Min(stepMultiplier, c.StepIncreaseMax)
	return stepMultiplier
}
