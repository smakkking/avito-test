package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func JSONwithCode(w http.ResponseWriter, r *http.Request, code int, v interface{}) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	if err := enc.Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(buf.Bytes()) //nolint:errcheck
}

func ErrMessage(text string) struct {
	Err string `json:"error"`
} {
	return struct {
		Err string `json:"error"`
	}{Err: text}
}
