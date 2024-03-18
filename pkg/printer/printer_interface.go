package printer

type Printer interface {
	PrintResult(ruleList map[int]string) error
}

func NewPrinter(printFormat string) Printer {
	switch printFormat {
	case "table":
		return NewTable()
	case "json":
		return NewJson()
	case "yaml":
		return NewYaml()
	default:
		return NewTable()
	}
}
