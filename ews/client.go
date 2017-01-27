package ews

// State that holds data between different functions (methods)
type EmailClient struct {
	Debug        int
	UserPassword string
	Url          string
	PostXML      string
	RespXML      string
	MessageList  []ItemId
}

// To use with posts that filter for certain emails
type Filter struct {
	Class       string
	Subject     string
	SenderName  string
	SenderEmail string
}
