package inflation

// CalculateInflation calculates the inflation rate based on an array of monthly inflation rates.
// The starting month and ending month are used as the range for the calculation.
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
