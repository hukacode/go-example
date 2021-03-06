package main

import (
	"log"
	"net/http"
	"text/template"
)

var name string = "World"

type User struct {
	Name   string
	Skills []string
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func handleHello(writer http.ResponseWriter, request *http.Request) {
	hello, err := template.ParseFiles("hello.html")
	check(err)

	user := User{Name: name, Skills: []string{"Java", "Go", "SQL"}}
	err = hello.Execute(writer, user)
	check(err)
}

func handleForm(writer http.ResponseWriter, request *http.Request) {
	form, err := template.ParseFiles("form.html")
	check(err)
	err = form.Execute(writer, nil)
	check(err)
}

func handleAdd(writer http.ResponseWriter, request *http.Request) {
	nameForm := request.FormValue("name")
	name = nameForm

	http.Redirect(writer, request, "hello", http.StatusFound)
}

func main() {
	http.HandleFunc("/hello", handleHello)
	http.HandleFunc("/form", handleForm)
	http.HandleFunc("/add", handleAdd)
	err := http.ListenAndServe("localhost:8080", nil)
	log.Fatal(err)
}
