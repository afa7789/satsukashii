package bigmac

import (
	"encoding/csv"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

// BigMacData holds the nested map for consulting
type BigMacData struct {
	data map[string]map[string]float64
}

const numberOfBits = 64

// NewBigMacData constructs the BigMacData from the given CSV file path
func NewBigMacData(filePath string) (*BigMacData, error) {
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

	data := make(map[string]map[string]float64)
	for _, record := range records[1:] { // Skip header
		// log record
		// log.Printf("record: %v", record)
		length := len(record)

		currencyCode := record[2]
		date := record[length-1]
		priceLocal, err := strconv.ParseFloat(record[3], numberOfBits)
		if err != nil {
			// first try position
			// log.Errorf("skipping record with invalid GDP local value: %v, invalid: %d,", record, length)
			// for i, v := range record {
			// 	log.Printf("record[%d]: %v", i, v)
			// }
			continue
		}

		if _, exists := data[currencyCode]; !exists {
			data[currencyCode] = make(map[string]float64)
		}
		data[currencyCode][date] = priceLocal
	}

	return &BigMacData{data: data}, nil
}

// NewBigMacDataFromURL constructs the BigMacData from the given CSV URL
func NewBigMacDataFromURL(url string) (*BigMacData, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	tmpFile, err := os.CreateTemp("", "bigmacdata-*.csv")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(body); err != nil {
		return nil, err
	}
	if err := tmpFile.Close(); err != nil {
		return nil, err
	}

	return NewBigMacData(tmpFile.Name())
}

// GetPrice returns the GDP local price for the given currency code and date
func (b *BigMacData) GetPrice(currencyCode, date string) (float64, bool) {
	if dateMap, exists := b.data[currencyCode]; exists {
		price, exists := dateMap[date]
		return price, exists
	}
	return 0, false
}

// GetPriceTimestamp returns the price and timestamp of the closest date for the given currency code and date
func (b *BigMacData) GetPriceTimestamp(currencyCode string, timestamp time.Time) (float64, time.Time, bool) {
	dateMap, exists := b.data[currencyCode]
	if !exists {
		return 0, time.Time{}, false
	}

	var closestTime time.Time
	var price float64

	for dateStr, priceValue := range dateMap {
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			continue
		}
		if date.Before(timestamp) || date.Equal(timestamp) {
			if closestTime.IsZero() || date.After(closestTime) {
				closestTime = date
				price = priceValue
			}
		}
	}

	if closestTime.IsZero() {
		return 0, time.Time{}, false
	}

	return price, closestTime, true
}

// GetPriceInBitcoin returns the price of the BigMac in Bitcoin for the given currency code and timestamp
func (b *BigMacData) GetPriceInBitcoin(currencyCode string, timestamp time.Time, bitcoinValueAtTime float64) (float64, bool) {
	dateMap, exists := b.data[currencyCode]
	if !exists {
		return 0, false
	}

	var closestDate string
	var closestTime time.Time

	for dateStr := range dateMap {
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			continue
		}
		if date.Before(timestamp) || date.Equal(timestamp) {
			if closestTime.IsZero() || date.After(closestTime) {
				closestTime = date
				closestDate = dateStr
			}
		}
	}

	if closestDate == "" {
		return 0, false
	}

	priceInLocalCurrency := dateMap[closestDate]
	priceInBitcoin := priceInLocalCurrency / bitcoinValueAtTime
	return priceInBitcoin, true
}

// func (bmd *BigMacData) GetTrendLine(currencyCode string) (float64, float64, error) {
// 	dateMap, exists := bmd.data[currencyCode]
// 	if !exists {
// 		return 0, 0, fmt.Errorf("currency code not found")
// 	}

// 	// Parse dates to find the earliest date
// 	var minDate time.Time
// 	for dateStr := range dateMap {
// 		date, err := time.Parse("2006-01-02", dateStr)
// 		if err != nil {
// 			return 0, 0, fmt.Errorf("invalid date format: %s", err)
// 		}
// 		if minDate.IsZero() || date.Before(minDate) {
// 			minDate = date
// 		}
// 	}

// 	if minDate.IsZero() {
// 		return 0, 0, fmt.Errorf("no valid dates found")
// 	}

// 	var sumX, sumY, sumXY, sumX2 float64
// 	var n float64

// 	// Use the number of days since the earliest date as x values
// 	for dateStr, price := range dateMap {
// 		date, err := time.Parse("2006-01-02", dateStr)
// 		if err != nil {
// 			return 0, 0, fmt.Errorf("invalid date format: %s", err)
// 		}

// 		// Days since the earliest date
// 		x := float64(date.Sub(minDate).Hours() / 24)
// 		y := price

// 		sumX += x
// 		sumY += y
// 		sumXY += x * y
// 		sumX2 += x * x
// 		n++
// 	}

// 	if n == 0 {
// 		return 0, 0, fmt.Errorf("no data points")
// 	}

// 	// Calculate slope (m) and intercept (b) for y = mx + b
// 	m := (n*sumXY - sumX*sumY) / (n*sumX2 - sumX*sumX)
// 	b := (sumY - m*sumX) / n

// 	return m, b, nil
// }
