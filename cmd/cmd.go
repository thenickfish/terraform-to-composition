package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/HewlettPackard/terraschema/pkg/jsonschema"
	pluralize "github.com/gertd/go-pluralize"
	t "github.com/nickfish/terraform-to-composition/pkg/template"
)

func Execute() error {

	//add cobra arg for xrd name
	xrdName := "example"
	cmd := &cobra.Command{
		Use:   "generate-schema",
		Short: "Generate a JSON schema for a Crossplane XRD",
		Run: func(cmd *cobra.Command, args []string) {
			generateSchema(xrdName)
		},
	}

	if err := cmd.Execute(); err != nil {
		return err
	}

	return nil
}

func generateSchema(xrdName string) {

	// Create a new schema
	path := "/Users/nick/Library/CloudStorage/SynologyDrive-development/crossplane-demo/terraform-module-example"
	options := jsonschema.CreateSchemaOptions{}
	schema, _ := jsonschema.CreateSchema(path, options)

	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)

	// encoder.SetEscapeHTML(escapeJSON)
	// encoder.SetIndent("", jsonIndent)

	err := encoder.Encode(schema)
	if err != nil {
		fmt.Printf("Error encoding schema to JSON: %v\n", err)
	}
	jsonOutput := buffer.Bytes()

	fmt.Println(string(jsonOutput))

	// generate yaml from the json schema
	propertiesfromjson := schema["properties"]
	//yaml
	// yamlOutput, err := yaml.Marshal(propertiesfromjson)

	fmt.Printf("Properties from JSON: %+v\n", propertiesfromjson)

	// assert schema["properties"] to the expected type for TemplateData
	propertiesMap, ok := schema["properties"].(map[string]interface{})
	if !ok {
		propertiesMap = map[string]interface{}{}
	}

	templateData := t.TemplateData{
		Name:         "xrdname",
		ModuleSource: "https://github.com/test",
		Group:        "thenickfish2.github.com",
		Version:      "v1alpha1",
		Kind:         "XMyDatabase",
		PluralKind:   pluralize.NewClient().Plural("xmydatabase"),
		// Properties:   string(yamlOutput),
		Properties: propertiesMap,
	}
	text, _ := t.GenerateXrd(templateData)
	fmt.Println(text)
	text, _ = t.GenerateXr(templateData)
	fmt.Println(text)
	// Print the schema
	// fmt.Printf("Schema: %+v\n", schema)

}
