package ews

import (
	"fmt"
	"log"
	"time"

	curl "github.com/andelf/go-curl"
)

// From curl.h c lang header
const CURLAUTH_NTLM = 0x08

var sent bool = false

func write_data(ptr []byte, userdata interface{}) bool {
	ch, ok := userdata.(chan string)
	if ok {
		ch <- string(ptr)
		return true // ok
	} else {
		log.Printf("ERROR!\n")
		return false
	}
	return false
}

func (client *EmailClient) ewsPost() (err error) {
	// init the curl session
	client.RespXML = ""
	easy := curl.EasyInit()
	defer easy.Cleanup()

	// Authentication
	easy.Setopt(curl.OPT_URL, client.Url)
	easy.Setopt(curl.OPT_USERPWD, client.UserPassword)
	easy.Setopt(curl.OPT_HTTPAUTH, CURLAUTH_NTLM)

	// Data to post
	easy.Setopt(curl.OPT_HTTPHEADER, []string{"Content-Type: text/xml;charset=utf-8"})
	easy.Setopt(curl.OPT_POST, true)
	easy.Setopt(curl.OPT_POSTFIELDS, client.PostXML)
	easy.Setopt(curl.OPT_POSTFIELDSIZE, len(client.PostXML))

	var reply string

	// Channel for reply
	ch := make(chan string)
	go func(ch chan string) {
		for {
			block := <-ch
			if client.Debug > 0 {
				log.Printf("Got block size=%d\n", len(block))
			}
			reply = reply + block
		}
	}(ch)

	// Write function for reply
	easy.Setopt(curl.OPT_WRITEFUNCTION, write_data)
	easy.Setopt(curl.OPT_WRITEDATA, ch)

	if client.Debug > 1 {
		easy.Setopt(curl.OPT_VERBOSE, true)
	}
	// Perform session
	if err = easy.Perform(); err != nil {
		return err
	}

	// Check HTML response code
	r, err := easy.Getinfo(curl.INFO_RESPONSE_CODE)
	responsecode := r.(int)
	if responsecode != 200 {
		return fmt.Errorf("Error HTML reponse %d", responsecode)
	}

	// This wait serializes the parallelism above
	// It's there to use later, but keep things simple to start
	time.Sleep(100000)

	client.RespXML = reply
	return nil
}
