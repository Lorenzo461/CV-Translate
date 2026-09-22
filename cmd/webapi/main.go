/*
webapi è il punto di ingresso del backend WASAText.

Apre il database SQLite, costruisce il router descritto da doc/api.yaml e
avvia il server HTTP, che viene spento in modo pulito alla ricezione di
SIGINT o SIGTERM.

Utilizzo:

	go run ./cmd/webapi [--api-host :3000] [--db-filename ./wasatext.db]
*/
package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	_ "github.com/mattn/go-sqlite3"

	"github.com/Lorenzo461/CV-Translate/service/api"
	"github.com/Lorenzo461/CV-Translate/service/database"
)

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "fatal error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	apiHost := flag.String("api-host", ":3000", "indirizzo di ascolto del server HTTP")
	dbFilename := flag.String("db-filename", "./wasatext.db", "percorso del file SQLite")
	flag.Parse()

	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	logger.SetLevel(logrus.InfoLevel)
	logger.Info("avvio di WASAText")

	// Database
	logger.WithField("filename", *dbFilename).Info("apertura del database")
	dbConn, err := sql.Open("sqlite3", *dbFilename+"?_foreign_keys=on")
	if err != nil {
		return fmt.Errorf("error opening the SQLite database: %w", err)
	}
	defer func() {
		if err := dbConn.Close(); err != nil {
			logger.WithError(err).Error("errore nella chiusura del database")
		}
	}()

	db, err := database.New(dbConn)
	if err != nil {
		return fmt.Errorf("error building the AppDatabase: %w", err)
	}

	// Router
	router, err := api.New(api.Config{Logger: logger, Database: db})
	if err != nil {
		return fmt.Errorf("error building the API router: %w", err)
	}
	defer func() {
		if err := router.Close(); err != nil {
			logger.WithError(err).Error("errore nella chiusura del router")
		}
	}()

	server := &http.Server{
		Addr:              *apiHost,
		Handler:           router.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
	}

	// Avvio del server e attesa di un segnale di terminazione.
	serverErrors := make(chan error, 1)
	go func() {
		logger.WithField("addr", server.Addr).Info("server HTTP in ascolto")
		serverErrors <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("error starting the HTTP server: %w", err)
		}
	case sig := <-shutdown:
		logger.WithField("signal", sig.String()).Info("spegnimento in corso")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			// Lo spegnimento pulito non è riuscito entro il timeout.
			_ = server.Close()
			return fmt.Errorf("error shutting down the HTTP server: %w", err)
		}
	}

	logger.Info("arrivederci")
	return nil
}
