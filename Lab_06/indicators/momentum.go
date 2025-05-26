package indicators

func RSI(prices []float64, period int) []float64 {
	if period <= 0 || len(prices) < period+1 {
		return nil
	}
	rsi := make([]float64, len(prices)-period)
	for i := period; i < len(prices); i++ {
		gain, loss := 0.0, 0.0
		for j := i - period + 1; j <= i; j++ {
			change := prices[j] - prices[j-1]
			if change > 0 {
				gain += change
			} else {
				loss -= change
			}
		}
		avgGain := gain / float64(period)
		avgLoss := loss / float64(period)
		if avgLoss == 0 {
			rsi[i-period] = 100
		} else {
			rs := avgGain / avgLoss
			rsi[i-period] = 100 - (100 / (1 + rs))
		}
	}
	return rsi
}
