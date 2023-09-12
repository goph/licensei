package licensei_test

import (
	"bytes"
	"os"
	"testing"

	. "github.com/goph/licensei/internal/licensei"
)

func TestJsonListView_Render(t *testing.T) {
	buf := new(bytes.Buffer)

	view := NewJSONListView(buf)

	model := ListViewModel{
		Dependencies: []ListDependencyItem{
			{
				Name:       "test",
				License:    "MIT",
				Confidence: 1.0,
			},
		},
	}

	result, err := os.ReadFile("testdata/list/golden0.json")
	if err != nil {
		t.Fatal(err)
	}

	err = view.Render(model)
	if err != nil {
		t.Fatal(err)
	}

	if buf.String() != string(result) {
		t.Errorf("unexpected result: %s", buf.String())
	}
}

func TestTableListView_Render(t *testing.T) {
	buf := new(bytes.Buffer)

	view := NewTableListView(buf)

	model := ListViewModel{
		Dependencies: []ListDependencyItem{
			{
				Name:       "test",
				License:    "MIT",
				Confidence: 1.0,
			},
		},
	}

	result, err := os.ReadFile("testdata/list/golden0.table")
	if err != nil {
		t.Fatal(err)
	}

	err = view.Render(model)
	if err != nil {
		t.Fatal(err)
	}

	if buf.String() != string(result) {
		t.Errorf("unexpected result: %s", buf.String())
	}
}
