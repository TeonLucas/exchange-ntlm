package ews

import (
	"encoding/xml"
)

// Namespace s:

type Envelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    Body
}

type Body struct {
	XMLName xml.Name `xml:"Body"`
	Raw     []byte   `xml:",innerxml"`
}

// Namespace m

type RootFolder struct {
	XMLName      xml.Name `xml:"RootFolder"`
	TotalItems   int      `xml:"TotalItemsInView,attr"`
	IncludesLast bool     `xml:"IncludesLastItemInRange,attr"`
	Raw          []byte   `xml:",innerxml"`
}

// Namespace t

type FolderId struct {
	Id        string `xml:"Id,attr"`
	ChangeKey string `xml:"ChangeKey,attr"`
}

type ItemId struct {
	Id        string `xml:"Id,attr"`
	ChangeKey string `xml:"ChangeKey,attr"`
}

type AttachmentId struct {
	Id string `xml:"Id,attr"`
}

type FileAttachment struct {
	XMLName      xml.Name `xml:"FileAttachment"`
	AttachmentId AttachmentId
	Name         string `xml:"Name"`
	Content      string `xml:"Content"`
	ContentType  string `xml:"ContentType"`
}
