package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/Lorenzo461/CV-Translate/service/api/reqcontext"
	"github.com/Lorenzo461/CV-Translate/service/database"
)

// setMyUserName implementa l'operazione `setMyUserName` della specifica OpenAPI.
//
//	PUT /users/me/username
//
// Aggiorna il nome utente dell'utente autenticato.
// Risposte dichiarate: 200, 400, 401, 409, 500.
// Il 401 è già gestito dal wrapper di autenticazione (bearerAuth).
func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// requestBody (application/json, required: true) -> ChangeUsernameRequest
	var body ChangeUsernameRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		ctx.Logger.WithError(err).Warn("setMyUserName: corpo della richiesta non decodificabile")
		writeError(w, http.StatusBadRequest, codeBadRequest, "corpo della richiesta non valido")
		return
	}

	// Vincoli dello schema Username (pattern, minLength, maxLength).
	if !body.IsValid() {
		writeError(w, http.StatusBadRequest, codeBadRequest,
			"il nome utente deve contenere da 3 a 16 caratteri fra lettere, numeri, punto e underscore")
		return
	}

	user, err := rt.db.SetMyUserName(ctx.UserID, body.Name)
	switch {
	case errors.Is(err, database.ErrUsernameTaken):
		// Risposta 409 dichiarata nello YAML.
		writeError(w, http.StatusConflict, codeConflict, "nome utente già in uso")
		return
	case errors.Is(err, database.ErrUserNotFound):
		// L'utente è stato autenticato poco fa: se è sparito, il token non è
		// più valido.
		writeError(w, http.StatusUnauthorized, codeUnauthorized, "token non più valido")
		return
	case err != nil:
		ctx.Logger.WithError(err).Error("setMyUserName: errore durante l'aggiornamento del nome utente")
		writeError(w, http.StatusInternalServerError, codeInternalError, "errore interno")
		return
	}

	// Risposta 200 -> schema User
	writeJSON(w, http.StatusOK, fromDatabaseUser(user))
}
