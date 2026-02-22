package pkg

import (
	"github.com/nickfish/terraform-to-composition/pkg/hclparser"
)

type TemplateData struct {
	Name         string
	ModuleSource string
	Group        string
	Scope        string
	Version      string
	Kind         string
	PluralKind   string
	Inputs       map[string]interface{}
	Outputs      []hclparser.Output
}
