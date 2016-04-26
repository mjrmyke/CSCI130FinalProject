package main

import (
	"html/template"
	"log"
	"net/http"
)

func serve(res http.ResponseWriter, req *http.Request) {
	temp, err := template.ParseFiles("index.html")
	if err != nil {
	}
	temp.Execute(res, nil)
}

func main() {
	http.HandleFunc("/", serve)
	err := http.ListenAndServe(":8080", http.FileServer(http.Dir(".")))
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
