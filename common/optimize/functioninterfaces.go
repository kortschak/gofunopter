package optimize

type UniObj interface {
	Objective(x float64) (obj float64, err error)
}

type UniGrad interface {
	Gradient(x float64) (obj float64, err error)
}

type UniObjGrad interface {
	ObjGrad(x float64) (obj float64, grad float64, err error)
}

type MultiObj interface {
	Objective(x []float64) (obj float64, err error)
}

type MultiGrad interface {
	Gradient(x []float64) (grad []float64, err error)
}

type MultiObjGrad interface {
	ObjGrad(x []float64) (obj float64, grad []float64, err error)
}
