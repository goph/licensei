package main

import (
	"fmt"

	"github.com/goph/lisensei/pkg/detector/srcd"
	"github.com/goph/lisensei/pkg/licensed"
	"gopkg.in/src-d/go-license-detector.v2/licensedb/filer"
)

func main() {
	var detector interface {
		Detect() (map[string]float32, error)
	}

	packages, err := licensed.DepSource()
	if err != nil {
		panic(err)
	}

	for _, pkg := range packages {
		f, err := filer.FromDirectory("vendor/" + pkg.Name)
		if err != nil {
			panic(err)
		}
		detector = srcd.NewDetector(f)

		matches, err := detector.Detect()
		if err != nil {
			panic(err)
		}

		var license string
		var confidence float32

		for l, c := range matches {
			if c > confidence {
				license = l
				confidence = c
			}
		}

		fmt.Printf("%s: %s (%d%%)\n", pkg.Name, license, int32(confidence*100))
	}
}
