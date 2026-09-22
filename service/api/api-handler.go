package api

import "net/http"

// Handler registra le rotte dell'applicazione. L'elenco ricalca uno a uno la
// sezione `paths` di doc/api.yaml: metodo HTTP, percorso e operationId.
func (rt *_router) Handler() http.Handler {
	// POST /session -> doLogin (security: [], quindi senza autenticazione)
	rt.router.POST("/session", rt.wrap(rt.doLogin, false))

	// PUT /users/me/username -> setMyUserName
	rt.router.PUT("/users/me/username", rt.wrap(rt.setMyUserName, true))

	// Endpoint di servizio, non descritto dalla specifica.
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
