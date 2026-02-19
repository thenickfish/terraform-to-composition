package pkg

type TemplateData struct {
	Name         string
	ModuleSource string
	Group        string
	Scope        string
	Version      string
	Kind         string
	PluralKind   string
	Properties   map[string]interface{}
}
