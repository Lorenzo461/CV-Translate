/*
webapi è il punto di ingresso del backend WASAText.

Usa esclusivamente la libreria standard di Go e non apre alcun database: si
limita a costruire il router del package service/api e ad avviare il server
HTTP, che si spegne in modo pulito alla ricezione di SIGINT (Ctrl+C) o
SIGTERM.

Utilizzo:

	go run ./cmd/webapi [--addr :3000]
*/
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Lorenzo461/CV-Translate/service/api"
)

func main() {
	if err := run(); err != nil {
		log.Printf("errore fatale: %v", err)
		os.Exit(1)
	}
}

func run() error {
	addr := flag.String("addr", ":3000", "indirizzo di ascolto del server HTTP")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.LstdFlags)

	// Le rotte vivono nel package service/api: main non le conosce, si limita
	// a costruire il router e a passarne l'handler al server.
	router := api.New(api.Config{Logger: logger})

	server := &http.Server{
		Addr:              *addr,
		Handler:           logRequests(logger, router.Handler()),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
	}

	// Il server gira in una goroutine, così main resta libero di attendere
	// il segnale di terminazione.
	serverErrors := make(chan error, 1)
	go func() {
		logger.Printf("server in ascolto su %s", browserURL(*addr))
		serverErrors <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("errore nell'avvio del server: %w", err)
		}
	case sig := <-shutdown:
		logger.Printf("ricevuto il segnale %s: spegnimento in corso", sig)

		// Diamo alle richieste in corso 10 secondi per concludersi.
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			_ = server.Close()
			return fmt.Errorf("errore nello spegnimento del server: %w", err)
		}
	}

	logger.Print("arrivederci")
	return nil
}

// logRequests stampa una riga per ogni richiesta ricevuta: durante lo sviluppo
// dice subito se il client sta davvero arrivando al server.
func logRequests(logger *log.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("%s %s da %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

// browserURL trasforma l'indirizzo di ascolto in un URL apribile nel browser:
// ":3000" diventa "http://localhost:3000", mentre un indirizzo che indica già
// un host viene usato così com'è.
func browserURL(addr string) string {
	if addr == "" {
		return "http://localhost:80"
	}
	if strings.HasPrefix(addr, ":") {
		return "http://localhost" + addr
	}
	return "http://" + addr
}
