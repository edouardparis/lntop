package models

type Transaction struct {
	// / The transaction hash
	TxHash string
	// / The transaction amount, denominated in satoshis
	Amount int64
	// / The number of confirmations
	NumConfirmations int32
	// / The hash of the block this transaction was included in
	BlockHash string
	// / The height of the block this transaction was included in
	BlockHeight int32
	// / Timestamp of this transaction
	TimeStamp int64
	// / Fees paid for this transaction
	TotalFees int64
	// / Addresses that received funds for this transaction
	DestAddresses []string
}
