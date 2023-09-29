package importexport

type Importable interface {
	Import(string) any
}

type Exportable interface {
	Export(string) error
}
