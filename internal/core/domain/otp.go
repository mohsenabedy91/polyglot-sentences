package domain

type OTP struct {
	Value        string
	Used         bool
	RequestCount int8
	CreatedAt    int64
	LastRequest  int64
}
