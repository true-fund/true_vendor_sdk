package true_vendor_sdk


type VendorUseCase interface {
	ReservationVerification(id string, trueAssistantUserId string)(string, error)
	ReservedCallback(id string)error
	CancelReservationCallback(id string)error
	ResolvingVerification(id string, token, verdict, note string)(string, error)
	ResolvingCallback(id string)error
}
