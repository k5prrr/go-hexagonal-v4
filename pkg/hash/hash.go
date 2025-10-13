package hash

import (
	"crypto/sha256" // 512
	"encoding/hex"
	"fmt"
	"strings"
)

func Numeric(text string, maxLength int) string {
	hash := sha256.Sum256([]byte(text))
	result := StringToNumeric(hash[:])

	if maxLength == 0 {
		maxLength = 64
	}
	if len(result) > maxLength {
		result = result[:64]
	}
	return result
}

func StringToNumeric(text []byte) string {
	var result strings.Builder
	for _, char := range text {
		result.WriteString(fmt.Sprintf("%d", char))
	}
	return result.String()
}

func String64(text string) string {
	hash := sha256.Sum256([]byte(text))
	return hex.EncodeToString(hash[:])
}
