package main

import (
	"log"
	"net/http"
)

func loggingMiddlewareHandler(handler http.HandlerFunc, args ...interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path, "\t", r.Proto)
		handler(w, r)
	}

}
