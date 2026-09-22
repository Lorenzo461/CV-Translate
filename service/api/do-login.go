package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/Lorenzo461/CV-Translate/service/api/reqcontext"
)

// doLogin implementa l'operazione `doLogin` della specifica OpenAPI.
//
//	POST /session
//
// Se l'utente esiste ne restituisce l'identificativo, altrimenti lo crea.
// Risposte dichiarate: 201, 400, 500.
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// requestBody (application/json, required: true) -> LoginRequest
	var body LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		ctx.Logger.WithError(err).Warn("doLogin: corpo della richiesta non decodificabile")
		writeError(w, http.StatusBadRequest, codeBadRequest, "corpo della richiesta non valido")
		return
	}

	// Vincoli dello schema Username (pattern, minLength, maxLength).
	if !body.IsValid() {
		writeError(w, http.StatusBadRequest, codeBadRequest,
			"il nome utente deve contenere da 3 a 16 caratteri fra lettere, numeri, punto e underscore")
		return
	}

	identifier, err := rt.db.DoLogin(body.Name)
	if err != nil {
		ctx.Logger.WithError(err).Error("doLogin: errore durante il login")
		writeError(w, http.StatusInternalServerError, codeInternalError, "errore interno")
		return
	}

	// Risposta 201 -> LoginResponse
	writeJSON(w, http.StatusCreated, LoginResponse{Identifier: identifier})
}
