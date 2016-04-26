package main

import (
	"html/template"
	"log"
	"net/http"
)

func serve(res http.ResponseWriter, req *http.Request) {
	//create template from index.html file
	temp, err := template.ParseFiles("index.html")
	if err != nil {
		log.Fatalln(err, "Failed Parsing File")
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
