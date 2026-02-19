package pkg

import (
	_ "embed"
	"fmt"
	"os"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"gopkg.in/yaml.v3"
)

type Patch struct {
	Type          string
	FromFieldPath string
	ToFieldPath   string
}

//go:embed xr.yaml
var xrTemplate string

func GenerateXr(xrData TemplateData) (string, error) {
	patches := []Patch{}
	for key, _ := range xrData.Properties {
		patch := Patch{
			Type:          "FromCompositeFieldPath",
			FromFieldPath: fmt.Sprintf("spec.%s", key),
			ToFieldPath:   fmt.Sprintf("spec.forProvider.varmap.%s", key),
		}
		patches = append(patches, patch)
	}
	fmt.Printf("Patches: %+v\n", patches)
	// xrData.Properties = fmt.Sprintf("%s", patches)
	patchesYaml, err := yaml.Marshal(patches)
	if err != nil {
		return "", err
	}
	fmt.Printf("Patches YAML:\n%s\n", string(patchesYaml))

	tmpl, err := template.New("example").Funcs(sprig.FuncMap()).Parse(xrTemplate)
	if err != nil {
		return "", err
	}

	err = tmpl.Execute(os.Stdout, map[string]interface{}{
		"Name":        xrData.Name,
		"Group":       xrData.Group,
		"Version":     xrData.Version,
		"Kind":        xrData.Kind,
		"PluralKind":  xrData.PluralKind,
		"PatchesYaml": string(patchesYaml),
	})
	if err != nil {
		return "", err
	}
	return "", nil
}
