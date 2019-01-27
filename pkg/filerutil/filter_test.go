package filerutil_test

import (
	"errors"
	"reflect"
	"testing"

	. "github.com/goph/licensei/pkg/filerutil"
	"gopkg.in/src-d/go-license-detector.v2/licensedb/filer"
)

type filerStub struct {
	t *testing.T

	files map[string][]filer.File
}

func (f *filerStub) ReadFile(path string) (content []byte, err error) {
	f.t.Fatal("this method should not be called: ReadFile")

	return nil, nil
}

func (f *filerStub) ReadDir(path string) ([]filer.File, error) {
	files, ok := f.files[path]
	if !ok {
		return nil, errors.New("file structure not found")
	}

	return files, nil
}

func (f *filerStub) Close() {
	f.t.Fatal("this method should not be called: Close")
}

func TestFilterFiler(t *testing.T) {
	type filter interface {
		Filter(file filer.File) bool
	}

	tests := map[string]struct {
		filer  filer.Filer
		files  []filer.File
		filter filter
	}{
		"correctness": {
			filer: &filerStub{
				t: t,
				files: map[string][]filer.File{
					"": {
						filer.File{
							Name:  "license.docs",
							IsDir: false,
						},
						filer.File{
							Name:  "license",
							IsDir: false,
						},
					},
				},
			},
			files: []filer.File{
				{
					Name:  "license",
					IsDir: false,
				},
			},
			filter: FilterFunc(CorrectnessFilter),
		},
		"directory": {
			filer: &filerStub{
				t: t,
				files: map[string][]filer.File{
					"": {
						filer.File{
							Name:  "cmd",
							IsDir: true,
						},
						filer.File{
							Name:  "license",
							IsDir: false,
						},
					},
				},
			},
			files: []filer.File{
				{
					Name:  "license",
					IsDir: false,
				},
			},
			filter: FilterFunc(DirFilter),
		},
		"filters": {
			filer: &filerStub{
				t: t,
				files: map[string][]filer.File{
					"": {
						filer.File{
							Name:  "cmd",
							IsDir: true,
						},
						filer.File{
							Name:  "license.docs",
							IsDir: false,
						},
						filer.File{
							Name:  "license",
							IsDir: false,
						},
					},
				},
			},
			files: []filer.File{
				{
					Name:  "license",
					IsDir: false,
				},
			},
			filter: Filters{FilterFunc(CorrectnessFilter), FilterFunc(DirFilter)},
		},
	}

	for name, test := range tests {
		name, test := name, test

		t.Run(name, func(t *testing.T) {
			f := NewFilterFiler(test.filer, test.filter)

			files, err := f.ReadDir("")
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(files, test.files) {
				t.Errorf("expected the returned file list to be filtered")
			}
		})
	}
}
