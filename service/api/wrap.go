package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"

	"github.com/Lorenzo461/CV-Translate/service/api/reqcontext"
)

// httpRouterHandler è la firma degli handler del package: come quella di
// httprouter.Handle, ma con in più il contesto di richiesta.
type httpRouterHandler func(http.ResponseWriter, *http.Request, httprouter.Params, reqcontext.RequestContext)

// wrap adatta un httpRouterHandler alla firma richiesta da httprouter,
// costruendo il contesto della richiesta.
//
// Il parametro authRequired implementa il `securityScheme` bearerAuth della
// specifica: vale true per le operazioni che ereditano la `security` globale
// (setMyUserName) e false per quelle che dichiarano `security: []` (doLogin).
func (rt *_router) wrap(fn httpRouterHandler, authRequired bool) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		reqUUID, err := uuid.NewV4()
		if err != nil {
			rt.baseLogger.WithError(err).Error("can't generate a request UUID")
			writeError(w, http.StatusInternalServerError, codeInternalError, "errore interno")
			return
		}

		ctx := reqcontext.RequestContext{ReqUUID: reqUUID}
		ctx.Logger = rt.baseLogger.WithFields(logFields(r, ctx.ReqUUID.String()))

		if authRequired {
			userID, ok := rt.authenticate(r)
			if !ok {
				// Risposta 401 dichiarata in doc/api.yaml.
				writeError(w, http.StatusUnauthorized, codeUnauthorized,
					"token di autenticazione mancante o non valido")
				return
			}
			ctx.UserID = userID
			ctx.Logger = ctx.Logger.WithField("userid", userID)
		}

		fn(w, r, ps, ctx)
	}
}

// authenticate estrae l'identificativo utente dall'header
// `Authorization: Bearer <identifier>` e verifica che l'utente esista.
func (rt *_router) authenticate(r *http.Request) (uint64, bool) {
	header := r.Header.Get("Authorization")
	if !strings.HasPrefix(header, "Bearer ") {
		return 0, false
	}

	token := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
	userID, err := strconv.ParseUint(token, 10, 64)
	if err != nil || !isValidIdentifier(userID) {
		return 0, false
	}

	if _, err := rt.db.GetUserByID(userID); err != nil {
		// Sia l'utente inesistente sia un errore del database portano a un
		// rifiuto: non riveliamo al client quale dei due si è verificato.
		return 0, false
	}

	return userID, true
}

// logFields raccoglie i campi comuni a tutte le righe di log di una richiesta.
func logFields(r *http.Request, reqID string) map[string]interface{} {
	return map[string]interface{}{
		"reqid":     reqID,
		"method":    r.Method,
		"path":      r.URL.Path,
		"remote-ip": r.RemoteAddr,
	}
}
