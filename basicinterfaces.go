package gofunopter

type Initializer interface {
	Initialize() err
}

func Initialize(initializers ...Initializer) (err error) {
	for _, initializer := range initializers {
		err := initializer.Initialize()
		if err != nil {
			return err
		}
	}
	return nil
}

type Iterator interface {
	Iterate() error
}

func Iterate(iterators ...Iterator) (err error) {
	for _, iterator := range iterators {
		err := iterator.Iterate()
		if err != nil {
			return err
		}
	}
	return nil
}

type Converger interface {
	Converged() string
}

func Converged(convergers ...Converger) (str string) {
	for _, converger := range convergers {
		str = converger.Converged()
		if str != "" {
			return str
		}
	}
	return ""
}

type Resulter interface {
	Result()
}

func SetResults(resulters ...Resulter) {
	for _, resulter := range resulters {
		resulter.Result()
	}
}
