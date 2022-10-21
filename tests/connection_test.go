//go:build !race
// +build !race

package tests

import (
	"errors"
	"strconv"
	"sync"
	"testing"

	"github.com/Nigel2392/quickproto"
	"github.com/Nigel2392/quickproto/client"
	"github.com/Nigel2392/quickproto/server"
)

// This test is not race condition safe, do not run with --race!
func TestConnection(t *testing.T) {
	//  "=",
	var BufSizes []int = []int{16, 128, 512, 1024, 2048, 16384}
	var UseB64 []bool = []bool{false, true}
	var DELIMITER_LIST []string = []string{
		"!", "@", "#", "$", "%", "^", "&", "*", "(", ")", "-", "_", "+", "[", "{", "]", "}", ";", ":", "'", "\"", ",", "<", ".", ">", "/", "?", "`", "~", "|", "\\", " ",
		// Ascii escape characters
		"\x1b", "\x1c", "\x1d", "\x1e", "\x1f",
		// Ascii control characters
		"\x00", "\x01", "\x02", "\x03", "\x04", "\x05", "\x06", "\x07", "\x08", "\x09", "\x0a", "\x0b", "\x0c", "\x0d", "\x0e", "\x0f", "\x10", "\x11", "\x12", "\x13", "\x14",
		"\x15", "\x16", "\x17", "\x18", "\x19", "\x1a",
		// Ascii non-printable characters
		"\x7f", "\x80", "\x81", "\x82", "\x83", "\x84", "\x85", "\x86", "\x87", "\x88", "\x89", "\x8a", "\x8b", "\x8c", "\x8d", "\x8e", "\x8f", "\x90", "\x91", "\x92", "\x93",
		"\x94", "\x95", "\x96", "\x97", "\x98", "\x99", "\x9a", "\x9b", "\x9c", "\x9d", "\x9e", "\x9f",
	}
	var ct int
	var FAILED_DELIMITERS []error = []error{}
	for _, BUFFER_SIZE := range BufSizes {
		wg := &sync.WaitGroup{}
		wg.Add(len(DELIMITER_LIST) * len(UseB64))
		for _, USAGE := range UseB64 {
			ct += len(DELIMITER_LIST) / 2
			go func(wg *sync.WaitGroup, DELIMITER_LIST []string, USAGE bool, BUFFER_SIZE int) {
				for _, DELIMITER := range DELIMITER_LIST {
					ct++

					conf := quickproto.Config{
						Delimiter: []byte(DELIMITER),
						UseBase64: USAGE,
						BufSize:   BUFFER_SIZE,
					}

					IP := "127.0.0.1"
					Port := 8080 + ct
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
					// Add headers to message
					msg.AddHeader("Test", "Test")
					msg.AddHeader("Test2", "Test2")
					msg.AddHeader("Test3", "Test3")
					// Add files to message
					msg.AddRawFile("test.txt", []byte("Hello World"))
					msg.AddRawFile("test2.txt", []byte("Hello World"))
					msg.AddRawFile("test3.txt", []byte("Hello World"))
					// Add body to message
					msg.AddContent("Hello World")
					c.Write(msg)
					newmsg, _ := c.Read()

					// Validate
					if newmsg.Headers["Test"][0] != "Test" {
						FAILED_DELIMITERS = append(FAILED_DELIMITERS, errors.New("Current Delimiter:"+DELIMITER+"\nHeader Test not equal to Test"))
					}
					if newmsg.Headers["Test2"][0] != "Test2" {
						FAILED_DELIMITERS = append(FAILED_DELIMITERS, errors.New("Current Delimiter:"+DELIMITER+"\nHeader Test2 not equal to Test2"))
					}
					if newmsg.Headers["Test3"][0] != "Test3" {
						FAILED_DELIMITERS = append(FAILED_DELIMITERS, errors.New("Current Delimiter:"+DELIMITER+"\nHeader Test3 not equal to Test3"))
					}
					if newmsg.Files["test.txt"].Name != "test.txt" {
						FAILED_DELIMITERS = append(FAILED_DELIMITERS, errors.New("Current Delimiter:"+DELIMITER+"\nFile test.txt no file."))
					}
					if newmsg.Files["test2.txt"].Name != "test2.txt" {
						FAILED_DELIMITERS = append(FAILED_DELIMITERS, errors.New("Current Delimiter:"+DELIMITER+"\nFile test2.txt no file."))
					}
					if newmsg.Files["test3.txt"].Name != "test3.txt" {
						FAILED_DELIMITERS = append(FAILED_DELIMITERS, errors.New("Current Delimiter:"+DELIMITER+"\nFile test3.txt no file."))
					}
					if string(newmsg.Files["test.txt"].Data) != "Hello World" {
						FAILED_DELIMITERS = append(FAILED_DELIMITERS, errors.New("Current Delimiter:"+DELIMITER+"\nFile test.txt data not equal to Hello World"))
					}
					if string(newmsg.Files["test2.txt"].Data) != "Hello World" {
						FAILED_DELIMITERS = append(FAILED_DELIMITERS, errors.New("Current Delimiter:"+DELIMITER+"\nFile test2.txt data not equal to Hello World"))
					}
					if string(newmsg.Files["test3.txt"].Data) != "Hello World" {
						FAILED_DELIMITERS = append(FAILED_DELIMITERS, errors.New("Current Delimiter:"+DELIMITER+"\nFile test3.txt data not equal to Hello World"))
					}
					if string(newmsg.Body) != "Hello World" {
						FAILED_DELIMITERS = append(FAILED_DELIMITERS, errors.New("Current Delimiter:"+DELIMITER+"\nBody is empty."))
					}
				}
				wg.Done()
			}(wg, DELIMITER_LIST, USAGE, BUFFER_SIZE)
			for _, err := range FAILED_DELIMITERS {
				t.Error("Failed Delimiter: ", err.Error(), "\nUse Base64: "+strconv.FormatBool(USAGE)+"\nBuffer size: "+strconv.Itoa(BUFFER_SIZE)+"\n__________________________________\n")
			}
		}
	}
}
