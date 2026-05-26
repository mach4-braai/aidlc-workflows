package trendreports

// Sparkline maps a slice of values to Unicode block characters (▁▂▃▄▅▆▇█).
// Each value is expected to be in [0, 1]. The output rune count equals len(values).
func Sparkline(values []float64) string {
	const blocks = "▁▂▃▄▅▆▇█"
	runes := []rune(blocks)
	n := len(runes)

	if len(values) == 0 {
		return ""
	}

	min, max := values[0], values[0]
	for _, v := range values[1:] {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}

	result := make([]rune, len(values))
	for i, v := range values {
		var idx int
		if max == min {
			idx = n / 2
		} else {
			idx = int((v - min) / (max - min) * float64(n-1))
		}
		if idx < 0 {
			idx = 0
		}
		if idx >= n {
			idx = n - 1
		}
		result[i] = runes[idx]
	}
	return string(result)
}

// CheckGate returns true if value meets or exceeds threshold.
func CheckGate(value, threshold float64) bool {
	return value >= threshold
}
