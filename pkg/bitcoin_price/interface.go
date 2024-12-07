package bitcoin_price

import "time"

// BitcoinPriceFetcher is an interface for fetching Bitcoin price data
type BitcoinPriceFetcher interface {
	FetchHistoricalData(startDate time.Time) (map[time.Time]BitcoinPrice, error)
	FetchPriceByDate(date time.Time) (BitcoinPrice, error)
}

// BitcoinPrice represents the structure of the Bitcoin price data
type BitcoinPrice struct {
	Start        time.Time
	Open         float64
	High         float64
	Low          float64
	Close        float64
	Volume       float64
	MarketCap    float64
	CurrencyCode string
}
