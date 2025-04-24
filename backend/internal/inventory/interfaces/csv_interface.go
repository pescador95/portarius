package interfaces

type ICSVInventoryImporter interface {
	ImportPetsFromCSV() error
}
