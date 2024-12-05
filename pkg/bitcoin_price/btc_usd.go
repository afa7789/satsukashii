package bitcoin_price

import (
	"encoding/csv"
	"os"
	"strconv"
	"time"
)

// BitcoinPrice represents the structure of the Bitcoin price data
type BitcoinPrice struct {
	Start     time.Time
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    float64
	MarketCap float64
}

// BitcoinPriceFetcher is an interface for fetching Bitcoin price data
type BitcoinPriceFetcher interface {
	FetchHistoricalData(startDate time.Time) (map[time.Time]BitcoinPrice, error)
	FetchPriceByDate(date time.Time) (BitcoinPrice, error)
}

// BTCPricesCSV fetches Bitcoin price data from a CSV file
type BTCPricesCSV struct {
	// filePath string
	prices map[time.Time]BitcoinPrice
}

// NewBTCPricesCSV creates a new BTCPricesCSV and loads the data
func NewBTCPricesCSV(filePath string) (*BTCPricesCSV, error) {
	fetcher := &BTCPricesCSV{
		// filePath: filePath,
		prices: make(map[time.Time]BitcoinPrice),
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	for _, record := range records[1:] { // Skip header
		start, _ := time.Parse("2006-01-02", record[0])
		// _ := time.Parse("2006-01-02", record[1]) end data
		open := parseFloat(record[2])
		high := parseFloat(record[3])
		low := parseFloat(record[4])
		close := parseFloat(record[5])
		volume := parseFloat(record[6])
		marketCap := parseFloat(record[7])
		fetcher.prices[start] = BitcoinPrice{
			Start:     start,
			Open:      open,
			High:      high,
			Low:       low,
			Close:     close,
			Volume:    volume,
			MarketCap: marketCap,
		}
	}

	return fetcher, nil
}

// FetchHistoricalData fetches historical Bitcoin price data from the preloaded data
func (f *BTCPricesCSV) FetchHistoricalData(startDate time.Time) (map[time.Time]BitcoinPrice, error) {
	prices := make(map[time.Time]BitcoinPrice)
	for date, price := range f.prices {
		if date.After(startDate) {
			prices[date] = price
		}
	}
	return prices, nil
}

// FetchPriceByDate fetches the Bitcoin price for a specific date
func (f *BTCPricesCSV) FetchPriceByDate(date time.Time) (BitcoinPrice, error) {
	price, exists := f.prices[date]
	if !exists {
		return BitcoinPrice{}, os.ErrNotExist
	}
	return price, nil
}

func parseFloat(value string) float64 {
	result, _ := strconv.ParseFloat(value, 64)
	return result
}
