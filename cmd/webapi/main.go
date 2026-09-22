/*
webapi avvia un server HTTP locale.

Usa esclusivamente la libreria standard di Go: nessuna dipendenza esterna e
nessuna persistenza. Il server espone due rotte di servizio e si spegne in
modo pulito alla ricezione di SIGINT (Ctrl+C) o SIGTERM.

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
	"syscall"
	"time"
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

	// Le rotte sono registrate su un ServeMux della libreria standard.
	// Da Go 1.22 il pattern può indicare anche il metodo HTTP.
	mux := http.NewServeMux()

	mux.HandleFunc("GET /liveness", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		// ServeMux instrada su "GET /" qualsiasi percorso non registrato:
		// distinguiamo la home dalle rotte inesistenti.
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = fmt.Fprintln(w, "WASAText: server attivo")
	})

	server := &http.Server{
		Addr:              *addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
	}

	// Il server gira in una goroutine, così main resta libero di attendere
	// il segnale di terminazione.
	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("server in ascolto su %s", server.Addr)
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
		log.Printf("ricevuto il segnale %s: spegnimento in corso", sig)

		// Diamo alle richieste in corso 10 secondi per concludersi.
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			_ = server.Close()
			return fmt.Errorf("errore nello spegnimento del server: %w", err)
		}
	}

	log.Print("arrivederci")
	return nil
}
