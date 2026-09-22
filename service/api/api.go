/*
Package api implementa lo strato HTTP descritto da doc/api.yaml.

Ogni `operationId` della specifica corrisponde a un metodo di _router definito
in un file dedicato (do-login.go, set-my-username.go, ...); la registrazione
delle rotte avviene in api-handler.go e ricalca la sezione `paths` dello YAML.
*/
package api

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"

	"github.com/Lorenzo461/CV-Translate/service/database"
)

// Router è l'interfaccia pubblica del package: espone l'handler HTTP
// dell'applicazione.
type Router interface {
	// Handler restituisce l'http.Handler con tutte le rotte registrate.
	Handler() http.Handler

	// Close rilascia le risorse associate al router.
	Close() error
}

// Config raccoglie le dipendenze necessarie a costruire il Router.
type Config struct {
	// Logger è il logger di base dell'applicazione (obbligatorio).
	Logger logrus.FieldLogger

	// Database è lo strato di persistenza (obbligatorio).
	Database database.AppDatabase
}

// New costruisce un nuovo Router a partire dalla configurazione indicata.
func New(cfg Config) (Router, error) {
	if cfg.Logger == nil {
		return nil, errors.New("logger is required")
	}
	if cfg.Database == nil {
		return nil, errors.New("database is required")
	}

	router := httprouter.New()
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false

	return &_router{
		router:     router,
		baseLogger: cfg.Logger,
		db:         cfg.Database,
	}, nil
}

// _router è l'implementazione non esportata di Router.
type _router struct {
	router     *httprouter.Router
	baseLogger logrus.FieldLogger
	db         database.AppDatabase
}

// Close implementa Router.
func (rt *_router) Close() error {
	return nil
}
