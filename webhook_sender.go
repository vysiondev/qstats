package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type WebhookMessage struct {
	Content string `json:"content"`
}

func (b *BaseHandler) SendMessageToWebhook(m string) {
	if len(b.Config.Webhook.URL) > 0 {
		strJson, err := json.Marshal(WebhookMessage{Content: m})
		if err != nil {
			fmt.Println("Could not marshal webhook JSON")
			return
		}
		resp, err := http.Post(b.Config.Webhook.URL, "application/json", bytes.NewBuffer(strJson))
		if err != nil {
			fmt.Println("Problem sending message to webhook URL: " + err.Error())
			return
		}
		closeErr := resp.Body.Close()
		if closeErr != nil {
			fmt.Println("Cound not close response body: " + closeErr.Error())
		}
	}
}
