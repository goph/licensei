package licensei

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/olekukonko/tablewriter"
)

// ListViewModel holds information for a list view.
type ListViewModel struct {
	Projects []ListProjectItem `json:"projects"`
}

// ListProjectItem represents an item in the list view.
type ListProjectItem struct {
	Name       string  `json:"name"`
	License    string  `json:"license"`
	Confidence float32 `json:"confidence"`
}

type jsonListView struct {
	output io.Writer
}

// NewJsonListView returns a view that outputs a license list as JSON.
func NewJsonListView(output io.Writer) *jsonListView {
	return &jsonListView{
		output: output,
	}
}

// Render renders a license list as JSON.
func (v *jsonListView) Render(model ListViewModel) error {
	encoder := json.NewEncoder(v.output)

	return encoder.Encode(model)
}

type tableListView struct {
	output io.Writer
}

// NewTableListView returns a view that outputs a license list as JSON.
func NewTableListView(output io.Writer) *tableListView {
	return &tableListView{
		output: output,
	}
}

// Render renders a license list as JSON.
func (v *tableListView) Render(model ListViewModel) error {
	table := tablewriter.NewWriter(v.output)
	table.SetHeader([]string{"Package", "License", "Confidence"})

	for _, project := range model.Projects {
		if project.License == "" {
			table.Append([]string{project.Name, "no license file was found", ""})

			continue
		}

		table.Append([]string{project.Name, project.License, fmt.Sprintf("%d%%", int(project.Confidence*100))})
	}

	table.Render()

	return nil
}
