package gofunopter

type Linesearcher interface {
	Optimizer
	WolfeFunConst() float64
	WolfeGradConst() float64
}

type CubicLinesearch struct {
	*Cubic
	WolfeFunConst
}
