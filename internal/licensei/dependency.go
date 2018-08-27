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
