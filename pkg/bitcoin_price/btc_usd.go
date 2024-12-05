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
}

// CSVBitcoinPriceFetcher fetches Bitcoin price data from a CSV file
type CSVBitcoinPriceFetcher struct {
	// filePath string
	prices map[time.Time]BitcoinPrice
}

// NewCSVBitcoinPriceFetcher creates a new CSVBitcoinPriceFetcher and loads the data
func NewCSVBitcoinPriceFetcher(filePath string) (*CSVBitcoinPriceFetcher, error) {
	fetcher := &CSVBitcoinPriceFetcher{
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
		open := parseFloat(record[1])
		high := parseFloat(record[2])
		low := parseFloat(record[3])
		close := parseFloat(record[4])
		volume := parseFloat(record[5])
		marketCap := parseFloat(record[6])
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
func (f *CSVBitcoinPriceFetcher) FetchHistoricalData(startDate time.Time) (map[time.Time]BitcoinPrice, error) {
	prices := make(map[time.Time]BitcoinPrice)
	for date, price := range f.prices {
		if date.After(startDate) {
			prices[date] = price
		}
	}
	return prices, nil
}

// FetchPriceByDate fetches the Bitcoin price for a specific date
func (f *CSVBitcoinPriceFetcher) FetchPriceByDate(date time.Time) (BitcoinPrice, error) {
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
