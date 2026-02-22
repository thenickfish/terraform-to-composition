package hclparser

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"os"
	"path"
	"path/filepath"

	"github.com/HewlettPackard/terraschema/pkg/jsonschema"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

var fileSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type:       "output",
			LabelNames: []string{"name"},
		},
	},
}

func ParseModule(modulePath string) (map[string]interface{}, []Output, error) {
	// check that directory exists
	if _, err := os.Stat(modulePath); err != nil {
		return nil, nil, err
	}

	inputs, err := getInputs(modulePath)
	if err != nil {
		return nil, nil, err
	}

	outputs, err := getOutputs(modulePath)
	if err != nil {
		return nil, nil, err
	}
	return inputs, outputs, err
}

func getInputs(modulePath string) (map[string]interface{}, error) {
	options := jsonschema.CreateSchemaOptions{}
	schema, err := jsonschema.CreateSchema(modulePath, options)
	if err != nil {
		return nil, err
	}

	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)

	err = encoder.Encode(schema)
	if err != nil {
		slog.Error("Error encoding schema to JSON.", "error", err)
		return nil, err
	}

	propertiesMap, ok := schema["properties"].(map[string]interface{})
	if !ok {
		propertiesMap = map[string]interface{}{}
	}
	slog.Info("props", "map", propertiesMap)

	return propertiesMap, nil
}

func getOutputs(modulePath string) ([]Output, error) {
	outputs := []Output{}
	files, err := filepath.Glob(path.Join(modulePath, "*.tf"))
	if err != nil {
		return nil, err
	}

	parser := hclparse.NewParser()
	for _, filename := range files {

		slog.Info("parsing file...", "filename", path.Base(filename))
		file, diag := parser.ParseHCLFile(filename)
		if diag.HasErrors() {
			return nil, diag
		}

		blocks, _, diag := file.Body.PartialContent(fileSchema)
		if diag.HasErrors() {
			return nil, diag
		}

		for _, block := range blocks.Blocks {
			name := block.Labels[0]
			attributes, diag := block.Body.JustAttributes()
			if diag.HasErrors() {
				return nil, diag
			}

			var val Output
			val.Name = name
			for name, attr := range attributes {
				if name == "value" {
					continue
				}

				value, diag := attr.Expr.Value(nil)
				if diag.HasErrors() {
					slog.Error("Error encountered while trying to calculate output value.", "diagnostics", diag, "block_name", name)
					return nil, diag
				}
				if name == "description" {
					err = gocty.FromCtyValue(value, &val.Description)
					if err != nil {
						return nil, err
					}
				}
				if name == "sensitive" && value.Type() == cty.Bool {
					gocty.FromCtyValue(value, &val.Sensitive)
				}
			}
			outputs = append(outputs, val)
		}
	}
	return outputs, nil
}
