Fast golang protocol parser

Supports:
* Headers (Support listed values, IE: Key=[Value, Value, Value])
* Files (Supports multiple files)
* Body (Can be encoded/decoded to any type)
  * !WARNING! !Make sure your delimiter not in the encoding's alphabet!
  * !WARNING! !Make sure your delimiter is not in the message headers, body or filenames!
* Delimiters
  * No alphabetic characters from [A-Z a-z 0-9 =]

Data is split apart by the delimiter.

Say the delimiter is a `$`:

The body and header will be split by the delimiter * 4

Each key value pair will be split by the delimiter * 2

Then the key values will be split by the delimiter.

Example:
```
key1$value1&value2$$key2$value2$$$$BODYBODYBODY
```

Usage:

Initialize a config like so:
```go
conf := quickproto.NewConfig([]byte(DELIMITER), USE_ENCODING, USE_CRYPTO, 2048, quickproto.Base16Encoding, quickproto.Base16Decoding)
// RSA only used if USE_CRYPTO is true, and when sending the AES key from client to server.
// The RSA keys are however not required, but highly recommended to securely send the AES key from client to server!
conf.PrivateKey = privkey // Client does not need the private key! This is a security risk!
conf.PublicKey = pubkey // Server does not need the public key, but it would not pose a security risk.

```
Then you can simply run a server with the following lines of code:
```go
s := server.New(IP, Port, conf)
s.Listen()
for {
	conn, client, err := s.Accept()
	msg, err := s.Read(client)
}
```

Or create a client like so:
```go
c := client.New(IP, Port, conf, nil)
c.Connect()
msg := c.CONFIG.NewMessage() // Use config to generate message to omit providing arguments
msg.AddHeader("Test", "Test")
msg.AddHeader("Test2", "Test2")
msg.AddHeader("Test3", "Test3")
msg.AddRawFile("test.txt", []byte("Hello World"))
msg.AddRawFile("test2.txt", []byte("Hello World"))
msg.AddRawFile("test3.txt", []byte("Hello World"))
msg.AddContent("Hello World")
c.Write(msg)
```

It is also possible to broadcast to multiple clients at once, 
simply use `s.Broadcast(msg)`.

To capture broadcasts on the client side and interact with them, run a goroutine like so:
```go
func OnBroadcast(msg *quickproto.Message) {
  // Do something with the message
}

c := client.New(IP, Port, conf, nil)
c.Connect()
c.OnMessage = OnBroadcast
go c.Listen()
```

