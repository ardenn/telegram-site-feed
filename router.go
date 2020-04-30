package main

import "net/http"

func prepareRoutes() {
	http.HandleFunc("/update", loggingMiddlewareHandler(updateHandler))
	http.HandleFunc("/webhook", loggingMiddlewareHandler(webHookHandler))
	http.HandleFunc("/incoming", loggingMiddlewareHandler(websiteHandler))
}
