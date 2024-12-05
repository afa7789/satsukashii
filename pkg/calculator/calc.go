package calc

// CalculateBitcoinAmount calculates the amount of Bitcoin that can be bought with a given amount of money at a given price.
func CalculateBitcoinAmount(price, amount float64) float64 {
	return amount / price
}

// CalculateMoneyAmount calculates the amount of money that would be obtained by selling a given amount of Bitcoin at a given price.
func CalculateMoneyAmount(price, bitcoinAmount float64) float64 {
	return bitcoinAmount * price
}

// CalculateBitcoinPrice calculates the price of Bitcoin given the amount of money and the amount of Bitcoin.
func CalculateBitcoinPrice(amount, bitcoinAmount float64) float64 {
	return amount / bitcoinAmount
}

// CalculateGains calculates the gains from an investment in Bitcoin given the initial amount of money invested, the amount of Bitcoin bought, and the current price of Bitcoin.
func CalculateGains(initialInvestment, bitcoinAmount, currentPrice float64) float64 {
	currentInvestmentValue := bitcoinAmount * currentPrice
	return currentInvestmentValue - initialInvestment
}
