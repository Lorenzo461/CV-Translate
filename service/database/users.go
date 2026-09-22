package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/mattn/go-sqlite3"
)

// DoLogin implementa l'operazione doLogin: restituisce l'identificativo
// dell'utente con il nome indicato e, se non esiste, lo crea.
func (db *appdbimpl) DoLogin(name string) (uint64, error) {
	var id uint64

	err := db.c.QueryRow(`SELECT id FROM users WHERE username = ?`, name).Scan(&id)
	switch {
	case err == nil:
		// L'utente esiste già: login senza creazione.
		return id, nil
	case !errors.Is(err, sql.ErrNoRows):
		return 0, fmt.Errorf("error looking up user %q: %w", name, err)
	}

	res, err := db.c.Exec(`INSERT INTO users (username) VALUES (?)`, name)
	if err != nil {
		// Due login simultanei con lo stesso nome: il secondo viola il
		// vincolo UNIQUE, quindi rileggiamo l'utente creato dall'altro.
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.Code == sqlite3.ErrConstraint {
			if err := db.c.QueryRow(`SELECT id FROM users WHERE username = ?`, name).Scan(&id); err != nil {
				return 0, fmt.Errorf("error looking up user %q after conflict: %w", name, err)
			}
			return id, nil
		}
		return 0, fmt.Errorf("error creating user %q: %w", name, err)
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error reading the identifier of the new user: %w", err)
	}
	return uint64(lastID), nil
}

// SetMyUserName implementa l'operazione setMyUserName.
func (db *appdbimpl) SetMyUserName(userID uint64, name string) (User, error) {
	res, err := db.c.Exec(`UPDATE users SET username = ? WHERE id = ?`, name, userID)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.Code == sqlite3.ErrConstraint {
			return User{}, ErrUsernameTaken
		}
		return User{}, fmt.Errorf("error updating the username of user %d: %w", userID, err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return User{}, fmt.Errorf("error checking the username update: %w", err)
	}
	if affected == 0 {
		// Nessuna riga aggiornata: o l'utente non esiste, o il nome era già
		// quello richiesto (l'operazione è idempotente).
		return db.GetUserByID(userID)
	}

	return db.GetUserByID(userID)
}

// GetUserByID restituisce il profilo dell'utente indicato.
func (db *appdbimpl) GetUserByID(userID uint64) (User, error) {
	var u User
	err := db.c.QueryRow(
		`SELECT id, username, photo_url FROM users WHERE id = ?`, userID,
	).Scan(&u.ID, &u.Username, &u.PhotoURL)

	if errors.Is(err, sql.ErrNoRows) {
		return User{}, ErrUserNotFound
	} else if err != nil {
		return User{}, fmt.Errorf("error reading user %d: %w", userID, err)
	}
	return u, nil
}
