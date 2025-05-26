package indicators

import "lab_06/models"

func ATR(data []models.StockData, period int) []float64 {
	if period <= 0 || len(data) < period+1 {
		return nil
	}
	tr := make([]float64, len(data))
	for i := 1; i < len(data); i++ {
		highLow := data[i].High - data[i].Low
		highClose := abs(data[i].High - data[i-1].Close)
		lowClose := abs(data[i].Low - data[i-1].Close)
		tr[i] = max(highLow, highClose, lowClose)
	}
	atr := make([]float64, len(data)-period)
	for i := period; i < len(data); i++ {
		sum := 0.0
		for j := i - period + 1; j <= i; j++ {
			sum += tr[j]
		}
		atr[i-period] = sum / float64(period)
	}
	return atr
}

func max(a, b, c float64) float64 {
	m := a
	if b > m {
		m = b
	}
	if c > m {
		m = c
	}
	return m
}

func abs(a float64) float64 {
	if a < 0 {
		return -a
	}
	return a
}
