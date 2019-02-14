package gomod

import (
	"testing"
)

func TestListDeps(t *testing.T) {
	_, err := ListDeps("./testdata/mod")
	if err != nil {
		t.Fatal(err)
	}

	// TODO: finish test
}
