package database

import (
	"sync"
)

// inMemoryDB è un'implementazione di AppDatabase che tiene i dati in memoria.
// I dati vivono quanto il processo: serve per sviluppo e test, non in
// produzione. L'implementazione persistente su SQLite resta in database.go e
// users.go, e si può riattivare sostituendo la costruzione dello store in
// cmd/webapi/main.go.
type inMemoryDB struct {
	mu     sync.RWMutex
	users  map[uint64]User   // id -> utente
	byName map[string]uint64 // username -> id, garantisce l'unicità
	nextID uint64
}

// NewInMemory costruisce un AppDatabase che non richiede alcuna risorsa
// esterna: nessun file, nessuna connessione, nessuno schema da creare.
func NewInMemory() AppDatabase {
	return &inMemoryDB{
		users:  make(map[uint64]User),
		byName: make(map[string]uint64),
		nextID: 1,
	}
}

// DoLogin restituisce l'identificativo dell'utente con il nome indicato,
// creandolo se non esiste ancora.
func (db *inMemoryDB) DoLogin(name string) (uint64, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if id, ok := db.byName[name]; ok {
		return id, nil
	}

	id := db.nextID
	db.nextID++
	db.users[id] = User{ID: id, Username: name}
	db.byName[name] = id

	return id, nil
}

// SetMyUserName aggiorna il nome utente e restituisce il profilo aggiornato.
func (db *inMemoryDB) SetMyUserName(userID uint64, name string) (User, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	user, ok := db.users[userID]
	if !ok {
		return User{}, ErrUserNotFound
	}

	// Il nome è già dell'utente stesso: operazione idempotente.
	if owner, taken := db.byName[name]; taken && owner != userID {
		return User{}, ErrUsernameTaken
	}

	delete(db.byName, user.Username)
	user.Username = name
	db.users[userID] = user
	db.byName[name] = userID

	return user, nil
}

// GetUserByID restituisce il profilo dell'utente indicato.
func (db *inMemoryDB) GetUserByID(userID uint64) (User, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	user, ok := db.users[userID]
	if !ok {
		return User{}, ErrUserNotFound
	}
	return user, nil
}

// Ping è sempre soddisfatto: non esiste alcuna risorsa esterna da verificare.
func (db *inMemoryDB) Ping() error {
	return nil
}
