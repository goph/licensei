package sourced

import (
	"gopkg.in/src-d/go-license-detector.v2/licensedb/filer"

	"github.com/goph/licensei/pkg/filerutil"
)

func FilerFromDirectory(path string) (filer.Filer, error) {
	filer, err := filer.FromDirectory(path)
	if err != nil {
		return filer, err
	}

	filters := filerutil.Filters{
		filerutil.FilterFunc(filerutil.CorrectnessFilter),
		filerutil.FilterFunc(filerutil.DirFilter),
	}

	return filerutil.NewFilterFiler(filer, filters), nil
}
