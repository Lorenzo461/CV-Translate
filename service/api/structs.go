package api

import (
	"regexp"

	"github.com/Lorenzo461/CV-Translate/service/database"
)

// Questo file contiene la traduzione in Go degli schemi dichiarati sotto
// `components/schemas` in doc/api.yaml. Per ogni schema ci sono:
//   - una struct con i tag `json` corrispondenti ai nomi delle proprietà;
//   - un metodo IsValid() che applica i vincoli (pattern, minLength,
//     maxLength, minimum, maximum) che lo YAML dichiara in modo dichiarativo
//     ma che Go non può far rispettare da solo.

// usernameRx corrisponde al `pattern` dello schema Username.
// I vincoli minLength/maxLength (3/16) sono già inclusi nel quantificatore.
var usernameRx = regexp.MustCompile(`^[a-zA-Z0-9_.]{3,16}$`)

// Limiti dello schema Identifier (minimum: 1, maximum: 9223372036854775807).
const (
	minIdentifier uint64 = 1
	maxIdentifier uint64 = 1<<63 - 1
)

// User corrisponde allo schema `User`.
// `photoUrl` non compare fra i campi `required`, quindi viene omesso dal JSON
// quando è vuoto.
type User struct {
	Identifier uint64 `json:"identifier"`
	Name       string `json:"name"`
	PhotoURL   string `json:"photoUrl,omitempty"`
}

// IsValid applica i vincoli degli schemi Identifier e Username.
func (u User) IsValid() bool {
	return isValidIdentifier(u.Identifier) && usernameRx.MatchString(u.Name)
}

// LoginRequest corrisponde al requestBody di doLogin.
type LoginRequest struct {
	Name string `json:"name"`
}

// IsValid verifica il campo obbligatorio `name` contro lo schema Username.
func (r LoginRequest) IsValid() bool {
	return usernameRx.MatchString(r.Name)
}

// LoginResponse corrisponde alla risposta 201 di doLogin.
type LoginResponse struct {
	Identifier uint64 `json:"identifier"`
}

// ChangeUsernameRequest corrisponde al requestBody di setMyUserName.
type ChangeUsernameRequest struct {
	Name string `json:"name"`
}

// IsValid verifica il campo obbligatorio `name` contro lo schema Username.
func (r ChangeUsernameRequest) IsValid() bool {
	return usernameRx.MatchString(r.Name)
}

// Error corrisponde allo schema `Error`, usato da tutte le risposte di errore.
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// isValidIdentifier applica i vincoli minimum/maximum dello schema Identifier.
func isValidIdentifier(id uint64) bool {
	return id >= minIdentifier && id <= maxIdentifier
}

// fromDatabaseUser converte il record di persistenza nella rappresentazione
// esposta dall'API. La conversione esiste per tenere separati i due modelli.
func fromDatabaseUser(u database.User) User {
	return User{
		Identifier: u.ID,
		Name:       u.Username,
		PhotoURL:   u.PhotoURL,
	}
}
