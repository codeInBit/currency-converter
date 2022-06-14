package utilities

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//JSON - This format JSON responses
func JSON(w http.ResponseWriter, statusCode int, result interface{}, wasCached bool, durationSec float64) {
	w.WriteHeader(statusCode)

	format := struct {
		Success     bool        `json:"success"`
		Result      interface{} `json:"result"`
		WasCached   bool        `json:"was_cached"`
		DurationSec float64     `json:"duration_sec"`
	}{
		Success:     true,
		Result:      result,
		WasCached:   wasCached,
		DurationSec: durationSec,
	}

	err := json.NewEncoder(w).Encode(format)
	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}
}

//ERROR - This format error messages
func ERROR(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)

	format := struct {
		Success bool   `json:"success"`
		Error   string `json:"error"`
	}{
		Success: false,
		Error:   message,
	}

	err := json.NewEncoder(w).Encode(format)
	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}
}
