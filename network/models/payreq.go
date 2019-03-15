package models

type PayReq struct {
	Destination     string
	PaymentHash     string
	Amount          int64
	Timestamp       int64
	Expiry          int64
	Description     string
	DescriptionHash string
	FallbackAddr    string
	CltvExpiry      int64
	String          string
}
