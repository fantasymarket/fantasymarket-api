package structs

// StockSettings is the type for storing information about stocks
type StockSettings struct {
	StockID   string          `yaml:"stockID"`    // Stock Symbol e.g GOOG
	Name      string          `yaml:"name"`       // Stock Name e.g Alphabet Inc.
	Index     int64           `yaml:"startPrice"` // Price per share
	Shares    int64           `yaml:"stockCount"` // Number per share
	Tags      map[string]bool `yaml:"tags"`       // A stock can have up to 5 tags
	Stability int64           `yaml:"stability"`  // Shows how many fluctuations the stock will have
	Trend     int64           `yaml:"trend"`      // Shows the generall trend of the Stock
}
