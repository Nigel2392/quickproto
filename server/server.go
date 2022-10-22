package server

import (
	"errors"
	"net"
	"strconv"

	"github.com/Nigel2392/quickproto"
)

type Server struct {
	IP          string
	PORT        int
	Listener    net.Listener
	UseEncoding bool
	UseCrypto   bool
	Delimiter   []byte
	BUF_SIZE    int
	Enc_func    func([]byte) []byte
	Dec_func    func([]byte) ([]byte, error)
	CONFIG      *quickproto.Config
	Clients     map[string]*Client
}

type Client struct {
	Conn net.Conn
	Key  *[32]byte
}

func New(ip string, port int, conf *quickproto.Config) *Server {
	return &Server{
		IP:          ip,
		PORT:        port,
		Listener:    nil,
		UseEncoding: conf.UseEncoding,
		UseCrypto:   conf.UseCrypto,
		Delimiter:   conf.Delimiter,
		BUF_SIZE:    conf.BufSize,
		Enc_func:    conf.Enc_func,
		Dec_func:    conf.Dec_func,
		CONFIG:      conf,
		Clients:     make(map[string]*Client),
	}
}

func (s *Server) Addr() string {
	return s.IP + ":" + strconv.Itoa(s.PORT)
}

func (s *Server) Listen() (net.Listener, error) {
	var err error
	s.Listener, err = net.Listen("tcp", s.Addr())
	return s.Listener, err
}

func (s *Server) Terminate() error {
	return s.Listener.Close()
}

func (s *Server) Accept() (net.Conn, *Client, error) {
	conn, err := s.Listener.Accept()
	if err != nil {
		return nil, &Client{}, err
	}
	var aes_key *[32]byte
	if s.UseCrypto {
		// read aes key from client
		msg, err := s.Read_c(conn)
		if err != nil {
			return nil, &Client{}, err
		}
		if msg.Headers["type"][0] != "aes_key" {
			return nil, &Client{}, errors.New("client did not send aes key")
		}
		// convert key to byte array
		aes_key = new([32]byte)
		copy(aes_key[:], msg.Body)
	}
	client := &Client{
		Conn: conn,
		Key:  aes_key,
	}
	s.Clients[conn.RemoteAddr().String()] = client
	return conn, client, nil
}

func (s *Server) Read_c(conn net.Conn) (*quickproto.Message, error) {
	return quickproto.ReadConn(conn, s.CONFIG, nil)
}

func (s *Server) Write_c(conn net.Conn, msg *quickproto.Message) error {
	return quickproto.WriteConn(conn, msg, nil)
}

func (s *Server) Read(client *Client) (*quickproto.Message, error) {
	return quickproto.ReadConn(client.Conn, s.CONFIG, client.Key)
}

func (s *Server) Write(client *Client, msg *quickproto.Message) error {
	return quickproto.WriteConn(client.Conn, msg, client.Key)
}

func (s *Server) RemoveClient(conn net.Conn) {
	delete(s.Clients, conn.RemoteAddr().String())
}

func (s *Server) Broadcast(msg *quickproto.Message) error {
	for _, client := range s.Clients {
		err := s.Write(client, msg)
		if err != nil {
			return err
		}
	}
	return nil
}
