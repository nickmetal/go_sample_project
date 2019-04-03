package handlers

import (
	"fmt"
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

type IndexPageData struct {
	StaticRoute string
}

func IndexPageHandler(w http.ResponseWriter, r *http.Request) {
	indexTemplatePath := filepath.Join(TEMPLATES_DIR_PATH, INDEX_HTML_PATH)
	_, err := os.Stat(TEMPLATES_DIR_PATH)
	templateData := IndexPageData{STATIC_DIR_PATH}

	if err != nil {
		if os.IsNotExist(err) {
			log.Println("TEMPLATES_DIR_PATH path not found in pwd")
			http.NotFound(w, r)
			return
		}
	}

	tmpl, err := template.ParseFiles(indexTemplatePath, indexTemplatePath)
	if err != nil {
		log.Printf("500: unable to parse indexTemplatePath(%s): %v", indexTemplatePath, err.Error())
		// Return a generic "Internal Server Error" message
		http.Error(w, http.StatusText(500), 500)
		return
	}

	if r.URL.Path == "/" {
		if err := tmpl.ExecuteTemplate(w, "index_page", templateData); err != nil {
			log.Println("template render was failed", err.Error())
			http.Error(w, http.StatusText(500), 500)
		}
	} else {
		errorMessage := fmt.Sprintf("404: No pages found for: %s path", r.URL.Path)
		log.Println(errorMessage)
		http.Error(w, errorMessage, 404)
	}

}
