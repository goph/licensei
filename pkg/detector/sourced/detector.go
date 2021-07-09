package sourced

import (
	"context"

	"github.com/go-enry/go-license-detector/v4/licensedb"
	"github.com/go-enry/go-license-detector/v4/licensedb/filer"
)

type detector struct {
	filer filer.Filer
}

func NewDetector(filer filer.Filer) *detector {
	d := &detector{
		filer: filer,
	}

	return d
}

func (d *detector) Detect() (map[string]float32, error) {
	matches, err := licensedb.Detect(d.filer)
	if err != nil {
		return nil, err
	}

	m := make(map[string]float32, len(matches))

	for l, v := range matches {
		m[l] = v.Confidence
	}

	return m, nil
}

func (d *detector) DetectContext(ctx context.Context) (map[string]float32, error) {
	return d.Detect()
}
