package gofunopt

import (
	"fmt"
	"math"
)

// Maybe all of these should go in their own subpackages?
// SISOGradBased, etc. That way have SISO.Optimizable?

type Cubic struct {
	// Basic values
	Loc  *OptFloat // Location
	Obj  *OptFloat // Function Value
	Grad *OptFloat // Gradient value
	Step *OptFloat // Step size
	*Common

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
	fun                       SISOProblem
}

func DefaultCubic() *Cubic {
	c := &Cubic{
		Loc:             DefaultInputFloat(),
		Obj:             DefaultObjectiveFloat(),
		Grad:            DefaultGradientFloat(),
		Step:            DefaultStepFloat(),
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
func (c *Cubic) Initialize(fun SISOOptimizable) (err error) {
	// Initialize takes all of these in so function evaluations can be saved if 
	// the information is already there
	c.Loc.Initialize()
	if math.IsNaN(f.Curr()) {
		// Initial function value hasn't been set, so do it.
		err = fun.Eval(x.Curr())
		if err != nil {
			return fmt.Errorf("Error evaluating the function at the set initial value %v", x.Curr())
		}
		c.Obj.Init = fun.Obj()
		c.Grad.Init = fun.Grad()
	}
	c.Obj.Initialize()
	c.Grad.Initialize()
	if step.Init <= 0 {
		return fmt.Errorf("Initial step must be positive")
	}
	c.Step.Initialize()
	c.initialGradNegative = (c.Grad.Curr() < 0)
	c.currStepDirectionPositive = true
	c.deltaCurrent = 0.0 // How far is the current point from the initial point

	// Add in some checking on the Step Increase and decrease sizes
	return nil
}

func (c *Cubic) CheckConvergence() string {
	str := CheckConvergence(c.Loc, c.Obj, c.Grad, c.Step)
	if str != "" {
		return str
	}
	_, ok = c.fun.(Converger)
	if ok {
		str := CheckConvergence(c.fun)
	}
	return ""
}

func (c *Cubic) DisplayHeadings() []string {
	headings := make([]string, 10)
	headings = AppendHeadings(headings, c.Common)
	headings = append(headings, "Grad", "StepSize")
	_, ok = c.fun.(Displayer)
	if ok {
		headings = AppendHeadings(headings, fun)
	}
	return headings
}

func (c *Cubic) DisplayValues() []interface{} {
	values := make([]interface{}, 10)
	values = AppendValues(values, c.Common)
	values = append(values, c.g.Curr(), c.step.Curr())
	_, ok = c.fun.(Displayer)
	if ok {
		values = AppendValues(values, fun)
	}
	return a
}

func (c *Cubic) Optimize(fun SISOOptimizable) (r SISOResult, err error) {
	err = c.Initialize(fun)
	if err != nil {
		return nil, err
	}
	for {
		str := c.CheckConvergence()
		if str != "" {
			c.Display.Values()
			break
		}
		err = c.Iterate()
		if err != nil {
			break
		}
	}
	// Want to return the result even if there is an error in case anything
	// gets lost (maybe a defer would be even better?)
	return c.Result(), err
}

func (cubic *Cubic) Iterate() (err error){

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
    
    err  = cubic.fun.Eval(trialX)
    trialF = cubic.fun.Obj()
    trailG = cubic.Fun.Grad()
	// Should this be embedded into Fun so every time eval is called
	// the count is updated?
    cubic.NumFunEvals.Add(1)
    
    /*
    fmt.Println("curr step size",cubic.step.Size())
    fmt.Println("LB", cubic.step.Lb())
    fmt.Println("UB", cubic.step.Ub())
    fmt.Println("initX", cubic.x.Init())
    fmt.Println("currX", cubic.x.Curr())
    fmt.Println("trialX", trialX)
    fmt.Println("InitF", cubic.f.Init())
    fmt.Println("currF", currF)
    fmt.Println("trialF",trialF)
        fmt.Println("InitG", cubic.g.Init())
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
		if (decreaseInValue && !changeInDerivSign) {
            // The trial point has lower function value
            // and steeper gradient. Set this location as a
            // lower bound for the minimum, and set the
            // next point farther in that direction.
            
            cubic.SetBound("Lower")
            // We know we need to increase in step, but unsure how much, so make a guess
            stepMultiplier = cubic.UnclearStepIncrease()
            if decreaseInDerivMagnitude{
                updateCurrPoint = true
            }
		}else{
            // All other conditions we want to decrease the step size, but the
            // cubic doesn't give an estimate of how much. Just do a binary search
            cubic.SetBound("Upper")
            for i:=0;i<10;i++{
                fmt.Println("Blah")
            }
            fmt.Println(cubic.step.Size())
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
			if decreaseInDerivMagnitude && decreaseInValue{
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
				if minCubic < cubic.increaseStepLowerBound {
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
				if minCubic < cubic.increaseStepLowerBound {
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
	deltaXTrialCurrent = cubic.step.Size() - cubic.deltaCurrent
	newDeltaXFromCurrent := deltaXTrialCurrent * stepMultiplier

	var newStepSize float64
	newStepSize = newDeltaXFromCurrent + cubic.deltaCurrent

	// Want to make sure that the new search location isn't pushing beyond
	// previously established bounds. If it is, just do a binary search between
	// the bounds
	if !cubic.step.WithinBounds(newStepSize){
		newStepSize = cubic.step.BoundMidpoint()
	}

	if updateCurrPoint {
        
        cubic.x.Set(trialX)
        cubic.f.Set(trialF)
        cubic.g.Set(trialG)
		cubic.deltaCurrent = trialX - cubic.x.Init()
		if cubic.initialGradNegative {
			cubic.deltaCurrent *= -1
		}
		if reverseDirection {
			cubic.currStepDirectionPositive = !cubic.currStepDirectionPositive
		}
	}
	cubic.step.SetSize(newStepSize)
    
    
    //fmt.Println("\n")
}

func (c *Cubic) UnclearStepIncrease() float64 {
	// Increase the step. If there is an upper bound, do a binary
	// search between the upper and lower bound. If there is no
	// upper bound, just double the step size.

	// Clean up this code!
	var stepMultiplier float64
	if c.currStepDirectionPositive {
		if math.IsInf(c.step.Ub(), 1) {
			stepMultiplier = (2*c.step.Size() - c.deltaCurrent) / (c.step.Size() - c.deltaCurrent)
		} else {
			stepMultiplier = (c.step.BoundMidpoint() - c.deltaCurrent) / (c.step.Size() - c.deltaCurrent)
		}
	} else {
		if math.IsInf(c.step.Lb(), -1) {
			stepMultiplier = (2*c.step.Size() - c.x.Curr()) / (c.step.Size() - c.x.Curr())
		} else {
			stepMultiplier = (c.step.BoundMidpoint() - c.x.Curr()) / (c.step.Size() - c.x.Curr())
		}
	}
	return stepMultiplier
}

func (c *Cubic) SetBound(dir string) {
	if dir == "Lower" {
		// Want to go farther in this direction
		if c.currStepDirectionPositive == true {
			c.step.SetLb(c.step.Size())
		} else {
			c.step.SetUb(c.step.Size())
		}
	} else if dir == "Upper" {
		if c.currStepDirectionPositive == true {
			c.step.SetUb(c.step.Size())
		} else {
			c.step.SetLb(c.step.Size())
		}
	}
	return
}

func (c *Cubic) SizeDecrease(minCubic float64) float64 {

	stepMultiplier := math.Max(minCubic, c.decreaseStepLowerBound)
	stepMultiplier = math.Min(stepMultiplier, c.decreaseStepUpperBound)
	return stepMultiplier
}

func (c *Cubic) SizeIncrease(minCubic float64) float64 {
	stepMultiplier := math.Max(minCubic, c.increaseStepLowerBound)
	stepMultiplier = math.Min(stepMultiplier, c.increaseStepUpperBound)
	return stepMultiplier
}
}