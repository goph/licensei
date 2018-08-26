package licensematch_test

import (
	"reflect"
	"testing"

	. "github.com/goph/licensei/pkg/licensematch"
)

func TestMerge(t *testing.T) {
	licenses := map[string]float32{
		"MIT":     0.98,
		"ECL-2.0": 0.81,
	}

	newlicenses := map[string]float32{
		"MIT": 1.0,
	}

	want := map[string]float32{
		"MIT":     (0.98 + 1.0) / 2,
		"ECL-2.0": 0.81 / 2,
	}

	got := Merge(licenses, newlicenses)

	if !reflect.DeepEqual(got, want) {
		t.Error("unexpected license match merge")
	}
}
