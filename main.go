package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi, Welcome home %s", r.URL.Path[1:])
}

func main() {
	apiKey := os.Getenv("API_KEY")
	botURL := "https://api.telegram.org/bot" + apiKey + "/"
	messageURL := botURL + "sendMessage"
	_ = messageURL
	http.HandleFunc("/", indexHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
