package server

import (
	"net"
	"strconv"

	"github.com/Nigel2392/quickproto"
)

type Server struct {
	IP          string
	PORT        int
	Listener    net.Listener
	UseEncoding bool
	Delimiter   []byte
	BUF_SIZE    int
	Enc_func    func([]byte) []byte
	Dec_func    func([]byte) ([]byte, error)
}

func New(ip string, port int, conf *quickproto.Config) *Server {
	return &Server{
		IP:          ip,
		PORT:        port,
		Listener:    nil,
		UseEncoding: conf.UseEncoding,
		Delimiter:   conf.Delimiter,
		BUF_SIZE:    conf.BufSize,
		Enc_func:    conf.Enc_func,
		Dec_func:    conf.Dec_func,
	}
}

func (s *Server) Addr() string {
	return s.IP + ":" + strconv.Itoa(s.PORT)
}

func (s *Server) Listen() error {
	var err error
	s.Listener, err = net.Listen("tcp", s.Addr())
	return err
}

func (s *Server) Terminate() error {
	return s.Listener.Close()
}

func (s *Server) Accept() (net.Conn, error) {
	return s.Listener.Accept()
}

func (s *Server) Read(conn net.Conn) (*quickproto.Message, error) {
	return quickproto.ReadConn(conn, s.Delimiter, s.UseEncoding, s.BUF_SIZE, s.Enc_func, s.Dec_func)
}

func (s *Server) Write(conn net.Conn, msg *quickproto.Message) error {
	return quickproto.WriteConn(conn, msg)
}
