## exchange-ntlm

### Quickstart to use Microsoft Exchange from Go Language (without Windows)

Outlook web services is a SOAP based protocol that uses the NT LAN Manager (NTLM) suite of Microsoft security protocols for authentication, integrity, and confidentiality. The libraries to authenticate with NTLM come built-in with the Windows operating system Security Support Provider Interface (SSPI) api.  On Linux, you'll need a library to re-implement this protocol.

Although NTLM is used over HTTP, it employs HTTP status codes and headers in its own way, not strictly following HTTP conventions. Your best bet to get in to your Exchange server from Linux, of course, is to reconfigure that server to allow other authentication protocols, so you don't have to use NTLM.

But let's say you are programming in Go-language, your Exchange server is already locked down with NTLM, yet you want to get in.  Hopefully this project might speed your research to get going quicker. 
 
## What is NTLM

NTLM uses three messages to authenticate a client in a connection oriented environment (connectionless is similar), and a fourth additional message if integrity is desired.

The client begins with a simple GET request to the NTLM server.  It responds with 401 Unauthorized, and reveals its authentication scheme in the response header:
* WWW-Authenticate: NTLM

#### The Handshake

Next you need to work a three-part handshake as follows:
1. First, the client establishes a network path to the server and sends a NEGOTIATE_MESSAGE advertising its capabilities.
1. Next, the server responds with CHALLENGE_MESSAGE which is used to establish the identity of the client.
1. Finally, the client responds to the challenge with an AUTHENTICATE_MESSAGE.

To understand the handshake in more detail, I found this write-up quite good: [NTLM Authentication Scheme for HTTP](https://www.innovation.ch/personal/ronald/ntlm.html)

[Figure 1](https://raw.githubusercontent.com/DavidSantia/exchange-ntlm/master/README-figure1.png)

#### Keeping the connection alive
NTLM authenticates connections, not requests. In Go-language, this means you are using the RoundTripper, i.e.:
* https://golang.org/pkg/net/http/#Transport

and then implementing the handshake natively.

When you begin your handshake, the network connection must be kept alive, i.e. between the receiving of the type-2 message from the server (step 2b) and the sending of the type-3 message (step 3). Each time the connection is closed this second part (steps 2 and 3) must be repeated over the new connection (i.e. it's not enough to just keep sending the last type-3 message).

Once the connection is authenticated, the Authorization header need not be sent anymore while the connection stays open, no matter what resource is accessed.

### One go libary to consider

This library is referred to a lot and has much of the parts for a native go implementation:
* https://github.com/ThomsonReutersEikon/go-ntlm

I tried it with an Exchange server, hacked around a bit, but didn't get it to work. I then read the fine print: "The major missing piece is the negotiation of capabilities between the client and the server, for our use we hardcoded a supported set of negotiation flags."  I then realized I needed a working start point.

#### Curl to the rescue

The curl program has a cool command-line option: -ntlm
Here is an example:
```sh
curl -v --ntlm \
  -u 'you@yourdomain.com:YourPassword' \
  -H 'Content-Type: text/xml;charset=UTF-8' \
  -d @finditem.xml \
  https://example.com/EWS/Exchange.asmx
```

Once have a way to do the handshake, like above, you post SOAP commands (which are just XML) and get back XML responses.

The most basic command check your inbox, which uses FindItem, see [finditem.xml](https://raw.githubusercontent.com/DavidSantia/exchange-ntlm/master/finditem.xml) from this repo.

Microsoft has all the XML you need to post nicely documented; below are some examples:

* [FindItem link](https://msdn.microsoft.com/en-us/library/office/aa566107(v=exchg.150).aspx#sectionSection1)
* [GetItem link](https://msdn.microsoft.com/en-us/library/office/aa566107(v=exchg.150).aspx#sectionSection1)
* [FindFolder link](https://msdn.microsoft.com/en-us/library/office/dd633627(v=exchg.80).aspx#Anchor_0)

#### Libcurl in Go

Once I had everything working with curl, I kept going. It turns out you can actually just use Libcurl from Go:
* https://github.com/andelf/go-curl

Using a C-library means you cannot use the CGO_ENABLE=0 when you compile.  This means if you are trying to build Go for an empty container that doesn't have a full OS, you won't have all the shared libraries libcurl uses, and it may take some work to build libcurl.a to go the static route.

But for the general case, and definately to jump-start your NTLM authentication, this is a really good route.

The sample Go code in this repo does the same thing as the curl command.
* It uses the go-curl library, Copyright 2014 Shuyu Wang (<andelf@gmail.com>)
* You will even notice you can turn on debug in go-curl (just like curl -v), and watch the handshake steps proceed.
* This makes it super easy to see if your go program's handshake matches your curl command.

To run:
```sh
$ go get https://github.com/andelf/go-curl
$ go build
$ ./exchange-ntlm
