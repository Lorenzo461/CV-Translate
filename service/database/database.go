/*
Package database contiene lo strato di persistenza dell'applicazione.

L'interfaccia AppDatabase espone un metodo per ogni operazione applicativa
richiesta dagli handler HTTP; l'implementazione concreta (appdbimpl) usa
SQLite. Gli handler non conoscono SQL: dialogano solo con questa interfaccia e
con gli errori di dominio dichiarati qui sotto, che vengono poi tradotti nei
codici di stato HTTP previsti dalla specifica OpenAPI.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// Errori di dominio restituiti dal database. Gli handler li riconoscono con
// errors.Is e li traducono nelle risposte dichiarate in doc/api.yaml.
var (
	// ErrUserNotFound indica che l'utente richiesto non esiste.
	ErrUserNotFound = errors.New("user not found")

	// ErrUsernameTaken indica che il nome utente è già usato da un altro
	// utente (risposta 409 di setMyUserName).
	ErrUsernameTaken = errors.New("username already taken")
)

// User è la rappresentazione di un utente nello strato di persistenza.
// È volutamente distinta dalla struct User del package api: lo schema del
// database può evolvere senza modificare il contratto REST.
type User struct {
	ID       uint64
	Username string
	PhotoURL string
}

// AppDatabase è l'interfaccia di alto livello verso il database.
type AppDatabase interface {
	// DoLogin restituisce l'identificativo dell'utente con il nome indicato,
	// creandolo se non esiste ancora.
	DoLogin(name string) (uint64, error)

	// SetMyUserName aggiorna il nome utente e restituisce il profilo
	// aggiornato. Restituisce ErrUsernameTaken se il nome è già occupato e
	// ErrUserNotFound se l'utente non esiste.
	SetMyUserName(userID uint64, name string) (User, error)

	// GetUserByID restituisce il profilo dell'utente, oppure ErrUserNotFound.
	GetUserByID(userID uint64) (User, error)

	// Ping verifica che la connessione al database sia utilizzabile.
	Ping() error
}

// appdbimpl è l'implementazione SQLite di AppDatabase.
type appdbimpl struct {
	c *sql.DB
}

// New costruisce un AppDatabase a partire da una connessione già aperta,
// creando lo schema se il database è vuoto.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building an AppDatabase")
	}

	// Le foreign key in SQLite sono disattivate per impostazione predefinita.
	if _, err := db.Exec(`PRAGMA foreign_keys = ON`); err != nil {
		return nil, fmt.Errorf("error enabling foreign keys: %w", err)
	}

	if err := createSchema(db); err != nil {
		return nil, fmt.Errorf("error creating database structure: %w", err)
	}

	return &appdbimpl{c: db}, nil
}

// createSchema crea le tabelle mancanti.
func createSchema(db *sql.DB) error {
	const usersTable = `CREATE TABLE IF NOT EXISTS users (
		id        INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		username  TEXT    NOT NULL UNIQUE,
		photo_url TEXT    NOT NULL DEFAULT ''
	);`

	if _, err := db.Exec(usersTable); err != nil {
		return fmt.Errorf("error creating table users: %w", err)
	}
	return nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
