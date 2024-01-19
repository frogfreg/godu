package utilities

import (
	"fmt"
	"testing"
)

func TestHumanReadableByteString(t *testing.T) {
	tests := []struct {
		byteCount int
		expected  string
	}{
		{1024, "1.00 KB"},
		{2048, "2.00 KB"},
		{1024 * 1024, "1.00 MB"},
		{2 * 1024 * 1024, "2.00 MB"},
		{1024 * 1024 * 1024, "1.00 GB"},
		{2 * 1024 * 1024 * 1024, "2.00 GB"},
		{0, "0.00 Bytes"},
		{500, "500.00 Bytes"},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Input_%d", test.byteCount), func(t *testing.T) {
			result := HumanReadableByteString(test.byteCount)
			if result != test.expected {
				t.Errorf("Expected %s, but got %s", test.expected, result)
			}
		})
	}
}
