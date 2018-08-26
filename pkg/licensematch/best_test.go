package licensematch_test

import (
	"reflect"
	"testing"

	. "github.com/goph/licensei/pkg/licensematch"
)

func TestBest(t *testing.T) {
	tests := map[string]struct {
		licenses map[string]float32

		license    string
		confidence float32
	}{
		"obvious best match": {
			licenses: map[string]float32{
				"MIT":     0.98,
				"ECL-2.0": 0.81,
			},
			license:    "MIT",
			confidence: 0.98,
		},
		"first best match": {
			licenses: map[string]float32{
				"MIT":        0.98,
				"Apache-2.0": 0.98,
				"ECL-2.0":    0.81,
			},
			license:    "Apache-2.0",
			confidence: 0.98,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			license, confidence := Best(test.licenses)

			if license != test.license {
				t.Errorf("unexpected license, got: %s, want: %s", license, test.license)
			}

			if confidence != test.confidence {
				t.Errorf("unexpected confidence, got: %f, want: %f", confidence, test.confidence)
			}
		})
	}
}

func TestBestAll(t *testing.T) {
	tests := map[string]struct {
		licenses map[string]float32

		bests      map[string]float32
		confidence float32
	}{
		"one best match": {
			licenses: map[string]float32{
				"MIT":     0.98,
				"ECL-2.0": 0.81,
			},
			bests: map[string]float32{
				"MIT": 0.98,
			},
			confidence: 0.98,
		},
		"two best matches": {
			licenses: map[string]float32{
				"MIT":        0.98,
				"Apache-2.0": 0.98,
				"ECL-2.0":    0.81,
			},
			bests: map[string]float32{
				"MIT":        0.98,
				"Apache-2.0": 0.98,
			},
			confidence: 0.98,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			bests, confidence := Bests(test.licenses)

			if !reflect.DeepEqual(bests, test.bests) {
				t.Errorf("unexpected results, got: %v, want: %v", bests, test.bests)
			}

			if confidence != test.confidence {
				t.Errorf("unexpected confidence, got: %f, want: %f", confidence, test.confidence)
			}
		})
	}
}
