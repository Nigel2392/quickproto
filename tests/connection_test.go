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
		// Ascii extended characters
		"\xa0", "\xa1", "\xa2", "\xa3", "\xa4", "\xa5", "\xa6", "\xa7", "\xa8", "\xa9", "\xaa",
		"\xab", "\xac", "\xad", "\xae", "\xaf", "\xb0", "\xb1", "\xb2", "\xb3", "\xb4", "\xb5",
		"\xb6", "\xb7", "\xb8", "\xb9", "\xba", "\xbb", "\xbc", "\xbd", "\xbe", "\xbf", "\xc0",
		"\xc1", "\xc2", "\xc3", "\xc4", "\xc5", "\xc6", "\xc7", "\xc8", "\xc9", "\xca", "\xcb",
		"\xcc", "\xcd", "\xce", "\xcf", "\xd0", "\xd1", "\xd2", "\xd3", "\xd4", "\xd5", "\xd6",
		"\xd7", "\xd8", "\xd9", "\xda", "\xdb", "\xdc", "\xdd", "\xde", "\xdf", "\xe0", "\xe1",
		"\xe2", "\xe3", "\xe4", "\xe5", "\xe6", "\xe7", "\xe8", "\xe9", "\xea", "\xeb", "\xec",
		"\xed", "\xee", "\xef", "\xf0", "\xf1", "\xf2", "\xf3", "\xf4", "\xf5", "\xf6", "\xf7",
		"\xf8", "\xf9", "\xfa", "\xfb", "\xfc", "\xfd", "\xfe", "\xff",
		// Unicode characters
		"€", "‚", "ƒ", "„", "…", "†", "‡", "ˆ", "‰", "Š", "‹", "Œ", "Ž", "‘", "’", "“", "”", "•",
		"–", "—", "˜", "™", "š", "›", "œ", "ž", "Ÿ", "¡", "¢", "£", "¤", "¥", "¦", "§", "¨", "©",
		"¼", "½", "¾", "¿", "À", "Á", "Â", "Ã", "Ä", "Å", "Æ", "Ç", "È", "É", "Ê", "Ë", "Ì", "Í",
		"Î", "Ï", "Ð", "Ñ", "Ò", "Ó", "Ô", "Õ", "Ö", "×", "Ø", "Ù", "Ú", "Û", "Ü", "Ý", "Þ", "ß",
		"à", "á", "â", "ã", "ä", "å", "æ", "ç", "è", "é", "ê", "ë", "ì", "í", "î", "ï", "ð", "ñ",
		"ò", "ó", "ô", "õ", "ö", "÷", "ø", "ù", "ú", "û", "ü", "ý", "þ", "ÿ", "Ā", "ā", "Ă", "ă",
		"Ą", "ą", "Ć", "ć", "Ĉ", "ĉ", "Ċ", "ċ", "Č", "č", "Ď", "ď", "Đ", "đ", "Ē", "ē", "Ĕ", "ĕ",
		"Ė", "ė", "Ę", "ę", "Ě", "ě", "Ĝ", "ĝ", "Ğ", "ğ", "Ġ", "ġ", "Ģ", "ģ", "Ĥ", "ĥ", "Ħ", "ħ",
		"Ĩ", "ĩ", "Ī", "ī", "Ĭ", "ĭ", "Į", "į", "İ", "ı", "Ĳ", "ĳ", "Ĵ", "ĵ", "Ķ", "ķ", "ĸ", "Ĺ",
		"ĺ", "Ļ", "ļ", "Ľ", "ľ", "Ŀ", "ŀ", "Ł", "ł", "Ń", "ń", "Ņ", "ņ", "Ň", "ň", "ŉ", "Ŋ", "ŋ",
		"Ō", "ō", "Ŏ", "ŏ", "Ő", "ő", "Œ", "œ", "Ŕ", "ŕ", "Ŗ", "ŗ", "Ř", "ř", "Ś", "ś", "Ŝ", "ŝ",
		"Ş", "ş", "Š", "š", "Ţ", "ţ", "Ť", "ť", "Ŧ", "ŧ", "Ũ", "ũ", "Ū", "ū", "Ŭ", "ŭ", "Ů", "ů",
		"Ű", "ű", "Ų", "ų", "Ŵ", "ŵ", "Ŷ", "ŷ", "Ÿ", "Ź", "ź", "Ż", "ż", "Ž", "ž", "ſ", "ƀ", "Ɓ",
		"Ƃ", "ƃ", "Ƅ", "ƅ", "Ɔ", "Ƈ", "ƈ", "Ɖ", "Ɗ", "Ƌ", "ƌ", "ƍ", "Ǝ", "Ə", "Ɛ", "Ƒ", "ƒ", "Ɠ",
		"Ɣ", "ƕ", "Ɩ", "Ɨ", "Ƙ", "ƙ", "ƚ", "ƛ", "Ɯ", "Ɲ", "ƞ", "Ɵ", "Ơ", "ơ", "Ƣ", "ƣ", "Ƥ", "ƥ",
		"Ʀ", "Ƨ", "ƨ", "Ʃ", "ƪ", "ƫ", "Ƭ", "ƭ", "Ʈ", "Ư", "ư", "Ʊ", "Ʋ", "Ƴ", "ƴ", "Ƶ", "ƶ", "Ʒ",
		"Ƹ", "ƹ", "ƺ", "ƻ", "Ƽ", "ƽ", "ƾ", "ƿ", "ǀ", "ǁ", "ǂ", "ǃ", "Ǆ", "ǅ", "ǆ", "Ǉ", "ǈ", "ǉ",
		"Ǌ", "ǋ", "ǌ", "Ǎ", "ǎ", "Ǐ", "ǐ", "Ǒ", "ǒ", "Ǔ", "ǔ", "Ǖ", "ǖ", "Ǘ", "ǘ", "Ǚ", "ǚ", "Ǜ",
		"ǜ", "ǝ", "Ǟ", "ǟ", "Ǡ", "ǡ", "Ǣ", "ǣ", "Ǥ", "ǥ", "Ǧ", "ǧ", "Ǩ", "ǩ", "Ǫ", "ǫ", "Ǭ", "ǭ",
		"Ǯ", "ǯ", "ǰ", "Ǳ", "ǲ", "ǳ", "Ǵ", "ǵ", "Ƕ", "Ƿ", "Ǹ", "ǹ", "Ǻ", "ǻ", "Ǽ", "ǽ", "Ǿ", "ǿ",
		"Ȁ", "ȁ", "Ȃ", "ȃ", "Ȅ", "ȅ", "Ȇ", "ȇ", "Ȉ", "ȉ", "Ȋ", "ȋ", "Ȍ", "ȍ", "Ȏ", "ȏ", "Ȑ", "ȑ",
		"Ȓ", "ȓ", "Ȕ", "ȕ", "Ȗ", "ȗ", "Ș", "ș", "Ț", "ț", "Ȝ", "ȝ", "Ȟ", "ȟ", "Ƞ", "ȡ", "Ȣ", "ȣ",
		"Ȥ", "ȥ", "Ȧ", "ȧ", "Ȩ", "ȩ", "Ȫ", "ȫ", "Ȭ", "ȭ", "Ȯ", "ȯ", "Ȱ", "ȱ", "Ȳ", "ȳ", "ȴ", "ȵ",
		"ȶ", "ȷ", "ȸ", "ȹ", "Ⱥ", "Ȼ", "ȼ", "Ƚ", "Ⱦ", "ȿ", "ɀ", "Ɂ", "ɂ", "Ƀ", "Ʉ", "Ʌ", "Ɇ", "ɇ",
		"Ɉ", "ɉ", "Ɋ", "ɋ", "Ɍ", "ɍ", "Ɏ", "ɏ", "ɐ", "ɑ", "ɒ", "ɓ", "ɔ", "ɕ", "ɖ", "ɗ", "ɘ", "ə",
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
