package client

import (
	"net"
	"strconv"

	"github.com/Nigel2392/quickproto"
)

type Client struct {
	IP         string
	PORT       int
	Conn       net.Conn
	Use_Base64 bool
	Delimiter  []byte
	BUF_SIZE   int
}

func NewClient(ip string, port int, use_b64 bool, delim []byte, buf_size int) *Client {
	return &Client{
		IP:         ip,
		PORT:       port,
		Conn:       nil,
		Use_Base64: use_b64,
		Delimiter:  delim,
		BUF_SIZE:   buf_size,
	}
}

func (c *Client) Addr() string {
	return c.IP + ":" + strconv.Itoa(c.PORT)
}

func (c *Client) Connect() error {
	var err error
	c.Conn, err = net.Dial("tcp", c.Addr())
	return err
}

func (c *Client) Terminate() error {
	return c.Conn.Close()
}

func (c *Client) Read() (*quickproto.Message, error) {
	return quickproto.ReadConn(c.Conn, c.Delimiter, c.Use_Base64, c.BUF_SIZE)
}

func (c *Client) Write(msg *quickproto.Message) error {
	return quickproto.WriteConn(c.Conn, msg)
}
