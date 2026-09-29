package api

import "sync"

// userStore conserva gli utenti in memoria: nessun database, nessun file.
// I dati vivono quanto il processo.
//
// Il server HTTP serve le richieste in goroutine concorrenti, quindi ogni
// accesso alle mappe è protetto da un mutex.
type userStore struct {
	mu     sync.Mutex
	byName map[string]uint64 // username -> identificativo, garantisce l'unicità
	nextID uint64
}

// newUserStore costruisce uno store vuoto.
func newUserStore() *userStore {
	return &userStore{
		byName: make(map[string]uint64),
		nextID: 1,
	}
}

// login restituisce l'identificativo dell'utente con il nome indicato,
// creandolo se non esiste ancora. Il secondo valore dice se l'utente è stato
// creato adesso: serve solo per il log.
func (s *userStore) login(name string) (uint64, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if id, ok := s.byName[name]; ok {
		return id, false
	}

	id := s.nextID
	s.nextID++
	s.byName[name] = id

	return id, true
}
