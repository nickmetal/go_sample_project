package main

import (
	"fmt"
	"githome/yalatask/handlers"
	"log"
	"net/http"
	"sync"
)

const (
	STATIC_DIR_PATH = "/static/"
)

// log to stderr request in [%s] - %s format: 2019/01/01 00:03:17 [POST] - /path/
func LogRequest(methdoName, path string) {
	log.Printf("[%s] - %s\n", methdoName, path)
}

// create a goroutine for each http request
func Middleware(handler func(resp http.ResponseWriter, req *http.Request)) func(resp http.ResponseWriter, req *http.Request) {
	return func(resp http.ResponseWriter, req *http.Request) {
		defer Write500ErrorOnPanic(resp, req)
		LogRequest(req.Method, req.URL.Path)

		var wg sync.WaitGroup
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			handler(resp, req)
		}(&wg)
		wg.Wait()
	}
}

// catch panic and write error as 500 http error
func Write500ErrorOnPanic(resp http.ResponseWriter, req *http.Request) {
	if err := recover(); err != nil {
		log.Printf("got a panic in %s: %v\n", req.URL.Path, err)
		http.Error(resp, fmt.Sprintf("Unhandled Error: %v", err), 500)
	}
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle(STATIC_DIR_PATH, http.StripPrefix(STATIC_DIR_PATH, fs))
	http.HandleFunc("/", Middleware(handlers.IndexPageHandler))
	http.HandleFunc("/transport-issue/", Middleware(handlers.TransportIssueHandler))

	log.Println("Listening...")
	http.ListenAndServe(":3000", nil)
}
