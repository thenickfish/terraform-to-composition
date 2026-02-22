package cmd

import (
	"os"

	"github.com/spf13/cobra"

	pluralize "github.com/gertd/go-pluralize"
	"github.com/nickfish/terraform-to-composition/pkg/hclparser"
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

	inputs, outputs, err := hclparser.ParseModule(modulePath)
	if err != nil {
		panic(err)
	}

	templateData := t.TemplateData{
		Name:         "bucket",
		ModuleSource: "https://github.com/thenickfish/crossplane-demo",
		Group:        "thenickfish.github.com",
		Version:      "v1alpha1",
		Kind:         "XBucket",
		PluralKind:   pluralize.NewClient().Plural("xbucket"),
		Inputs:       inputs,
		Outputs:      outputs,
	}

	data, err := t.GenerateXrd(templateData)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("out/xrd.yaml", data, 0644)
	if err != nil {
		panic(err)
	}

	data, err = t.GenerateXr(templateData)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("out/xr.yaml", data, 0644)
	if err != nil {
		panic(err)
	}

}
