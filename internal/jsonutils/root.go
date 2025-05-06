package jsonutils

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Return indented, formatted JSON from bytes
func FormatJson(raw []byte) (string, error) {
	var prettyJSON bytes.Buffer

	// Indent JSON with 2 spaces
	err := json.Indent(&prettyJSON, raw, "", "    ")
	if err != nil {
		return "", fmt.Errorf("error indenting JSON: %v", err)
	}

	return prettyJSON.String(), nil
}
