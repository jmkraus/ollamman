package formatbytes

import (
	"fmt"
)

// FormatBytes wandelt eine Größe in Bytes in eine lesbare Form um
func FormatBytes(bytes int64) string {
	const unit = 1000
	if bytes < unit {
		return fmt.Sprintf("%dB", bytes) // Weniger als 1KB → einfach "X B"
	}

	sizes := []string{"", "KB", "MB", "GB", "TB", "PB", "EB"}
	sizeIndex := 0
	value := float64(bytes)

	for value >= unit && sizeIndex < len(sizes)-1 {
		value /= unit
		sizeIndex++
	}

	return fmt.Sprintf("%.1f %s", value, sizes[sizeIndex])
}
