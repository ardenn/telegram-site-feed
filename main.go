package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

type webHookData struct {
	URL string `json:"url"`
}

type chat struct {
	ChatID int `json:"id"`
}
type message struct {
	MessageID   int    `json:"message_id"`
	MessageChat chat   `json:"chat"`
	Text        string `json:"text"`
}
type update struct {
	UpdateID int     `json:"update_id"`
	TMessage message `json:"message"`
}

type requestPayload struct {
	ChatID    int    `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

type websiteForm struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Body  string `json:"message"`
	IP    string `json:"ip"`
}

type netlifyData struct {
	Data websiteForm `json:"data"`
}

func sendMessage(chatID int, text string) {
	requestData := requestPayload{ChatID: chatID, Text: text, ParseMode: "Markdown"}
	payload, _ := json.Marshal(requestData)
	_, _ = http.Post(botURL+"sendMessage", "application/json", bytes.NewBuffer(payload))
}

func buildMessage(data websiteForm) (text string) {
	text = fmt.Sprintf("*%s*\n*Name*: %s\n*Email*: %s\n%s.\n - sent from *%s*", "New Website Contact", data.Name, data.Email, data.Body, data.IP)
	return
}

// Handler for website callback (Receives messages from the website form)
func websiteHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var data netlifyData
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error processing request")
			return
		}
		err = json.Unmarshal(body, &data)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error decoding request body")
			return
		}
		textBody := buildMessage(data.Data)
		sendMessage(int(chatInt), textBody)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, string(successPayload))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Method not Allowed")
	}
}

// Handler for updates sent to telegram. Receives messages sent to bot
func updateHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var data update
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error processing request")
			return
		}
		err = json.Unmarshal(body, &data)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error decoding request body")
			return
		}
		sendMessage(data.TMessage.MessageChat.ChatID, data.TMessage.Text)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, string(successPayload))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Method not Allowed")
	}
}

// Performs GET and POST to view current webhook details set webhook endpoint
func webHook(w http.ResponseWriter, r *http.Request) {
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
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error processing request")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, string(body))
	case "POST":
		var requestData webHookData
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error processing request")
			return
		}
		if err = json.Unmarshal(body, &requestData); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error decoding body json!")
			return
		}
		if requestData.URL == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "url is required!")
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
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error processing request")
		}
		if resp.StatusCode >= 400 {
			log.Println("Error occurred", respBody)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, string(successPayload))

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Method not Allowed")
	}
}
func main() {
	http.HandleFunc("/update", updateHandler)
	http.HandleFunc("/webhook", webHook)
	http.HandleFunc("/incoming", websiteHandler)
	log.Println("Server started on port 8080 ...")
	http.ListenAndServe(":8080", nil)
}
