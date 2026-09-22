package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// liveness risponde 204 se l'applicazione è in grado di servire richieste.
// Non fa parte della specifica OpenAPI: serve agli orchestratori.
func (rt *_router) liveness(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if err := rt.db.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
