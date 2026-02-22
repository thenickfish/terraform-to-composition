package hclparser

type Output struct {
	Name        string `hcl:"name,label"`
	Description string `hcl:"description,optional"`
	Sensitive   bool   `hcl:"sensitive,optional"`
}
