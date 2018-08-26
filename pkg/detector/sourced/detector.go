package sourced

import (
	"context"

	"gopkg.in/src-d/go-license-detector.v2/licensedb"
	"gopkg.in/src-d/go-license-detector.v2/licensedb/filer"
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
	return licensedb.Detect(d.filer)
}

func (d *detector) DetectContext(ctx context.Context) (map[string]float32, error) {
	return d.Detect()
}
