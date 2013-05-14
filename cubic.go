package gofunopter

import (
	"fmt"
	"math"
)

func DefaultCubic() *Cubic {
	c := &Cubic{
		loc:             DefaultLocationFloat(),
		obj:             DefaultObjectiveFloat(),
		grad:            DefaultGradientFloat(),
		step:            DefaultBoundedStepFloat(),
		StepDecreaseMin: 1E-4,
		StepDecreaseMax: 0.9,
		StepIncreaseMin: 1.25,
		StepIncreaseMax: 1E3,
		Common:          DefaultCommon(),
	}
	SetDisplayMethods(c)
	return c
}

// Maybe all of these should go in their own subpackages?
// SISOGradBased, etc. That way have SISO.Optimizable?

type Cubic struct {
	// Basic values (should these be interfaces?)
	loc  OptFloat        // Location
	obj  OptTolFloat     // Function Value
	grad OptTolFloat     // Gradient value
	step BoundedOptFloat // Step size
	*Common
	fun SisoGradBasedProblem

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

func (c *Cubic) Loc() OptFloat {
	return c.loc
}

func (c *Cubic) Obj() OptTolFloat {
	return c.obj
}

func (c *Cubic) Grad() OptTolFloat {
	return c.grad
}

func (c *Cubic) Fun() SisoGradBasedProblem {
	return c.fun
}

func (c *Cubic) SetFun(fun SisoGradBasedProblem) {
	c.fun = fun
}

// Should we add error checking to the evaluations?
func (c *Cubic) Initialize() (err error) {
	// Initialize takes all of these in so function evaluations can be saved if 
	// the information is already there
	c.Common.Initialize()
	c.loc.Initialize()
	if math.IsNaN(c.obj.Init()) || math.IsNaN(c.grad.Curr()) {
		// Initial function value hasn't been set, so do it.
		f, g, err := c.fun.Eval(c.loc.Init())
		if err != nil {
			return fmt.Errorf("Error evaluating the function at the set initial value %v", c.loc.Curr)
		}
		c.FunEvals().Add(1)
		c.obj.SetInit(f)
		c.grad.SetInit(g)
	}
	c.obj.Initialize()
	c.grad.Initialize()
	if c.step.Init() <= 0 {
		return fmt.Errorf("Initial step must be positive")
	}
	c.step.SetLb(0.0)
	c.step.SetUb(math.Inf(1))
	c.step.Initialize()
	c.initialGradNegative = (c.grad.Curr() < 0)
	c.currStepDirectionPositive = true
	c.deltaCurrent = 0.0 // How far is the current point from the initial point
	// Add in some checking on the Step Increase and decrease sizes
	return nil
}

func (c *Cubic) Converged() Convergence {
	conv := Converged(c.obj, c.grad, c.step)
	if conv != nil {
		return conv
	}
	s, ok := c.fun.(Converger)
	if ok {
		conv = s.Converged()
		if conv != nil {
			return conv
		}
	}
	return nil
}

func (c *Cubic) AppendHeadings(headings []string) []string {
	headings = AppendHeadings(headings, c.Common, c.loc, c.obj, c.grad, c.step)
	s, ok := c.fun.(Displayer)
	if ok {
		headings = AppendHeadings(headings, s)
	}
	return headings
}

func (c *Cubic) AppendValues(values []interface{}) []interface{} {
	values = AppendValues(values, c.Common)
	values = append(values, c.grad.Curr(), c.step.Curr)
	s, ok := c.fun.(Displayer)
	if ok {
		values = AppendValues(values, s)
	}
	return values
}

func (c *Cubic) SetResult() {
	SetResults(c.Common, c.loc, c.obj, c.grad, c.step)
}

type CubicResult struct {
	StepHist []float64
}

func (cubic *Cubic) Iterate() (err error) {
	// Initialize
	var stepMultiplier float64
	updateCurrPoint := false
	reverseDirection := false
	currG := cubic.grad.Curr()
	currF := cubic.obj.Curr()

	// Evaluate trial point
	// Step Size is from the original point
	var trialX float64

	if cubic.initialGradNegative {
		trialX = cubic.step.Curr() + cubic.loc.Init()
	} else {
		trialX = -cubic.step.Curr() + cubic.loc.Init()
	}

	trialF, trialG, err := cubic.fun.Eval(trialX)
	if err != nil {
		return &FunctionError{Err: err, Loc: trialX}
	}
	// Should this be embedded into Fun so every time eval is called
	// the count is updated?
	cubic.FunEvals().Add(1)
	cubic.loc.AddToHist(trialX)
	cubic.loc.AddToHist(trialF)
	cubic.loc.AddToHist(trialG)

	/*
		fmt.Println("curr step size", cubic.step.Curr())
		fmt.Println("LB", cubic.step.Lb())
		fmt.Println("UB", cubic.step.Ub())
		fmt.Println("initX", cubic.loc.Init())
		fmt.Println("currX", cubic.loc.Curr())
		fmt.Println("trialX", trialX)
		fmt.Println("InitF", cubic.obj.Init())
		fmt.Println("currF", currF)
		fmt.Println("trialF", trialF)
		fmt.Println("InitG", cubic.grad.Init())
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
			fmt.Println(cubic.step.Curr)
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
	deltaXTrialCurrent = cubic.step.Curr() - cubic.deltaCurrent
	newDeltaXFromCurrent := deltaXTrialCurrent * stepMultiplier

	var newStepSize float64
	newStepSize = newDeltaXFromCurrent + cubic.deltaCurrent

	// Want to make sure that the new search location isn't pushing beyond
	// previously established bounds. If it is, just do a binary search between
	// the bounds
	if !cubic.step.WithinBounds(newStepSize) {
		newStepSize = cubic.step.Midpoint()
	}

	if updateCurrPoint {

		cubic.loc.SetCurr(trialX)
		cubic.obj.SetCurr(trialF)
		cubic.grad.SetCurr(trialG)
		cubic.deltaCurrent = trialX - cubic.loc.Init()
		if cubic.initialGradNegative {
			cubic.deltaCurrent *= -1
		}
		if reverseDirection {
			cubic.currStepDirectionPositive = !cubic.currStepDirectionPositive
		}
	}
	//panic("wa")
	cubic.step.SetCurr(newStepSize)
	return nil
}

func (c *Cubic) UnclearStepIncrease() float64 {
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
			stepMultiplier = (2*c.step.Curr() - c.loc.Curr()) / (c.step.Curr() - c.loc.Curr())
		} else {
			stepMultiplier = (c.step.Midpoint() - c.loc.Curr()) / (c.step.Curr() - c.loc.Curr())
		}
	}
	return stepMultiplier
}

func (c *Cubic) SetBound(dir string) {
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
