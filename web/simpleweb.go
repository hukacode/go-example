package main

import (
	"log"
	"net/http"
)

func handleHello(writer http.ResponseWriter, request *http.Request) {
	hello := []byte("Hello web")
	_, err := writer.Write(hello)
	if err != nil {
		log.Fatal(err)
	}
}
func main() {
	http.HandleFunc("/hello", handleHello)
	err := http.ListenAndServe("localhost:8080", nil)
	log.Fatal(err)
}
