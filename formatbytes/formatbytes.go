package formatbytes

import (
	"fmt"
)

// FormatBytes converts a value of bytes into a readable output
// Uses decimal prefixes (KB, MB, GB...) with base 1000
func FormatBytes(bytes int64) string {
	const unit = 1000

	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	units := []string{"KB", "MB", "GB", "TB", "PB", "EB"}
	size := float64(bytes) / unit
	unitIndex := 0

	for size >= unit && unitIndex < len(units)-1 {
		size /= unit
		unitIndex++
	}

	return fmt.Sprintf("%.1f %s", size, units[unitIndex])
}
