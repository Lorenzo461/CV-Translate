package api

import "regexp"

// Questo file contiene la traduzione in Go degli schemi dichiarati sotto
// `components/schemas` in doc/api.yaml. Per ogni schema ci sono:
//   - una struct con i tag `json` corrispondenti ai nomi delle proprietà;
//   - un metodo IsValid() che applica i vincoli (pattern, minLength,
//     maxLength, minimum, maximum) che lo YAML dichiara in modo dichiarativo
//     ma che Go non può far rispettare da solo.

// usernameRx corrisponde al `pattern` dello schema Username.
// I vincoli minLength/maxLength (3/16) sono già inclusi nel quantificatore.
var usernameRx = regexp.MustCompile(`^[a-zA-Z0-9_.]{3,16}$`)

// maxRequestBody limita la dimensione dei corpi JSON accettati: senza questo
// limite un client potrebbe far crescere la memoria del server a piacere.
const maxRequestBody = 1 << 20 // 1 MiB

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

// Error corrisponde allo schema `Error`, usato da tutte le risposte di errore.
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
