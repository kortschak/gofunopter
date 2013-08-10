package linesearch

import (
	"errors"
	"github.com/gonum/floats"

	//"gofunopter/common"
	"github.com/btracey/gofunopter/common/optimize"
	"github.com/btracey/gofunopter/common/status"
	"github.com/btracey/gofunopter/common/uni"
	"github.com/btracey/gofunopter/univariate"

	"fmt"
)

// Result is a struct for returning the result from a linesearch
type LinesearchResult struct {
	Loc       []float64
	Obj       float64
	Grad      []float64
	Step      float64
	NFunEvals int
}

// LinesearchFun is a type which is the one-dimensional funciton that projects
// the multidimensional function onto the line
type linesearchFun struct {
	fun         optimize.MultiObjGrad
	wolfe       WolfeConditioner
	direction   []float64 // unit vector
	initLoc     []float64
	currLoc     []float64
	currLocCopy []float64 // In case the user-defined function changes the value
	currGrad    []float64
}

func (l *linesearchFun) ObjGrad(step float64) (f float64, g float64, err error) {
	// Take the step (need to add back in the scaling)
	for i, val := range l.direction {
		l.currLoc[i] = val*step + l.initLoc[i]
	}
	// Copy the location (in case the user-defined function modifies it)
	copy(l.currLocCopy, l.currLoc)
	f, gVec, err := l.fun.ObjGrad(l.currLocCopy)
	if err != nil {
		return f, g, errors.New("linesearch: error during user defined function")
	}
	// Add the function to the history so that it isn't thrown out
	// Copy the gradient vector (in case Fun modifies it)
	n := copy(l.currGrad, gVec)
	if n != len(l.currLocCopy) {
		return f, g, errors.New("linesearch: user defined function returned incorrect gradient length")
	}

	// Find the gradient in the direction of the search vector
	g = floats.Dot(l.direction, l.currGrad)
	l.wolfe.SetCurrState(f, g, step)
	return f, g, nil
}

func (l *linesearchFun) Status() status.Status {
	// Set the function and gradient values for the line searcher
	return l.wolfe.Status()
}

type LinesearchMethod interface {
	univariate.UniGradOptimizer
	Step() *uni.BoundedStep
}

// Linesearch performs a linesearch. Optimizer should turn off all non-wolfe status patterns for the gradient and step
func Linesearch(multifun optimize.MultiObjGrad, method LinesearchMethod, settings *univariate.UniGradSettings, wolfe WolfeConditioner, searchVector []float64, initLoc []float64, initObj float64, initGrad []float64) (*LinesearchResult, error) {

	// Linesearch modifies the values of the slices, but should revert the changes by the end

	// Find the norm of the search direction
	normSearchVector := floats.Norm(searchVector, 2)

	// Find the search direction (replace this with an input to avoid make?)
	direction := make([]float64, len(searchVector))
	copy(direction, searchVector)
	floats.Scale(1/normSearchVector, direction)

	// Find the initial projection of the gradient into the search direction
	initDirectionalGrad := floats.Dot(direction, initGrad)

	if initDirectionalGrad > 0 {
		return &LinesearchResult{}, errors.New("initial directional gradient must be negative")
	}

	// Set wolfe constants
	wolfe.SetInitState(initObj, initDirectionalGrad)
	wolfe.SetCurrState(initObj, initDirectionalGrad, 1.0)
	fun := &linesearchFun{
		fun:         multifun,
		wolfe:       wolfe,
		direction:   direction,
		initLoc:     initLoc,
		currLoc:     make([]float64, len(initLoc)),
		currLocCopy: make([]float64, len(initLoc)),
		currGrad:    make([]float64, len(initLoc)),
	}

	settings.InitialGradient = initDirectionalGrad
	settings.InitialObjective = initObj
	method.Step().SetInit(normSearchVector)

	// Run optimization, initial location is zero
	optVal, optLoc, result, err := univariate.OptimizeGrad(fun, 0, settings, method)
	//status, err := optimize.OptimizeOpter(method, fun)

	// Regerate results structure (do this before returning error in case optimizer can recover from it)
	// need to scale alpha_k because linesearch is x_k + alpha_k p_k
	r := &LinesearchResult{
		Loc:  fun.currLoc,
		Obj:  optVal,
		Grad: fun.currGrad,
		Step: optLoc / normSearchVector,
	}

	if err != nil {
		fmt.Println("Error in linsearch")
		return r, errors.New("linesearch: error during linesearch optimization: " + err.Error())
	}
	stat := result.Status
	// Check to make sure that the status due to wolfe status
	if stat != status.WolfeConditionsMet {
		// If the status wasn't because of wolfe conditions, see if they are met anyway
		c := wolfe.Status()
		if c == status.WolfeConditionsMet {
			// Conditions met, no problem
			return r, nil
		}
		// Conditions not met
		return r, errors.New("linesearch: status not because of wolfe conditions.")
	}
	return r, nil
}
