package main

import (
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := envelop{
		"statis": "OK",
		"info": map[string]string{
			"env": app.config.env,
		},
	}

	err := app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
