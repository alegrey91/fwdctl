package template

type Template interface {
	GenerateTemplate(output string) error
}