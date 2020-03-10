package models

// Stock is the Stock "Class"
type Stock struct {
	// Stock Symbol e.g GOOG
	StockID string

	// Stock Name e.g Alphabet Inc.
	Name    string

	// Price per share
	Index   int64

	// Volume since last tick, we'll have to invent this shit
	//    Calculated based on
	// 		-	the change of the index from the last tick,
	// 		- total index (so expensive stocks have larger volume than cheaper ones)
	// 		- random fluctuation
	Volume  int64
}
