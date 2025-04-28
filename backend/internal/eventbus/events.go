package eventbus

type PackageCreatedEvent struct {
	PackageID uint
	Channel   string
	Recipient string
}

type ReservationCreatedEvent struct {
	ReservationID uint
	Channel       string
	Recipient     string
}
