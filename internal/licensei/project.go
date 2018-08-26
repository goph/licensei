package licensei

type Project struct {
	Name       string  `toml:"name"`
	Revision   string  `toml:"revision"`
	License    string  `toml:"license"`
	Confidence float32 `toml:"confidence"`
}

type projectSource interface {
	Projects() ([]Project, error)
}
