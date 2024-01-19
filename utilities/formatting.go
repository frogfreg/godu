package utilities

import "fmt"

func HumanReadableByteString(byteCount int) string {
	var quantity float64
	var suffix string

	kilobytes := float64(byteCount) / 1024
	megabytes := kilobytes / 1024
	gigabytes := megabytes / 1024

	if gigabytes >= 1 {
		quantity = gigabytes
		suffix = "GB"
	} else if megabytes >= 1 {
		quantity = megabytes
		suffix = "MB"
	} else if kilobytes >= 1 {
		quantity = kilobytes
		suffix = "KB"
	} else {
		quantity = float64(byteCount)
		suffix = "Bytes"
	}

	return fmt.Sprintf("%.2f %v", quantity, suffix)
}
