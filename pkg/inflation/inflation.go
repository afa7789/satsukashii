package inflation

import (
	"encoding/csv"
	"os"
	"strconv"
	"time"
)

// CalculateInflation calculates the inflation rate based on an array of monthly inflation rates.
// The starting month and ending month are used as the range for the calculation.
// this works with the array arthur did in the bar with elton
func CalculateInflation(monthlyInflation []float64, startingMonth, endingMonth int) []float64 {
	if startingMonth < 0 || startingMonth >= len(monthlyInflation) || endingMonth < 0 || endingMonth >= len(monthlyInflation) || startingMonth > endingMonth {
		return nil
	}

	baseInflation := monthlyInflation[startingMonth]
	inflationRates := make([]float64, endingMonth-startingMonth+1)

	for i := startingMonth; i <= endingMonth; i++ {
		inflationRates[i-startingMonth] = (monthlyInflation[i] - baseInflation) / baseInflation * 100
	}

	return inflationRates
}

type CPIInflation struct {
	data map[int][]float64
}

func NewCPIInflation(filePath string) (*CPIInflation, error) {
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

	cpiInflation := &CPIInflation{
		data: make(map[int][]float64),
	}

	for _, record := range records[1:] {
		year, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, err
		}

		monthlyInflation := make([]float64, 12)
		for i := 1; i <= 12; i++ {
			if record[i] != "" {
				inflation, err := strconv.ParseFloat(record[i], 64)
				monthlyInflation[i-1] = 1 + inflation
				if err != nil {
					return nil, err
				}
			}
		}

		cpiInflation.data[year] = monthlyInflation
	}

	return cpiInflation, nil
}

func (cpi *CPIInflation) CalculateInflationSince(start, end time.Time, currentValue float64) float64 {
	accumulatedInflation := 1.0 // start with 1 (no inflation)

	// Loop through the years from start year to end year
	for i := start.Year(); i <= end.Year(); i++ {
		monthlyInflation := cpi.data[i]
		if monthlyInflation == nil {
			continue
		}

		// Loop through months (January is 0, December is 11)
		for j := 0; j < len(monthlyInflation); j++ {
			// Generate the month and year to compare with the range
			newDate := time.Date(i, time.Month(j+1), 1, 0, 0, 0, 0, time.UTC)

			// Only apply inflation if the month is within the range [start, end]
			if newDate.After(start) && newDate.Before(end) || newDate.Equal(start) || newDate.Equal(end) {
				accumulatedInflation *= monthlyInflation[j]
			}
		}
	}

	// Apply the accumulated inflation to the current value
	return currentValue * accumulatedInflation
}
