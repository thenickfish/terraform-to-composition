package pkg

import (
	_ "embed"
	"os"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"gopkg.in/yaml.v3"
)

//go:embed xrd.yaml
var xrdTemplate string

func GenerateXrd(xrdData TemplateData) (string, error) {
	yamlProperties, err := yaml.Marshal(xrdData.Properties)
	if err != nil {
		return "", err
	}
	tmpl, err := template.New("example").Funcs(sprig.FuncMap()).Parse(xrdTemplate)
	if err != nil {
		return "", err
	}

	err = tmpl.Execute(os.Stdout, map[string]interface{}{
		"Name":           xrdData.Name,
		"Group":          xrdData.Group,
		"Version":        xrdData.Version,
		"Kind":           xrdData.Kind,
		"PluralKind":     xrdData.PluralKind,
		"YamlProperties": string(yamlProperties),
	})
	if err != nil {
		return "", err
	}

	return "", nil
}
