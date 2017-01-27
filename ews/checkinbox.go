package ews

import (
	"encoding/xml"
	"errors"
	"fmt"
	"log"
)

func (client *EmailClient) CheckInbox() (err error) {

	var err_txt string
	var offset, retrieved int

	offset = 0
	log.Printf("Reading Inbox at %s\n", client.Url)
	for {
		client.PostXML = fmt.Sprintf(postFindItem, offset)
		err = client.ewsPost()
		if err != nil {
			return
		}

		// One Unmarshall per namespace

		// xmlns:s
		envl := new(Envelope)
		xml.Unmarshal([]byte(client.RespXML), envl)

		if envl.Body.XMLName.Local != "Body" {
			err_txt = fmt.Sprintf("Error parsing SOAP body: invalid XML syntax")
			err = errors.New(err_txt)
			return
		}

		// xmlns:m
		resp := new(FindItemResp)
		xml.Unmarshal(envl.Body.Raw, resp)

		if resp.Msg.Class != "Success" {
			err_txt = fmt.Sprintf("Error finding items in Inbox: Class=%s Code=%s",
				resp.Msg.Class, resp.Msg.Code)
			err = errors.New(err_txt)
			return
		}

		if resp.Msg.Folder.XMLName.Local != "RootFolder" {
			err_txt = fmt.Sprintf("Error parsing SOAP Inbox root folder: invalid XML syntax")
			err = errors.New(err_txt)
			return
		}

		// xmlns:t
		items := new(Items)
		xml.Unmarshal(resp.Msg.Folder.Raw, items)

		if items.XMLName.Local != "Items" {
			err_txt = fmt.Sprintf("Error parsing SOAP Inbox Items: invalid XML syntax")
			err = errors.New(err_txt)
			return
		}

		if offset == 0 {
			log.Printf("Inbox contains %d items\n", resp.Msg.Folder.TotalItems)
		}
		retrieved = len(items.Message)
		offset += retrieved
		if client.Debug > 0 {
			log.Printf("DEBUG: Folder Items=%d, Retrieved=%d, IncludesLast=%t\n",
				resp.Msg.Folder.TotalItems, retrieved, resp.Msg.Folder.IncludesLast)
		}

		// Validate messages
		for i := 0; i < retrieved; i++ {
			if items.Message[i].XMLName.Local != "Message" {
				log.Printf("Error parsing SOAP Message #%d from Inbox: invalid XML syntax", i+1)
				continue
			}
			// Copy ItemId to context
			client.MessageList = append(client.MessageList, items.Message[i].ItemId)
		}

		if resp.Msg.Folder.IncludesLast {
			log.Printf("Retrieved %d Id's\n", len(client.MessageList))
			return
		}
	}
}
