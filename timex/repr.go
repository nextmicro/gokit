package timex

import (
	"fmt"
	"strconv"
	"time"
)

// Duration returns the float64 representation of given duration in ms.
func Duration(duration time.Duration) float64 {
	v := fmt.Sprintf("%.3f", float32(duration)/float32(time.Millisecond))
	float, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return 0
	}
	return float
}

// ReprOfDuration returns the string representation of given duration in ms.
func ReprOfDuration(duration time.Duration) string {
	return fmt.Sprintf("%.3fms", float32(duration)/float32(time.Millisecond))
}
