/*
Package api implementa lo strato HTTP descritto da doc/api.yaml.

Usa soltanto la libreria standard. Le rotte sono registrate su un
http.ServeMux in New(); ogni `operationId` della specifica corrisponde a un
metodo di Router definito in un file dedicato (do-login.go, ...).
*/
package api

import (
	"log"
	"net/http"
)

// Config raccoglie le dipendenze necessarie a costruire il Router.
type Config struct {
	// Logger è il logger usato dagli handler. Se è nil viene usato quello
	// predefinito del package log.
	Logger *log.Logger
}

// Router espone le rotte dell'applicazione.
type Router struct {
	mux    *http.ServeMux
	logger *log.Logger
	users  *userStore
}

// New costruisce il Router e registra le rotte.
func New(cfg Config) *Router {
	logger := cfg.Logger
	if logger == nil {
		logger = log.Default()
	}

	rt := &Router{
		mux:    http.NewServeMux(),
		logger: logger,
		users:  newUserStore(),
	}

	// I pattern contengono solo il percorso: la forma "POST /session"
	// richiede Go 1.22 e su versioni precedenti verrebbe letta come un nome
	// di host. Il metodo HTTP lo verifica ogni handler.
	//
	//	POST /session  -> doLogin
	rt.mux.HandleFunc("/session", rt.doLogin)
	rt.mux.HandleFunc("/liveness", rt.liveness)
	rt.mux.HandleFunc("/", rt.root)

	return rt
}

// Handler restituisce l'http.Handler con tutte le rotte registrate.
func (rt *Router) Handler() http.Handler {
	return rt.mux
}

// root risponde sulla radice del sito. Non fa parte della specifica.
func (rt *Router) root(w http.ResponseWriter, r *http.Request) {
	// Il pattern "/" raccoglie ogni percorso non registrato altrove:
	// distinguiamo la radice dalle rotte inesistenti.
	if r.URL.Path != "/" {
		writeError(w, http.StatusNotFound, codeNotFound, "risorsa non trovata")
		return
	}
	if !methodIs(w, r, http.MethodGet) {
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte("WASAText: server attivo\n"))
}

// liveness segnala che il server è in grado di rispondere.
func (rt *Router) liveness(w http.ResponseWriter, r *http.Request) {
	if !methodIs(w, r, http.MethodGet) {
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
