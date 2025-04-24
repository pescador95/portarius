package interfaces

type ICSVReservationImporter interface {
	ImportReservationsFromCSV() error
}
