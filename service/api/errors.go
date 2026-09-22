package api

import (
	"encoding/json"
	"net/http"
)

// Codici simbolici usati nel campo `code` dello schema Error. Rispettano il
// pattern `^[a-z_]{3,32}$` dichiarato in doc/api.yaml.
const (
	codeBadRequest    = "bad_request"
	codeUnauthorized  = "unauthorized"
	codeConflict      = "conflict"
	codeInternalError = "internal_error"
)

// writeJSON serializza il corpo della risposta con il codice di stato indicato.
func writeJSON(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	// L'errore di scrittura non è recuperabile: l'header è già stato inviato.
	_ = json.NewEncoder(w).Encode(body)
}

// writeError produce una risposta conforme allo schema Error.
func writeError(w http.ResponseWriter, status int, code string, message string) {
	writeJSON(w, status, Error{Code: code, Message: message})
}
