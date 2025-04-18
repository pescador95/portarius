package resident

type IResidentImporter interface {
	ImportResidentsFromCSV() error
}
