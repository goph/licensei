package licensei

import (
	"bufio"
	"bytes"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

type gomodDependencySource struct {
}

func NewGoModDependencySource() *gomodDependencySource {
	return new(gomodDependencySource)
}

func (s *gomodDependencySource) Dependencies() ([]Dependency, error) {
	modlistCmd := exec.Command("go", "list", "-deps", "-f", "{{with .Module}}{{if not .Main}}{{.Path}} {{.Version}}{{end}}{{end}}", "./...")

	var buf bytes.Buffer

	modlistCmd.Stdout = &buf

	err := modlistCmd.Run()
	if err != nil {
		return nil, errors.Wrap(err, "failed to list modules: "+buf.String())
	}

	var dependencies []Dependency
	depMap := map[string]bool{}

	scanner := bufio.NewScanner(&buf)
	for scanner.Scan() {
		rawMod := scanner.Text()
		if rawMod == "" {
			continue
		}

		mod := strings.Split(rawMod, " ")

		if _, ok := depMap[mod[0]]; ok {
			continue
		}

		depMap[mod[0]] = true

		dependencies = append(
			dependencies,
			Dependency{
				Name:     mod[0],
				Revision: mod[1],
				Type:     "gomod",
			},
		)
	}

	return dependencies, nil
}
