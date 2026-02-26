package history

import (
	"math"
	"sync"
)

// History maintains a rolling window of metric values for graphing
type History struct {
	mu       sync.RWMutex
	values   []float64
	maxSize  int
	label    string
	min, max float64
}

// Sparkline symbols (Unicode block elements)
var sparkSymbols = []rune{' ', '▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}

// NewHistory creates a new history tracker
func NewHistory(maxSize int, label string) *History {
	return &History{
		values:  make([]float64, 0, maxSize),
		maxSize: maxSize,
		label:   label,
		min:     math.MaxFloat64,
		max:     -math.MaxFloat64,
	}
}

// Add adds a new value to the history
func (h *History) Add(value float64) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.values = append(h.values, value)

	// Update min/max
	if value < h.min {
		h.min = value
	}
	if value > h.max {
		h.max = value
	}

	// Remove oldest if exceeding max size
	if len(h.values) > h.maxSize {
		h.values = h.values[1:]
		// Recalculate min/max
		h.recalcMinMax()
	}
}

// Values returns a copy of the values
func (h *History) Values() []float64 {
	h.mu.RLock()
	defer h.mu.RUnlock()

	result := make([]float64, len(h.values))
	copy(result, h.values)
	return result
}

// Len returns the number of values
func (h *History) Len() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.values)
}

// Average returns the average of all values
func (h *History) Average() float64 {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if len(h.values) == 0 {
		return 0
	}

	sum := 0.0
	for _, v := range h.values {
		sum += v
	}
	return sum / float64(len(h.values))
}

// Min returns the minimum value
func (h *History) Min() float64 {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if h.min == math.MaxFloat64 {
		return 0
	}
	return h.min
}

// Max returns the maximum value
func (h *History) Max() float64 {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if h.max == -math.MaxFloat64 {
		return 0
	}
	return h.max
}

// Current returns the latest value
func (h *History) Current() float64 {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if len(h.values) == 0 {
		return 0
	}
	return h.values[len(h.values)-1]
}

// Clear resets the history
func (h *History) Clear() {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.values = h.values[:0]
	h.min = math.MaxFloat64
	h.max = -math.MaxFloat64
}

// Sparkline generates a sparkline string from the values
func (h *History) Sparkline(width int) string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if len(h.values) == 0 {
		return ""
	}

	// Determine range
	minVal := h.min
	maxVal := h.max
	rangeVal := maxVal - minVal

	// If all values are the same, use middle symbol
	if rangeVal == 0 {
		rangeVal = 1
	}

	result := make([]rune, 0, width)
	step := float64(len(h.values)) / float64(width)

	for i := 0; i < width && i*int(step) < len(h.values); i++ {
		idx := int(float64(i) * step)
		if idx >= len(h.values) {
			idx = len(h.values) - 1
		}
		v := h.values[idx]

		// Normalize to 0-8 range
		normalized := (v - minVal) / rangeVal
		symbolIdx := int(normalized * 8)
		if symbolIdx > 8 {
			symbolIdx = 8
		}
		if symbolIdx < 0 {
			symbolIdx = 0
		}
		result = append(result, sparkSymbols[symbolIdx])
	}

	return string(result)
}

// recalcMinMax recalculates min and max from current values
func (h *History) recalcMinMax() {
	h.min = math.MaxFloat64
	h.max = -math.MaxFloat64
	for _, v := range h.values {
		if v < h.min {
			h.min = v
		}
		if v > h.max {
			h.max = v
		}
	}
}

// Trend returns the trend direction: 1=increasing, -1=decreasing, 0=stable
func (h *History) Trend() int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if len(h.values) < 2 {
		return 0
	}

	// Compare last 5 values average with previous 5 values average
	n := len(h.values)
	window := 5
	if n < window*2 {
		window = n / 2
	}
	if window < 1 {
		return 0
	}

	recent := 0.0
	previous := 0.0

	for i := 0; i < window; i++ {
		recent += h.values[n-1-i]
		previous += h.values[n-window-1-i]
	}

	recent /= float64(window)
	previous /= float64(window)

	diff := recent - previous
	if diff > 1.0 {
		return 1
	}
	if diff < -1.0 {
		return -1
	}
	return 0
}
