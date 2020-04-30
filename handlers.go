package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func sendMessage(chatID int, text string) {
	requestData := TMessagePayload{ChatID: chatID, Text: text, ParseMode: "Markdown"}
	payload, _ := json.Marshal(requestData)
	_, _ = http.Post(botURL+"sendMessage", "application/json", bytes.NewBuffer(payload))
}

func buildMessage(data WebsiteForm) (text string) {
	text = fmt.Sprintf("*%s*\n*Name*: %s\n*Email*: %s\n%s.\n - sent from *%s*", "New Website Contact", data.Name, data.Email, data.Body, data.IP)
	return
}

// Handler for website callback (Receives messages from the website form)
func websiteHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var data NetlifyData
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			http.Error(w, "Error processing request", http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal(body, &data)
		if err != nil {
			log.Println(err)
			http.Error(w, "Error decoding request body", http.StatusInternalServerError)
			return
		}
		textBody := buildMessage(data.Data)
		sendMessage(int(chatInt), textBody)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, string(successPayload))
	default:
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
	}
}

// Handler for updates sent to telegram. Receives messages sent to bot
func updateHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var data Update
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			http.Error(w, "Error processing request", http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal(body, &data)
		if err != nil {
			log.Println(err)
			http.Error(w, "Error decoding request body", http.StatusInternalServerError)
			return
		}
		sendMessage(data.TMessage.TChat.ChatID, data.TMessage.Text)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, string(successPayload))
	default:
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
	}
}

// Performs GET and POST to view current webhook details set webhook endpoint
func webHookHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		resp, err := http.Get(botURL + "getWebhookInfo")
		if err != nil {
			log.Println("Error reaching Telegram", err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			http.Error(w, "Error processing request", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, string(body))
	case "POST":
		var requestData SetWebHookData
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error processing request")
			return
		}
		if err = json.Unmarshal(body, &requestData); err != nil {
			log.Println(err)
			http.Error(w, "Error decoding request body", http.StatusInternalServerError)
			return
		}
		if requestData.URL == "" {
			http.Error(w, "url is required", http.StatusBadRequest)
			return
		}
		resp, err := http.Post(botURL+"setWebhook", "application/json", bytes.NewBuffer(body))
		defer resp.Body.Close()
		if err != nil {
			log.Println("Error reaching Telegram")
		}
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			http.Error(w, "Error processing request", http.StatusInternalServerError)
			return
		}
		if resp.StatusCode >= 400 {
			log.Println("Error occurred", respBody)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, string(successPayload))

	default:
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
	}
}
