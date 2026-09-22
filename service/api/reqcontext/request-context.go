/*
Package reqcontext contiene il contesto della singola richiesta HTTP.

Il contesto viene creato dal wrapper di routing (service/api/wrap.go) prima di
invocare l'handler e trasporta le informazioni che derivano dalla richiesta e
non dai suoi parametri: l'identificativo univoco della richiesta, il logger
già arricchito con tale identificativo e, per le operazioni autenticate,
l'identificativo dell'utente estratto dal token (securityScheme `bearerAuth`
della specifica OpenAPI).
*/
package reqcontext

import (
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
)

// RequestContext è il contesto della richiesta HTTP in corso.
type RequestContext struct {
	// ReqUUID è l'identificativo univoco della richiesta, usato per correlare
	// le righe di log appartenenti alla stessa chiamata.
	ReqUUID uuid.UUID

	// Logger è il logger da usare all'interno dell'handler: contiene già il
	// campo `reqid`.
	Logger logrus.FieldLogger

	// UserID è l'identificativo dell'utente autenticato. Vale 0 per le
	// operazioni che nella specifica OpenAPI dichiarano `security: []`
	// (cioè doLogin).
	UserID uint64
}
