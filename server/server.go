package server

import (
	"crypto/rsa"
	"errors"
	"net"
	"strconv"

	"github.com/Nigel2392/quickproto"
	simple_rsa "github.com/Nigel2392/simplecrypto/rsa"
)

type Server struct {
	IP            string
	PORT          int
	Listener      net.Listener
	UseEncoding   bool
	UseCrypto     bool
	Delimiter     []byte
	BUF_SIZE      int
	Enc_func      func([]byte) []byte
	Dec_func      func([]byte) ([]byte, error)
	CONFIG        *quickproto.Config
	Clients       map[string]*Client
	RSAPrivateKey *rsa.PrivateKey
}

type Client struct {
	Conn net.Conn
	Key  *[32]byte
}

func New(ip string, port int, conf *quickproto.Config) *Server {
	return &Server{
		IP:            ip,
		PORT:          port,
		Listener:      nil,
		UseEncoding:   conf.UseEncoding,
		UseCrypto:     conf.UseCrypto,
		Delimiter:     conf.Delimiter,
		BUF_SIZE:      conf.BufSize,
		Enc_func:      conf.Enc_func,
		Dec_func:      conf.Dec_func,
		RSAPrivateKey: conf.PrivateKey,
		CONFIG:        conf,
		Clients:       make(map[string]*Client),
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
	// If we are using crypto, the first message sent by the client will be the AES key.
	// This key will be used to encrypt all future messages.
	// If we are provided with a private key, we will use it to decrypt the AES key.
	// If we are not provided with a private key, we will assume that the client is not using RSA encryption.
	var aes_key *[32]byte
	if s.UseCrypto {
		// read aes key from client
		msg, err := s.Read_c(conn)
		if s.RSAPrivateKey != nil {
			msgbody, err := simple_rsa.Decrypt(msg.Body, s.RSAPrivateKey)
			if err != nil {
				return nil, &Client{}, err
			}
			msg.Body = msgbody
		}
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
