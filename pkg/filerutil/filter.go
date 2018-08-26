package filerutil

import (
	"gopkg.in/src-d/go-license-detector.v2/licensedb/filer"
)

type filter interface {
	// Filter decides whether a file should be included in the file list or not.
	// If the returned value is true, the file should included in the file list.
	Filter(file filer.File) bool
}

// FilterFunc makes a filter from a function when it's signature is compatible.
type FilterFunc func(file filer.File) bool

// Filter calls the underlying filter function.
func (fn FilterFunc) Filter(file filer.File) bool {
	return fn(file)
}

// Filters wraps a number of filters and executes all of them for each file.
type Filters []filter

// Filter calls all filters for a file.
func (f Filters) Filter(file filer.File) bool {
	for _, filter := range f {
		if !filter.Filter(file) {
			return false
		}
	}

	return true
}

type filterFiler struct {
	filer  filer.Filer
	filter filter
}

// NewFilterFiler returns a filer that wraps another filer
// and filters out common sources of wrong license detections.
func NewFilterFiler(filer filer.Filer, filter filter) *filterFiler {
	return &filterFiler{
		filer:  filer,
		filter: filter,
	}
}

// ReadFile returns the contents of a file given it's path.
func (f *filterFiler) ReadFile(path string) (content []byte, err error) {
	return f.filer.ReadFile(path)
}

// ReadDir lists a directory.
func (f *filterFiler) ReadDir(path string) ([]filer.File, error) {
	files, err := f.filer.ReadDir(path)
	if err != nil {
		return files, err
	}

	var filteredFiles []filer.File

	for _, file := range files {
		if f.filter.Filter(file) {
			filteredFiles = append(filteredFiles, file)
		}
	}

	return filteredFiles, nil
}

// Close frees all the resources allocated by this Filer.
func (f *filterFiler) Close() {
	f.filer.Close()
}
