package client

import (
	"net"
	"strconv"

	"github.com/Nigel2392/quickproto"
)

type Client struct {
	IP          string
	PORT        int
	Conn        net.Conn
	UseEncoding bool
	Delimiter   []byte
	BUF_SIZE    int
	Enc_func    func([]byte) []byte
	Dec_func    func([]byte) ([]byte, error)
	CONFIG      *quickproto.Config
}

func New(ip string, port int, conf *quickproto.Config) *Client {
	return &Client{
		IP:          ip,
		PORT:        port,
		Conn:        nil,
		UseEncoding: conf.UseEncoding,
		Delimiter:   conf.Delimiter,
		BUF_SIZE:    conf.BufSize,
		Enc_func:    conf.Enc_func,
		Dec_func:    conf.Dec_func,
		CONFIG:      conf,
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
	return quickproto.ReadConn(c.Conn, c.CONFIG)
}

func (c *Client) Write(msg *quickproto.Message) error {
	return quickproto.WriteConn(c.Conn, msg)
}
