package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// doLogin implementa l'operazione `doLogin` della specifica OpenAPI.
//
//	POST /session
//
// Se l'utente esiste ne restituisce l'identificativo, altrimenti lo crea.
// Risposte dichiarate: 201, 400, 500.
func (rt *Router) doLogin(w http.ResponseWriter, r *http.Request) {
	if !methodIs(w, r, http.MethodPost) {
		return
	}

	// requestBody (application/json, required: true) -> LoginRequest
	var body LoginRequest
	decoder := json.NewDecoder(io.LimitReader(r.Body, maxRequestBody))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&body); err != nil {
		rt.logger.Printf("doLogin: corpo della richiesta non valido: %v", err)

		if errors.Is(err, io.EOF) {
			writeError(w, http.StatusBadRequest, codeBadRequest, "corpo della richiesta mancante")
			return
		}
		writeError(w, http.StatusBadRequest, codeBadRequest, "corpo della richiesta non valido")
		return
	}

	// Vincoli dello schema Username (pattern, minLength, maxLength).
	if !body.IsValid() {
		writeError(w, http.StatusBadRequest, codeBadRequest,
			"il nome utente deve contenere da 3 a 16 caratteri fra lettere, numeri, punto e underscore")
		return
	}

	identifier, created := rt.users.login(body.Name)
	if created {
		rt.logger.Printf("doLogin: creato l'utente %q con identificativo %d", body.Name, identifier)
	} else {
		rt.logger.Printf("doLogin: accesso dell'utente %q (identificativo %d)", body.Name, identifier)
	}

	// Risposta 201 -> LoginResponse
	writeJSON(w, http.StatusCreated, LoginResponse{Identifier: identifier})
}
