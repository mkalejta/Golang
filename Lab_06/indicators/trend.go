package indicators

func SMA(prices []float64, period int) []float64 {
	if period <= 0 || len(prices) < period {
		return nil
	}
	sma := make([]float64, len(prices)-period+1)
	for i := 0; i <= len(prices)-period; i++ {
		sum := 0.0
		for j := 0; j < period; j++ {
			sum += prices[i+j]
		}
		sma[i] = sum / float64(period)
	}
	return sma
}
