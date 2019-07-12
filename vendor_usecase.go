package true_vendor_sdk

// VendorUseCase is an interface which should implemented
type VendorUseCase interface {
	ReservationVerification(id string, trueAssistantUserID string) (string, error)
	ReservedCallback(id string) error
	CancelReservationCallback(id string) error
	ResolvingVerification(id string, token, verdict, note string) (string, error)
	ResolvingCallback(id string) error
}
