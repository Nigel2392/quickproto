package tests

import (
	"testing"

	"github.com/Nigel2392/quickproto"
	"github.com/Nigel2392/quickproto/client"
	"github.com/Nigel2392/quickproto/server"
)

func TestConnection(t *testing.T) {
	conf := quickproto.Config{
		Delimiter: []byte("&"),
		UseBase64: true,
		BufSize:   4096,
	}

	IP := "127.0.0.1"
	Port := 8080
	go func() {
		s := server.New(IP, Port, &conf)
		s.Listen()
		for {
			conn, _ := s.Accept()
			msg, _ := s.Read(conn)
			s.Write(conn, msg)
		}
	}()
	c := client.New(IP, Port, &conf)
	c.Connect()
	msg := quickproto.NewMessage(c.Delimiter, c.Use_Base64)
	msg.AddHeader("Test", "Test")
	msg.AddHeader("Test2", "Test2")
	msg.AddHeader("Test3", "Test3")
	msg.AddRawFile("test.txt", []byte("Hello World"))
	msg.AddRawFile("test2.txt", []byte("Hello World"))
	msg.AddRawFile("test3.txt", []byte("Hello World"))
	msg.Body = []byte("Hello World")
	c.Write(msg)
	newmsg, _ := c.Read()

	// Validate
	if newmsg.Headers["Test"][0] != "Test" {
		t.Error("Header Test not equal to Test")
	}
	if newmsg.Headers["Test2"][0] != "Test2" {
		t.Error("Header Test2 not equal to Test2")
	}
	if newmsg.Headers["Test3"][0] != "Test3" {
		t.Error("Header Test3 not equal to Test3")
	}
	if newmsg.Files["test.txt"].Name == "" {
		t.Error("File test.txt no file.")
	}
	if newmsg.Files["test2.txt"].Name == "" {
		t.Error("File test2.txt no file.")
	}
	if newmsg.Files["test3.txt"].Name == "" {
		t.Error("File test3.txt no file.")
	}
}
