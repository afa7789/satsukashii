package bitcoin_price

import (
	"fmt"
	"sort"
	"time"
)

// BTCPricesRanged holds the ranged Bitcoin prices and corresponding dates.
type BTCPricesRanged struct {
	Prices []float64
	Dates  []string
}

// BtcRange samples Bitcoin price data over a specified number of items, starting from a fixed date.
// It fetches historical data from the provided BitcoinPriceFetcher and returns a subset of dates and prices,
// sorted chronologically. If numberOfItems is greater than the available data, it returns all data.
// Errors from date parsing or data fetching are returned if they occur.
func BtcRange(btcData BitcoinPriceFetcher, numberOfItems int) (BTCPricesRanged, error) {
	// Validate input
	if numberOfItems <= 0 {
		return BTCPricesRanged{}, fmt.Errorf("numberOfItems must be positive, got %d", numberOfItems)
	}

	// Parse start date with error handling
	timeStart, err := time.Parse("2006-01-02", "2010-07-17")
	if err != nil {
		return BTCPricesRanged{}, fmt.Errorf("failed to parse start date: %v", err)
	}

	// Fetch historical data with error handling
	returned, err := btcData.FetchHistoricalData(timeStart)
	if err != nil {
		return BTCPricesRanged{}, fmt.Errorf("failed to fetch historical data: %v", err)
	}

	// Handle empty data case
	if len(returned) == 0 {
		return BTCPricesRanged{Prices: []float64{}, Dates: []string{}}, nil
	}

	// Convert the map to a slice of key-value pairs and sort by date
	type datePricePair struct {
		date  time.Time
		price BitcoinPrice
	}
	pairs := make([]datePricePair, 0, len(returned))
	for date, price := range returned {
		pairs = append(pairs, datePricePair{date: date, price: price})
	}
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].date.Before(pairs[j].date)
	})

	// Determine the number of items to sample
	totalItems := len(pairs)
	if numberOfItems > totalItems {
		numberOfItems = totalItems // Return all items if requested more than available
	}

	// Sample data points
	dates := make([]string, 0, numberOfItems)
	prices := make([]float64, 0, numberOfItems)
	step := float64(totalItems-1) / float64(numberOfItems-1) // More precise step size
	for i := 0; i < numberOfItems; i++ {
		index := int(float64(i) * step)
		if index >= totalItems {
			index = totalItems - 1
		}
		dates = append(dates, pairs[index].date.Format("2006-01-02"))
		prices = append(prices, pairs[index].price.Close)
		// log.Printf("Date: %v, Price: %v", pairs[index].date, pairs[index].price.Close)
	}

	return BTCPricesRanged{Prices: prices, Dates: dates}, nil
}
