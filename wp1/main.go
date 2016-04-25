package main

import(
	"html/template"
	"log"
	"net/http"
)

func serve(res http.ResponseWriter, req *http.Request) {
	temp := template.New("index.html")
	temp, err := template.ParseFiles("index.html")
	if err != nil {
		log.Fatalln(err,"Failed!")
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

