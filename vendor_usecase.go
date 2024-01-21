package true_vendor_sdk

import "context"

// VendorUseCase is an interface which should implemented
type VendorUseCase interface {
	ReservationVerification(ctx context.Context, id string, trueAssistantUserID string) (string, error)
	ReservedCallback(ctx context.Context, id string) error
	CancelReservationCallback(ctx context.Context, id string) error
	ResolvingVerification(ctx context.Context, id string, token, verdict, note string) (string, error)
	ResolvingCallback(ctx context.Context, id string) error
}
