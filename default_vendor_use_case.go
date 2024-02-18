package true_vendor_sdk

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"github.com/techpro-studio/gohttplib"
)

type ApprovalData struct {
	Verdict         *string `bson:"verdict"`
	VerdictNote     *string `bson:"verdict_note"`
	TrueEntityId    string  `bson:"true_entity_id"`
	TrueAssistantID *string `bson:"true_assistant_id"`
	ResolvingHash   *string `bson:"resolving_hash"`
	ResolvedHash    *string `bson:"resolved_hash"`
}

type VendorModel interface {
	SetVerificationState()
	SetPreVerificationState()
	SetApprovedState()
	SetRejectedState()
	GetApprovalData() *ApprovalData
}

func IsEqualStrings(lhs *string, rhs *string) bool {
	if lhs == rhs {
		return true
	} else {
		if lhs == nil || rhs == nil {
			return false
		} else {
			return *lhs == *rhs
		}
	}
}

type BaseVendorModelRepository[M VendorModel] interface {
	GetOneVendorModel(ctx context.Context, id string) (M, error)
	SaveVendorModel(ctx context.Context, object M) error
}

type DefaultVendorUseCase[M VendorModel] struct {
	repository BaseVendorModelRepository[M]
	auth       Authenticating
}

func NewDefaultVendorUseCase[M VendorModel](repository BaseVendorModelRepository[M], auth Authenticating) *DefaultVendorUseCase[M] {
	return &DefaultVendorUseCase[M]{repository: repository, auth: auth}
}

func (d *DefaultVendorUseCase[M]) ReservationVerification(ctx context.Context, id string, trueAssistantUserID string) (string, error) {
	obj, err := d.repository.GetOneVendorModel(ctx, id)
	if err != nil {
		return "", err
	}
	obj.SetVerificationState()
	data := obj.GetApprovalData()
	data.TrueAssistantID = &trueAssistantUserID
	err = d.repository.SaveVendorModel(ctx, obj)
	if err != nil {
		return "", err
	}
	return d.auth.GenerateTokenFor(trueAssistantUserID, obj.GetApprovalData().TrueEntityId)
}

func (d *DefaultVendorUseCase[M]) ReservedCallback(ctx context.Context, id string) error {
	return nil
}

func (d *DefaultVendorUseCase[M]) CancelReservationCallback(ctx context.Context, id string) error {
	obj, err := d.repository.GetOneVendorModel(ctx, id)
	if err != nil {
		return err
	}
	obj.SetPreVerificationState()
	return d.repository.SaveVendorModel(ctx, obj)
}

func (d *DefaultVendorUseCase[M]) ResolvingVerification(ctx context.Context, id string, token, verdict, note string) (string, error) {
	obj, err := d.repository.GetOneVendorModel(ctx, id)
	if err != nil {
		return "", err
	}
	if !IsEqualStrings(obj.GetApprovalData().ResolvingHash, &token) {
		return "", gohttplib.HTTP400("TOKEN_MISMATCH")
	}

	data := obj.GetApprovalData()

	hasher := sha256.New()

	hasher.Write([]byte(*data.ResolvingHash + *data.TrueAssistantID))
	hash := hex.EncodeToString(hasher.Sum(nil))

	data.ResolvedHash = &hash
	data.Verdict = &verdict
	data.VerdictNote = &note

	if verdict == "halal" {
		obj.SetApprovedState()
	} else {
		obj.SetRejectedState()
	}
	err = d.repository.SaveVendorModel(ctx, obj)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func (d *DefaultVendorUseCase[M]) ResolvingCallback(ctx context.Context, id string) error {
	return nil
}
