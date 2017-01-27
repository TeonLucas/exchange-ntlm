package ews

import (
	"encoding/xml"
)

// SOAP request for getting list of inbox folder items
// From https://msdn.microsoft.com/en-us/library/office/aa566107(v=exchg.150).aspx#sectionSection1

const postFindItem = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"
               xmlns:t="http://schemas.microsoft.com/exchange/services/2006/types">
  <soap:Body>
    <FindItem xmlns="http://schemas.microsoft.com/exchange/services/2006/messages"
               xmlns:t="http://schemas.microsoft.com/exchange/services/2006/types"
              Traversal="Shallow">
      <ItemShape>
        <t:BaseShape>IdOnly</t:BaseShape>
      </ItemShape>
      <IndexedPageItemView MaxEntriesReturned="100" BasePoint="Beginning" Offset="%d" />
      <ParentFolderIds>
        <t:DistinguishedFolderId Id="inbox"/>
      </ParentFolderIds>
    </FindItem>
  </soap:Body>
</soap:Envelope>`

// Namespace m

type FindItemResp struct {
	XMLName xml.Name        `xml:"FindItemResponse"`
	Msg     FindItemRespMsg `xml:"ResponseMessages>FindItemResponseMessage"`
}

type FindItemRespMsg struct {
	XMLName xml.Name `xml:"FindItemResponseMessage"`
	Class   string   `xml:"ResponseClass,attr"`
	Code    string   `xml:"ResponseCode"`
	Folder  RootFolder
}

// Namespace t

type Items struct {
	XMLName xml.Name `xml:"Items"`
	Message []Message
}

type Message struct {
	XMLName xml.Name `xml:"Message"`
	ItemId  ItemId
}
