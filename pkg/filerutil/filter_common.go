package filerutil

import (
	"fmt"
	"path"
	"regexp"
	"strings"

	"gopkg.in/src-d/go-license-detector.v2/licensedb/filer"
)

var licenseDocsRe = regexp.MustCompile(fmt.Sprintf("^(|.*[-_. ])(%s)(\\.docs?.*)$", "li[cs]en[cs]e(s?)"))

// CorrectnessFilter does some corrections based on known issues (like detecting multiple licenses).
func CorrectnessFilter(file filer.File) bool {
	return !licenseDocsRe.MatchString(strings.ToLower(path.Base(file.Name)))
}

// DirFilter excludes directories from the file list.
// This is usually a good decision since most projects store their license files
// in the project root.
func DirFilter(file filer.File) bool {
	return !file.IsDir
}
