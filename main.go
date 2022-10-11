package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func main() {
	http.HandleFunc("/", serveTemplate)

	log.Print("Listening on :3000...")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	fp := filepath.Join("static", "index.html")
	tmpl, _ := template.ParseFiles(fp)
	tmpl.Execute(w, nil)
}
