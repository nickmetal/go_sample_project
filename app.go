package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const (
	INDEX_HTML_PATH    = "index.html"
	TEMPLATES_DIR_PATH = "templates"
	STATIC_DIR_PATH    = "/static/"
)

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle(STATIC_DIR_PATH, http.StripPrefix(STATIC_DIR_PATH, fs))
	http.HandleFunc("/", serveTemplate)
	http.HandleFunc("/api/", serveApi)

	log.Println("Listening...")
	http.ListenAndServe(":3000", nil)
}

func serveApi(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("zbs"))
}

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	lp := filepath.Join(TEMPLATES_DIR_PATH, INDEX_HTML_PATH)
	fp := filepath.Join(TEMPLATES_DIR_PATH, filepath.Clean(r.URL.Path))

	// Return a 404 if the template doesn't exist
	info, err := os.Stat(fp)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
	}

	// Return a 404 if the request is for a directory
	if info.IsDir() {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles(lp, fp)
	if err != nil {
		// Log the detailed error
		log.Println(err.Error())
		// Return a generic "Internal Server Error" message
		http.Error(w, http.StatusText(500), 500)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "index_page", nil); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}
