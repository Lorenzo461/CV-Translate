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
	"strings"
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

	// I pattern sono solo percorsi, senza il metodo HTTP davanti: quella forma
	// richiede Go 1.22 e su versioni precedenti verrebbe interpretata come un
	// nome di host, facendo rispondere 404 a ogni richiesta. Il metodo lo
	// controlliamo dentro l'handler, così il server funziona con qualsiasi
	// versione di Go.
	mux := http.NewServeMux()
	mux.HandleFunc("/liveness", handleLiveness)
	mux.HandleFunc("/", handleRoot)

	server := &http.Server{
		Addr:              *addr,
		Handler:           logRequests(mux),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
	}

	// Il server gira in una goroutine, così main resta libero di attendere
	// il segnale di terminazione.
	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("server in ascolto su %s", browserURL(*addr))
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

// handleRoot risponde sulla radice del sito.
func handleRoot(w http.ResponseWriter, r *http.Request) {
	// Il pattern "/" raccoglie ogni percorso non registrato altrove:
	// distinguiamo la radice dalle rotte inesistenti.
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "metodo non consentito", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = fmt.Fprintln(w, "WASAText: server attivo")
}

// handleLiveness segnala che il server è in grado di rispondere.
func handleLiveness(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "metodo non consentito", http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// logRequests stampa una riga per ogni richiesta ricevuta: durante lo sviluppo
// dice subito se il browser sta davvero arrivando al server.
func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s da %s", r.Method, r.URL.Path, r.RemoteAddr)
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
