//go:build !race
// +build !race

package tests

import (
	"errors"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/Nigel2392/quickproto"
	"github.com/Nigel2392/quickproto/client"
	"github.com/Nigel2392/quickproto/server"
	simple_rsa "github.com/Nigel2392/simplecrypto/rsa"
)

// This test is not race condition safe, do not run with --race!
func TestConnection(t *testing.T) {
	var UseB64 []bool = []bool{false, true}
	var CONNTYPE string = "udp"
	var DELIMITER_LIST []string = []string{
		"!", "@", "#", "$", "%", "^", "&", "*", "+", "[", "{", "]", "}", ";", ":", "'", "\"", ",", "<", ".", ">", "/", "?", "`", "~", "|", "\\", " ",
		// Ascii escape characters
		"\x1b", "\x1c", "\x1d", "\x1f",
		//// Ascii control characters
		"\x01", "\x02", "\x04", "\x05", "\x06",
		"\x07", "\x09", "\x0a", "\x0b", "\x0c", "\x0e", "\x0f", "\x10", "\x11", "\x12", "\x13", "\x14",
		"\x15", "\x17", "\x19", "\x1a",
		// Ascii non-printable characters
		// "\x7f", "\x80", "\x81", "\x82", "\x83", "\x84", "\x85", "\x86", "\x87", "\x88", "\x89",
		// "\x8a", "\x8b", "\x8c", "\x8d", "\x8e", "\x8f", "\x90", "\x91", "\x92", "\x93", "\x94",
		// "\x95", "\x96", "\x97", "\x98", "\x99", "\x9a", "\x9b", "\x9c", "\x9d", "\x9e", "\x9f",
		//// Ascii extended characters
		//"\xa0", "\xa1", "\xa2", "\xa3", "\xa4", "\xa5", "\xa6", "\xa7", "\xa8", "\xa9", "\xaa",
		//"\xab", "\xac", "\xad", "\xae", "\xaf", "\xb0", "\xb1", "\xb2", "\xb3", "\xb4", "\xb5",
		//"\xb6", "\xb7", "\xb8", "\xb9", "\xba", "\xbb", "\xbc", "\xbd", "\xbe", "\xbf", "\xc0",
		//"\xc1", "\xc2", "\xc3", "\xc4", "\xc5", "\xc6", "\xc7", "\xc8", "\xc9", "\xca", "\xcb",
		//"\xcc", "\xcd", "\xce", "\xcf", "\xd0", "\xd1", "\xd2", "\xd3", "\xd4", "\xd5", "\xd6",
		//"\xd7", "\xd8", "\xd9", "\xda", "\xdb", "\xdc", "\xdd", "\xde", "\xdf", "\xe0", "\xe1",
		//"\xe2", "\xe3", "\xe4", "\xe5", "\xe6", "\xe7", "\xe8", "\xe9", "\xea", "\xeb", "\xec",
		//"\xed", "\xee", "\xef", "\xf0", "\xf1", "\xf2", "\xf3", "\xf4", "\xf5", "\xf6", "\xf7",
		//"\xf8", "\xf9", "\xfa", "\xfb", "\xfc", "\xfd", "\xfe", "\xff",

		//// Unicode characters
		//"???", "???", "??", "???", "???", "???", "???", "??", "???", "??", "???", "??", "??", "???", "???", "???", "???", "???",
		//"???", "???", "??", "???", "??", "???", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??",
		//"??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??",
		//"??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??",
		//"??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??",
		//"??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??",
		//"??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??",
		//"??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??",
		//"??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??",
		//"??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??",
		//"??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??",
		//"??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??",
		//"??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??",
		//"??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??",
		//"??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??",
		//"??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??",
		//"??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??",
		//"??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??",
		//"??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??",
		//// "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??",
		//"??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??",
		//"??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??",
		//"??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??",
		//"??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??",
		//"??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??", "??",
	}
	var ct int
	var FAILED_DELIMITERS []error = []error{}
	privkey, pubkey, _ := simple_rsa.GenKeypair(2048)
	// var BufSizes []int = []int{16, 128, 512, 1024, 2048, 16384}
	// for _, BUFFER_SIZE := range BufSizes {
	// wg := &sync.WaitGroup{}
	// wg.Add(len(DELIMITER_LIST) * len(UseB64))
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, USE_CRYPTO := range UseB64 {
		if USE_CRYPTO {
			CONNTYPE = "tcp"
		} else {
			CONNTYPE = "udp"
		}
		for _, USAGE := range UseB64 {
			ct += len(DELIMITER_LIST) / 2
			// go func(wg *sync.WaitGroup, DELIMITER_LIST []string, USAGE bool, BUFFER_SIZE int) {
			wg.Add(1)
			go func(wag *sync.WaitGroup, mut *sync.Mutex, DELIMITER_LIST []string, USAGE bool, USE_CRYPTO bool) {
				defer wg.Done()
				for _, DELIMITER := range DELIMITER_LIST {
					var hasBody bool = USAGE && !USE_CRYPTO
					conf := quickproto.NewConfig([]byte(DELIMITER), USAGE, USE_CRYPTO, 2048, quickproto.Base64Encoding, quickproto.Base64Decoding)
					conf.PrivateKey = privkey
					conf.PublicKey = pubkey
					conf.Compressed = true
					IP := "127.0.0.1"
					mut.Lock()
					ct++
					Port := 8080 + ct
					s := server.New(IP, Port, conf)
					go func(t *testing.T, s *server.Server) {
						s.Listen(CONNTYPE)
						for {
							_, client, err := s.Accept()
							if err != nil {
								t.Error(err)
							}
							msg, err := s.Read(client)
							t.Log("Server key:", client.Key)
							if err != nil {
								t.Error(err)
							}
							newmsg := conf.NewMessage()
							_, err = msg.Generate()
							if err != nil {
								t.Error(err)
							}
							for k, v := range msg.Headers {
								for _, v2 := range v {
									newmsg.AddHeader(k, v2)
								}
							}
							newmsg.AddContent(msg.Body)
							for _, v := range msg.Files {
								newmsg.AddFile(v)
							}
							client.AddCookie("test", "test")
							client.AddCookie("test", "test2")
							client.AddCookie("test2", "test")

							err = s.Write(client, newmsg)
							if err != nil {
								t.Error(err)
							}
						}
					}(t, s)
					c := client.New(IP, Port, conf, nil)
					mut.Unlock()
					c.Connect(CONNTYPE)
					t.Log("Client key: ", c.AesKey)
					msg := conf.NewMessage()
					// Add headers to message
					msg.AddHeader("Test", "Test")
					msg.AddHeader("Test2", "Test2")
					msg.AddHeader("Test3", "Test3")
					// Add files to message
					if !hasBody {
						msg.AddRawFile("test.txt", []byte("Hello World"))
						msg.AddRawFile("test2.txt", []byte("Hello World"))
						msg.AddRawFile("test3.txt", []byte("Hello World"))
					}
					// Add body to message
					if hasBody {
						msg.AddContent("Hello World")
					}
					// msg.AddContent("Hello World")
					c.Write(msg)
					newmsg, err := c.Read()
					if err != nil {
						t.Error(err)
					}
					// t.Log(strings.Repeat("-", 50))
					// t.Log("Message struct: ", newmsg)
					// t.Log("Message Headers: ", newmsg.Headers)
					// t.Log("Using crypto: ", USE_CRYPTO)
					// t.Log("Using base64: ", USAGE)
					// t.Log("Message Body: ", string(newmsg.Body))
					// t.Log("Message Files: ")
					// for _, file := range newmsg.Files {
					// t.Log("  File Name: ", file.Name)
					// t.Log("  File Data: ", string(file.Data)+"\n")
					// }
					// t.Log("Message Delimiter: ", string(newmsg.Delimiter), "(bytes:", newmsg.Delimiter, ")")
					t.Log("Client Cookies: ", c.Cookies)
					t.Log("Message Data: ", string(newmsg.Data))
					t.Log(strings.Repeat("-", 50))
					if !hasBody {
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
						if newmsg.Files["test.txt"] == nil {
							t.Error("newmsg.Files[\"test.txt\"] == nil")
						}
						if newmsg.Files["test.txt"].Name != "test.txt" {
							FAILED_DELIMITERS = append(FAILED_DELIMITERS, errors.New("Current Delimiter:"+DELIMITER+"\nFile test.txt no file."))
						}
						if newmsg.Files["test2.txt"] == nil {
							t.Error("newmsg.Files[\"test2.txt\"] == nil")
						}
						if newmsg.Files["test2.txt"].Name != "test2.txt" {
							FAILED_DELIMITERS = append(FAILED_DELIMITERS, errors.New("Current Delimiter:"+DELIMITER+"\nFile test2.txt no file."))
						}
						if newmsg.Files["test3.txt"] == nil {
							t.Error("newmsg.Files[\"test3.txt\"] == nil")
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
						if c.GetCookies("test") == nil {
							t.Error("client.GetCookies(\"test\") == nil")
						} else {
							for _, cookie := range c.GetCookies("test") {
								if cookie != "test" && cookie != "test2" {
									t.Error("client.GetCookies(\"test\")[\"test\"] != \"test\" || client.GetCookies(\"test\")[\"test\"] != \"test2\"")
								}
							}
						}
					}
					if string(newmsg.Body) != "" && !hasBody {
						FAILED_DELIMITERS = append(FAILED_DELIMITERS, errors.New("Current Delimiter:"+DELIMITER+"\nBody is not empty."))
					} else if string(newmsg.Body) != "Hello World" && hasBody {
						FAILED_DELIMITERS = append(FAILED_DELIMITERS, errors.New("Current Delimiter:"+DELIMITER+"\nBody not equal to Hello World"))
					}
				}
			}(&wg, &mu, DELIMITER_LIST, USAGE, USE_CRYPTO)
			// wg.Done()
			// }(wg, DELIMITER_LIST, USAGE, BUFFER_SIZE)
			for _, err := range FAILED_DELIMITERS {
				t.Error("Failed Delimiter: ", err.Error(), "\nUse Base64: "+strconv.FormatBool(USAGE)+"\nUse Crypto:"+strconv.FormatBool(USE_CRYPTO)+"\nBuffer size: "+strconv.Itoa(2048)+"\n__________________________________\n")
			}
		}
	}
	wg.Wait()
}
