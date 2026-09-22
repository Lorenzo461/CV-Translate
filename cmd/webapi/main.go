/*
webapi è il punto di ingresso del backend WASAText.

Costruisce il router descritto da doc/api.yaml e avvia il server HTTP, che
viene spento in modo pulito alla ricezione di SIGINT o SIGTERM. I dati sono
tenuti in memoria: l'entrypoint non apre né configura alcun database.

Utilizzo:

	go run ./cmd/webapi [--api-host :3000]
*/
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

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
	flag.Parse()

	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	logger.SetLevel(logrus.InfoLevel)
	logger.Info("avvio di WASAText")

	// Store in memoria: nessuna risorsa esterna da aprire o chiudere.
	router, err := api.New(api.Config{
		Logger:   logger,
		Database: database.NewInMemory(),
	})
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
