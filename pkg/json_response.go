package pkg

import (
	"encoding/json"
	"net/http"
)

func WriteJson(w http.ResponseWriter, status int, msg any) error {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(msg)
}
