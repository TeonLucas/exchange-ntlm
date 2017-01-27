## exchange-ntlm

This is an example of using Go-language to access Microsoft Exchange WebServices (Outlook)

To run, clone or download this repo, and type:
```sh
$ go get https://github.com/andelf/go-curl
$ go build
$ ./exchange-ntlm
```

### Quickstart to use Microsoft Exchange via NTLM on Linux

The Microsoft Exchange WebServices is a SOAP-based api that uses the NT LAN Manager (NTLM) suite of Microsoft security protocols for authentication, integrity, and confidentiality. The libraries to authenticate with NTLM come built-in with the Windows operating system Security Support Provider Interface (SSPI) api. On Linux, you'll need a library to re-implement this protocol.

Although NTLM is used over HTTP, it employs HTTP status codes and headers in its own way, not strictly following HTTP conventions. The easiest way to access your Exchange server from Linux, of course, is to reconfigure that server to allow other authentication protocols, so you don't have to use NTLM.

But let's say you are programming in Go-language, your Exchange server is already locked down with NTLM, and you want to get in.  Hopefully this project can speed your research to get going quicker. 
 
## What is NTLM

NTLM uses three messages to authenticate a client in a connection oriented environment (connectionless is similar), and a fourth additional message if integrity is desired.

The client begins with a simple GET request to the NTLM server.  It responds with 401 Unauthorized, and reveals its authentication scheme in the response header:
* WWW-Authenticate: NTLM

#### The Handshake

Next you need to work a three-part handshake as follows:

1. First, the client establishes a network path to the server and sends a NEGOTIATE_MESSAGE advertising its capabilities.
2. Next, the server responds with CHALLENGE_MESSAGE which is used to establish the identity of the client.
3. Finally, the client responds to the challenge with an AUTHENTICATE_MESSAGE.

To understand the handshake in more detail, I found this write-up quite good: [NTLM Authentication Scheme for HTTP](https://www.innovation.ch/personal/ronald/ntlm.html)

![Figure 1](https://raw.githubusercontent.com/DavidSantia/exchange-ntlm/master/README-figure1.png)

#### Keeping the connection alive
NTLM authenticates connections, not requests. In Go-language, this means you are using the RoundTripper, i.e.:
* https://golang.org/pkg/net/http/#Transport

and then implementing the handshake natively.

When you begin your handshake, the network connection must be kept alive, i.e. between the receiving of the type-2 message from the server (step 2b) and the sending of the type-3 message (step 3). Each time the connection is closed this second part (steps 2 and 3) must be repeated over the new connection (i.e. it's not enough to just keep sending the last type-3 message).

Once the connection is authenticated, the Authorization header need not be sent anymore while the connection stays open, no matter what Exchange resource is accessed.

### One Go libary to consider

This library is referred to a lot and has much of the parts for a native Go implementation:
* https://github.com/ThomsonReutersEikon/go-ntlm

I tried it with an Exchange server, hacked around a bit, but didn't get it to work. I then read the fine print: "The major missing piece is the negotiation of capabilities between the client and the server, for our use we hardcoded a supported set of negotiation flags."

#### Curl to the rescue
I needed a working start point. The curl program has a cool command-line option: **-ntlm**

Here is a curl example:
```sh
curl -v --ntlm \
  -u 'you@yourdomain.com:YourPassword' \
  -H 'Content-Type: text/xml;charset=UTF-8' \
  -d @finditem.xml \
  https://example.com/EWS/Exchange.asmx
```

Once you establish the handshake, like above, you post SOAP commands.  These are just XML, and get back XML responses.

#### Checking your inbox

The most basic command is check your inbox.  It uses the FindItem operation, which I copied into the file finditem.xml in this repo.

To run the curl command to check your inbox:

1. Copy [finditem.xml](https://raw.githubusercontent.com/DavidSantia/exchange-ntlm/master/finditem.xml) into your current directory.
2. Then just type the curl example shown above. 

Microsoft has all the XML you need nicely documented; below are links to examples:

* [FindItem operation](https://msdn.microsoft.com/en-us/library/office/aa566107(v=exchg.150).aspx#sectionSection1)
* [GetItem operation](https://msdn.microsoft.com/en-us/library/office/aa566013(v=exchg.150).aspx#Anchor_1)
* [FindFolder operation](https://msdn.microsoft.com/en-us/library/office/dd633627(v=exchg.80).aspx#Anchor_0)

#### Libcurl in Go

Once I had everything working with curl, I kept going. It turns out you can actually just use the C library *libcurl* in Go:
* https://github.com/andelf/go-curl

For the general case, and definately to jump-start your NTLM authentication, libcurl is a really good route.

The sample Go code in this repo does the same thing as the curl example command.
* It uses the **go-curl** library, Copyright 2014 Shuyu Wang (<andelf@gmail.com>)
* You will notice you can turn on debug (which behaves just like curl -v), and watch the handshake steps proceed.
* This makes it super easy to see if your Go program's handshake matches your curl command.

#### Notes on using minimal containers and C-libraries in Go

Using a C-library in GO means you cannot disable **cgo** when you compile.  Let's say you are trying to build with Go for an empty container, one that doesn't have a full OS.  It's nice to disable cgo (setting **CGO_ENABLE=0**) to avoid dynamic libraries becoming dependencies for you executable.  In your minimal container, the shared libraries of an OS don't exist. But that setting contradicts using a C-library in Go.

You can still use libcurl in Go and link with static libraries.
```sh
go build --ldflags '-extldflags "-static"'
```

It just may take some work to build **libcurl.a** from scratch.
