package server

import (
	"errors"
	"net"
	"strings"

	"github.com/Nigel2392/quickproto"
	simple_rsa "github.com/Nigel2392/simplecrypto/rsa"
)

// Server struct.
type Server struct {
	// Address of the server.
	IP   string
	PORT any
	// Listener for connections.
	Listener net.Listener
	// General configuration.
	CONFIG  *quickproto.Config
	Clients map[string]*Client
}

// Server-side client.
type Client struct {
	Conn net.Conn
	Key  *[32]byte
	// Cookies
	Cookies    map[string][]string
	delCookies []string
	setCookies map[string][]string
	// Data is used for storing extra data about the client server side.
	Data any
}

func (c *Client) AddCookie(key string, value string) {
	c.setCookies[key] = append(c.setCookies[key], value)
}

func (c *Client) SetCookies(key string, values []string) {
	c.setCookies[key] = values
}

func (c *Client) DeleteCookie(key string) {
	c.delCookies = append(c.delCookies, key)
}

func (c *Client) GetCookie(key string) []string {
	return c.Cookies[key]
}

// Initialize a new server.
func New(ip string, port any, conf *quickproto.Config) *Server {
	return &Server{
		IP:       ip,
		PORT:     port,
		Listener: nil,
		CONFIG:   conf,
		Clients:  make(map[string]*Client),
	}
}

// Get the address of the server.
func (s *Server) Addr() string {
	return quickproto.CraftAddr(s.IP, s.PORT)
}

// Listen for connections
func (s *Server) Listen() (net.Listener, error) {
	var err error
	s.Listener, err = net.Listen("tcp", s.Addr())
	return s.Listener, err
}

// Close the server
func (s *Server) Terminate() error {
	return s.Listener.Close()
}

// Accept a new client connection.
// If the server is using crypto, the first message received from the client will be the AES key.
// If the server is provided with a private key, it will use it to decrypt the AES key.
// The server will then use the AES key to decrypt and encrypt all future messages.
func (s *Server) Accept() (net.Conn, *Client, error) {
	conn, err := s.Listener.Accept()
	if err != nil {
		return nil, &Client{}, err
	}
	// If we are using crypto, the first message sent by the client will be the AES key.
	// This key will be used to encrypt all future messages.
	// If we are provided with a private key, we will use it to decrypt the AES key.
	// If we are not provided with a private key, we will assume that the client is not using RSA encryption.
	client := &Client{
		Conn:       conn,
		Cookies:    make(map[string][]string),
		setCookies: make(map[string][]string),
		delCookies: make([]string, 0),
	}
	if s.CONFIG.UseCrypto {
		// read aes key from client.
		msg, err := s.Read(client)
		if err != nil {
			return nil, &Client{}, err
		}
		if s.CONFIG.PrivateKey != nil {
			if msg.Body, err = quickproto.Base64Decoding(msg.Body); err != nil {
				return nil, &Client{}, err
			}
			if msg.Body, err = simple_rsa.Decrypt(msg.Body, s.CONFIG.PrivateKey); err != nil {
				return nil, &Client{}, err
			}
		}
		typ, ok := msg.Headers["type"]
		if !ok {
			return nil, &Client{}, errors.New("no type header")
		}
		if typ[0] != "aes_key" {
			return nil, &Client{}, errors.New("client did not send aes key")
		}
		// convert key to byte array.
		aes_key := new([32]byte)
		copy(aes_key[:], msg.Body)
		client.Key = aes_key
	}
	s.Clients[conn.RemoteAddr().String()] = client
	return conn, client, nil
}

// Read a message from a client.
func (s *Server) Read(client *Client) (*quickproto.Message, error) {
	msg, err := quickproto.ReadConn(client.Conn, s.CONFIG, client.Key, s.CONFIG.Compressed)
	if err != nil {
		return nil, err
	}
	// Logic for handling cookies.
	for key, cookie := range msg.Headers {
		if strings.HasPrefix(key, "Q-COOKIES-") {
			n_key := strings.TrimPrefix(key, "Q-COOKIES-")
			client.Cookies[n_key] = cookie
			delete(msg.Headers, key)
		}
	}
	return msg, nil
}

// Write a message to a client.
func (s *Server) Write(client *Client, msg *quickproto.Message) error {
	for key, cookie := range client.setCookies {
		msg.Headers["Q-SET-COOKIES-"+key] = append(msg.Headers["Q-SET-COOKIES-"+key], cookie...)
	}
	for _, key := range client.delCookies {
		msg.Headers["Q-DEL-COOKIES-"+key] = []string{"\x00"}
	}
	return quickproto.WriteConn(client.Conn, msg, client.Key, s.CONFIG.Compressed)
}

// Close a client connection.
func (s *Server) RemoveClient(conn net.Conn) error {
	delete(s.Clients, conn.RemoteAddr().String())
	return conn.Close()
}

// Broadcast a message to all clients.
func (s *Server) Broadcast(msg *quickproto.Message) error {
	for _, client := range s.Clients {
		if err := s.Write(client, msg); err != nil {
			return err
		}
	}
	return nil
}
