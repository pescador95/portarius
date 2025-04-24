package inventory

type IInventoryImporter interface {
	ImportPetsFromCSV() error
}
