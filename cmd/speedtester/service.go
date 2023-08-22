package main

import (
	"log"
	"net/http"
)

func main() {
	s := &http.Server{
		Addr: "localhost:8080",
	}
	http.HandleFunc("/download", downloadHandler())
	http.HandleFunc("/upload", uploadHandler())

	log.Fatal(s.ListenAndServe())
}
