package licensei_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	. "github.com/goph/licensei/internal/licensei"
)

func TestJsonListView_Render(t *testing.T) {
	buf := new(bytes.Buffer)

	view := NewJsonListView(buf)

	model := ListViewModel{
		Projects: []ListProjectItem{
			{
				Name:       "test",
				License:    "MIT",
				Confidence: 1.0,
			},
		},
	}

	result, err := ioutil.ReadFile("testdata/list/golden0.json")
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
		Projects: []ListProjectItem{
			{
				Name:       "test",
				License:    "MIT",
				Confidence: 1.0,
			},
		},
	}

	result, err := ioutil.ReadFile("testdata/list/golden0.table")
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
