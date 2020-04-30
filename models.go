package main

// SetWebHookData defines the request payload required to set s webhook
type SetWebHookData struct {
	URL string `json:"url"`
}

// Chat is a Telegram Chat Object
type Chat struct {
	ChatID int `json:"id"`
}

// Message is Telegram Message Object
type Message struct {
	MessageID int    `json:"message_id"`
	TChat     Chat   `json:"chat"`
	Text      string `json:"text"`
}

// Update is a Telegram Update Object
type Update struct {
	UpdateID int     `json:"update_id"`
	TMessage Message `json:"message"`
}

// TMessagePayload defines the request payload to sendMessage on Telegram
type TMessagePayload struct {
	ChatID    int    `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

// WebsiteForm defines the data from the website form
type WebsiteForm struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Body  string `json:"message"`
	IP    string `json:"ip"`
}

// NetlifyData describes the structure of the payload from Netlify's webhook
type NetlifyData struct {
	Data WebsiteForm `json:"data"`
}
