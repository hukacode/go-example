package main

import "net/http"

func (app *application) routes() *http.NewServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/todo", app.read)
	mux.HandleFunc("/todo/create", app.create)
	mux.HandleFunc("/todo/update", app.update)

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
