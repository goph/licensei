package licensei

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type HeaderChecker struct {
	IgnorePaths []string
	IgnoreFiles []string
}

type HeaderViolations map[string]string

// nolint: gocognit
func (c HeaderChecker) Check(root string, template string) (HeaderViolations, error) {
	violations := HeaderViolations{}

	template = regexp.QuoteMeta(template)
	// `([0-9]{4}[,]?[\s]?(\band \b)?)+` allows to match multiple appearances of `:YEAR:` that a header might have
	// for example, "2019 and 2020", or "2019, 2020, and 2021"
	template = strings.Replace(template, ":YEAR:", "([0-9]{4}[,]?[\\s]?(\\band \\b)?)+", -1)

	err := filepath.Walk(root, func(path string, info os.FileInfo, _ error) error {
		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".go" {
			return nil
		}

		for _, ignorePath := range c.IgnorePaths {
			if !filepath.IsAbs(ignorePath) {
				ignorePath = filepath.Join(root, ignorePath)
			}

			if strings.HasPrefix(path, ignorePath+"/") {
				return nil
			}
		}

		for _, glob := range c.IgnoreFiles {
			if matched, err := filepath.Match(glob, filepath.Base(path)); err == nil && matched {
				return nil
			}
		}

		file, err := os.Open(path)
		if err != nil {
			violations[path] = err.Error()

			return nil
		}
		defer file.Close()

		b := bufio.NewReader(file)
		var lines string

		for len(lines) <= len(template) {
			line, err := b.ReadString('\n')
			if err == io.EOF {
				break
			}
			if err != nil {
				violations[path] = err.Error()

				return nil
			}

			if strings.HasPrefix(line, "// +build") {
				continue
			}

			if strings.HasPrefix(line, "// nolint") {
				continue
			}

			if line == "" {
				continue
			}

			lines += line
		}

		matched, err := regexp.MatchString(template, lines)
		if err != nil {
			violations[path] = err.Error()

			return nil
		}

		if !matched {
			violations[path] = "incorrect license header"

			return nil
		}

		return nil
	})

	return violations, err
}
