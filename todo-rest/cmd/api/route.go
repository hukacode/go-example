package main

import (
	"expvar"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodPost, "/v1/task", app.requirePermission("task:read", app.createTaskHandler))
	router.HandlerFunc(http.MethodGet, "/v1/task", app.requirePermission("task:write", app.readListTasksHandler))
	router.HandlerFunc(http.MethodGet, "/v1/task/:id", app.requirePermission("task:read", app.readTaskHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/task/:id", app.requirePermission("task:write", app.updateTaskHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/task/:id", app.requirePermission("task:write", app.deleteTaskHandler))

	router.HandlerFunc(http.MethodPost, "/v1/user", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/user/activated", app.activateUserHandler)

	router.HandlerFunc(http.MethodPost, "/v1/token/authentication", app.createAuthenticationTokenHandler)

	router.Handler(http.MethodGet, "/v1/metrics", expvar.Handler())

	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router)))))
}
