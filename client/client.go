package client

import (
	"net"
	"strings"

	"github.com/Nigel2392/quickproto"
	"github.com/Nigel2392/simplecrypto/aes"
	simple_rsa "github.com/Nigel2392/simplecrypto/rsa"
)

// Client struct for connecting to a server.
type Client struct {
	IP        string
	PORT      any
	Conn      net.Conn
	CONFIG    *quickproto.Config
	OnMessage func(*quickproto.Message)
	AesKey    *[32]byte
	Cookies   map[string][]string
}

// Initiate a new client.
func New(ip string, port any, conf *quickproto.Config, onmessage func(*quickproto.Message)) *Client {
	return &Client{
		IP:        ip,
		PORT:      port,
		Conn:      nil,
		CONFIG:    conf,
		OnMessage: onmessage,
		Cookies:   make(map[string][]string),
	}
}

// Addr returns the address of the server, in the form "ip:port".
func (c *Client) Addr() string {
	return quickproto.CraftAddr(c.IP, c.PORT)
}

// Connect to the server.
func (c *Client) Connect() error {
	// If we are using crypto, the first message sent by the client will be the AES key.
	// If the client is provided with a public key, it will use it to encrypt the AES key.
	// The server will then use its private key to decrypt the AES key.
	// Then, the server will use the AES key to decrypt all future messages.
	var err error
	c.Conn, err = net.Dial("tcp", c.Addr())
	if c.CONFIG.UseCrypto && c.AesKey == nil {
		// Generate new aes key each session
		aes_key := aes.NewEncryptionKey()
		// Generate message to send to server
		msg := c.CONFIG.NewMessage()
		msg.AddHeader("type", "aes_key")
		msg.Body = aes_key[:]
		// Encrypt body with public key when one is provided
		if c.CONFIG.PublicKey != nil {
			msg.Body, err = simple_rsa.Encrypt(msg.Body, c.CONFIG.PublicKey)
			if err != nil {
				return err
			}
			// Encode message body
			msg.Body = quickproto.Base64Encoding(msg.Body)
		}
		// Send message
		err = c.Write(msg)
		if err != nil {
			return err
		}
		c.AesKey = aes_key
	}
	return err
}

// Terminate the connection.
func (c *Client) Terminate() error {
	return c.Conn.Close()
}

// Read a message from the server.
func (c *Client) Read() (*quickproto.Message, error) {
	msg, err := quickproto.ReadConn(c.Conn, c.CONFIG, c.AesKey, c.CONFIG.Compressed)
	if err != nil {
		return nil, err
	}
	for k, v := range msg.Headers {
		if strings.HasPrefix(k, "Q-SET-COOKIES-") {
			c.Cookies[strings.TrimPrefix(k, "Q-SET-COOKIES-")] = v
		} else if strings.HasPrefix(k, "Q-DEL-COOKIES-") {
			delete(c.Cookies, strings.TrimPrefix(k, "Q-DEL-COOKIES-"))
		}
	}
	return msg, nil
}

// Write a message to the server.
func (c *Client) Write(msg *quickproto.Message) error {
	for k, v := range c.Cookies {
		for _, v2 := range v {
			msg.AddHeader("Q-COOKIES-"+k, v2)
		}
	}
	return quickproto.WriteConn(c.Conn, msg, c.AesKey, c.CONFIG.Compressed)
}

// Listen for messages from the server.
func (c *Client) Listen() error {
	for {
		msg, err := c.Read()
		if err != nil {
			return err
		}
		c.OnMessage(msg)
	}
}

func (c *Client) GetCookies(key string) []string {
	values, ok := c.Cookies[key]
	if !ok {
		return nil
	}
	return values
}
