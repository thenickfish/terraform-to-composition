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

var modulePath string

var rootCmd = &cobra.Command{
	Use:   "terraform-to-composition",
	Short: "A tool to convert Terraform schemas to Crossplane compositions",
	Long:  `A tool to convert Terraform schemas to Crossplane compositions`,
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a Crossplane XRD and XR from a Terraform schema",
	Long:  `Generate a Crossplane XRD and XR from a Terraform schema`,
	Run: func(cmd *cobra.Command, args []string) {
		generateSchema("test-xrd-name")
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringVarP(&modulePath, "module-path", "p", ".", "Path to the Terraform module to generate the XRD and XR from")
}

func Execute() error {
	return rootCmd.Execute()
}

func generateSchema(xrdName string) {

	// Create a new schema
	options := jsonschema.CreateSchemaOptions{}
	schema, _ := jsonschema.CreateSchema(modulePath, options)

	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)

	err := encoder.Encode(schema)
	if err != nil {
		fmt.Printf("Error encoding schema to JSON: %v\n", err)
	}

	propertiesMap, ok := schema["properties"].(map[string]interface{})
	if !ok {
		propertiesMap = map[string]interface{}{}
	}

	templateData := t.TemplateData{
		Name:         "xrdname",
		ModuleSource: "https://github.com/test",
		Group:        "thenickfish.github.com",
		Version:      "v1alpha1",
		Kind:         "XMyDatabase",
		PluralKind:   pluralize.NewClient().Plural("xmydatabase"),
		Properties:   propertiesMap,
	}

	text, _ := t.GenerateXrd(templateData)
	fmt.Println(text)

	text, _ = t.GenerateXr(templateData)
	fmt.Println(text)
}
