package sourced

import (
	"github.com/go-enry/go-license-detector/v4/licensedb/filer"

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
