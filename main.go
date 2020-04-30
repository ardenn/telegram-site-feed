package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
)

var apiKey string
var botURL string
var successPayload string
var chatInt int64

func init() {
	apiKey = os.Getenv("API_KEY")
	chatInt, _ = strconv.ParseInt(os.Getenv("CHAT_ID"), 10, 64)
	botURL = "https://api.telegram.org/bot" + apiKey + "/"
	successResponse := make(map[string]string)
	successResponse["message"] = "success"
	payload, _ := json.Marshal(successResponse)
	successPayload = string(payload)
}

func main() {
	log.Println("Server started on port 8080 ...")
	prepareRoutes()
	http.ListenAndServe(":8080", nil)
}
