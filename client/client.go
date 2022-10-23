package client

import (
	"crypto/rsa"
	"errors"
	"net"
	"strconv"

	"github.com/Nigel2392/quickproto"
	"github.com/Nigel2392/simplecrypto/aes"
	simple_rsa "github.com/Nigel2392/simplecrypto/rsa"
)

type Client struct {
	IP           string
	PORT         int
	Conn         net.Conn
	UseEncoding  bool
	UseCrypto    bool
	Delimiter    []byte
	BUF_SIZE     int
	Encode_func  func([]byte) []byte
	Decode_func  func([]byte) ([]byte, error)
	CONFIG       *quickproto.Config
	OnMessage    func(*quickproto.Message)
	AesKey       *[32]byte
	RsaPublicKey *rsa.PublicKey
}

func New(ip string, port int, conf *quickproto.Config, onmessage func(*quickproto.Message)) *Client {
	return &Client{
		IP:           ip,
		PORT:         port,
		Conn:         nil,
		UseEncoding:  conf.UseEncoding,
		UseCrypto:    conf.UseCrypto,
		Delimiter:    conf.Delimiter,
		BUF_SIZE:     conf.BufSize,
		Encode_func:  conf.Encode_func,
		Decode_func:  conf.Decode_func,
		CONFIG:       conf,
		OnMessage:    onmessage,
		RsaPublicKey: conf.PublicKey,
	}
}

func (c *Client) Addr() string {
	return c.IP + ":" + strconv.Itoa(c.PORT)
}

func (c *Client) Connect() error {
	var err error
	c.Conn, err = net.Dial("tcp", c.Addr())
	if c.UseCrypto && c.AesKey == nil {
		// Generate new aes key each session
		aes_key := aes.NewEncryptionKey()
		// Generate message to send to server
		msg := c.CONFIG.NewMessage()
		msg.Headers["type"] = []string{"aes_key"}
		msg.Body = aes_key[:]
		// Encrypt body with public key when one is provided
		if c.RsaPublicKey != nil {
			msg.Body, err = simple_rsa.Encrypt(msg.Body, c.RsaPublicKey)
			if err != nil {
				return err
			}
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

func (c *Client) Terminate() error {
	return c.Conn.Close()
}

func (c *Client) Read() (*quickproto.Message, error) {
	msg, err := quickproto.ReadConn(c.Conn, c.CONFIG, c.AesKey)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (c *Client) Write(msg *quickproto.Message) error {
	return quickproto.WriteConn(c.Conn, msg, c.AesKey)
}

func (c *Client) Listen() error {
	for {
		msg, err := c.Read()
		if err != nil {
			break
		}
		c.OnMessage(msg)
	}
	return errors.New("connection closed")
}
