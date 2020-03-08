package models

// Stock is the Stock "Class"
type Stock struct {
	StockID string // Stock Symbol e.g GOOG
	Name    string // Stock Name e.g Alphabet Inc.
	Index   int64  // Price per share
	Shares  int64  // Number of shares
}
