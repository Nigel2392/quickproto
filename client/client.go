package client

import (
	"errors"
	"net"
	"strconv"

	"github.com/Nigel2392/quickproto"
	"github.com/Nigel2392/simplecrypto/aes"
)

type Client struct {
	IP          string
	PORT        int
	Conn        net.Conn
	UseEncoding bool
	UseCrypto   bool
	Delimiter   []byte
	BUF_SIZE    int
	Enc_func    func([]byte) []byte
	Dec_func    func([]byte) ([]byte, error)
	CONFIG      *quickproto.Config
	OnMessage   func(*quickproto.Message)
	AesKey      *[32]byte
}

func New(ip string, port int, conf *quickproto.Config, onmessage func(*quickproto.Message)) *Client {
	return &Client{
		IP:          ip,
		PORT:        port,
		Conn:        nil,
		UseEncoding: conf.UseEncoding,
		UseCrypto:   conf.UseCrypto,
		Delimiter:   conf.Delimiter,
		BUF_SIZE:    conf.BufSize,
		Enc_func:    conf.Enc_func,
		Dec_func:    conf.Dec_func,
		CONFIG:      conf,
		OnMessage:   onmessage,
	}
}

func (c *Client) Addr() string {
	return c.IP + ":" + strconv.Itoa(c.PORT)
}

func (c *Client) Connect() error {
	var err error
	c.Conn, err = net.Dial("tcp", c.Addr())
	if c.UseCrypto && c.AesKey == nil {
		aes_key := aes.NewEncryptionKey()
		// send aes key to server
		msg := c.CONFIG.NewMessage()
		msg.Headers["type"] = []string{"aes_key"}
		msg.Body = aes_key[:]
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
