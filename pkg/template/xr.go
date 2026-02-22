package pkg

import (
	"bytes"
	_ "embed"
	"fmt"
	"log/slog"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"gopkg.in/yaml.v3"
)

type Patch struct {
	Type          string `yaml:"type"`
	FromFieldPath string `yaml:"fromFieldPath"`
	ToFieldPath   string `yaml:"toFieldPath"`
}

//go:embed xr.yaml
var xrTemplate string

func GenerateXr(xrData TemplateData) ([]byte, error) {
	slog.Debug("Generating Composition (XR)...")
	patches := []Patch{}
	for key := range xrData.Inputs {
		patches = append(patches, Patch{
			Type:          "FromCompositeFieldPath",
			FromFieldPath: fmt.Sprintf(key),
			ToFieldPath:   fmt.Sprintf("spec.forProvider.varmap.%s", key),
		})
	}
	for _, output := range xrData.Outputs {
		if output.Sensitive {
			continue
		}
		patches = append(patches, Patch{
			Type:          "ToCompositeFieldPath",
			FromFieldPath: fmt.Sprintf("status.atProvider.outputs.%s", output.Name),
			ToFieldPath:   fmt.Sprintf("status.%s", output.Name),
		})
	}
	slog.Debug("Patches generated", "patches", patches)
	patchesYaml, err := yaml.Marshal(patches)
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New("xr.yaml").Funcs(sprig.FuncMap()).Parse(xrTemplate)
	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, map[string]interface{}{
		"Name":         xrData.Name,
		"Group":        xrData.Group,
		"Version":      xrData.Version,
		"Kind":         xrData.Kind,
		"PluralKind":   xrData.PluralKind,
		"PatchesYaml":  string(patchesYaml),
		"ModuleSource": xrData.ModuleSource,
	})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
