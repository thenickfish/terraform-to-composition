package pkg

import (
	"bytes"
	_ "embed"
	"log/slog"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/nickfish/terraform-to-composition/pkg/hclparser"
	"gopkg.in/yaml.v3"
)

//go:embed xrd.yaml
var xrdTemplate string

func GenerateXrd(xrdData TemplateData) ([]byte, error) {
	slog.Debug("Generating Composite Resource Definition (XRD)...")
	yamlProperties, err := yaml.Marshal(xrdData.Inputs)
	if err != nil {
		return nil, err
	}
	tmpl, err := template.New("xrd.yaml").Funcs(sprig.FuncMap()).Parse(xrdTemplate)
	if err != nil {
		return nil, err
	}

	outputs, err := getStatusProperties(xrdData.Outputs)
	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, map[string]interface{}{
		"Name":             xrdData.Name,
		"Group":            xrdData.Group,
		"Version":          xrdData.Version,
		"Kind":             xrdData.Kind,
		"PluralKind":       xrdData.PluralKind,
		"YamlProperties":   string(yamlProperties),
		"StatusProperties": outputs,
	})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func getStatusProperties(outputs []hclparser.Output) (string, error) {
	statusProps := make(map[string]interface{})
	for _, output := range outputs {
		if output.Sensitive {
			continue
		}

		statusProps[output.Name] = &statusProperty{
			Type: "string",
		}
	}
	yaml, err := yaml.Marshal(statusProps)
	if err != nil {
		return "", err
	}
	return string(yaml), nil
}

type statusProperty struct {
	Type string
}
