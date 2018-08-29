package licensei

type Dependency struct {
	Name       string  `toml:"name"`
	Revision   string  `toml:"revision"`
	License    string  `toml:"license"`
	Confidence float32 `toml:"confidence"`
	Type       string  `toml:"type"`
}

type dependencySource interface {
	Dependencies() ([]Dependency, error)
}

func indexDependencies(dependencies []Dependency) map[string]int {
	index := make(map[string]int, len(dependencies))

	for key, dep := range dependencies {
		index[dep.Name] = key
	}

	return index
}
