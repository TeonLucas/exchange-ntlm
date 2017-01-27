package main

import (
	"github.com/DavidSantia/exchange-ntlm/ews"
	"log"
)

func main() {
	client := &ews.EmailClient{
		UserPassword: "you@yourdomain.com:YourPassword",
		Url:          "https://example.com/EWS/Exchange.asmx",
	}

	// Set to 1 for debug, 2 for extra verbosity (like curl -v)
	client.Debug = 1

	// Check Inbox
	err := client.CheckInbox()
	if err != nil {
		log.Printf("Inbox result: %v\n", err)
	}

	// Show Id for each inbox item
	for i, item := range client.MessageList {
		log.Printf("Item %d: %#v\n", i+1, item)
	}
}
