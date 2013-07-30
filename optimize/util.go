package optimize

type Initializer interface {
	Initialize() error
}

func Initialize(i ...Initializer) error {
	for _, val := range i {
		err := val.Initialize()
		if err != nil {
			return err
		}
	}
	return nil
}

type SetResulter interface {
	SetResult()
}

func SetResult(resulters ...SetResulter) {
	for _, val := range resulters {
		val.SetResult()
	}
}
