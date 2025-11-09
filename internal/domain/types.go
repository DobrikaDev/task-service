package domain

type VerificationType string

const (
	VerificationTypeKYC   VerificationType = "kyc"
	VerificationTypeNone  VerificationType = "none"
	VerificationTypeOther VerificationType = "other"
)

func (v VerificationType) String() string {
	return string(v)
}
