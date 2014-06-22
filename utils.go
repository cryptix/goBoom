package goBoom

import (
	"encoding/json"
	"fmt"
)

func reverse(input string) string {
	runes := make([]rune, len(input))
	n := 0
	for _, r := range input {
		runes[n] = r
		n++
	}
	runes = runes[0:n]

	// Reverse
	for i := 0; i < n/2; i++ {
		runes[i], runes[n-1-i] = runes[n-1-i], runes[i]
	}

	// Convert back to UTF-8.
	return string(runes)
}

// ugly way to help with multiple unknown json return values
func jsonRemarshal(in, out interface{}) (err error) {
	// ugly remarshall...
	tmp, err := json.Marshal(in)
	if err != nil {
		return fmt.Errorf("Remashal failed during Marshal():%s\n", err)
	}

	err = json.Unmarshal(tmp, out)
	if err != nil {
		return fmt.Errorf("Remashal failed during Unmarshal():%s\n", err)
	}

	return nil
}
