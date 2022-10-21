package quickproto

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"strings"
	"sync"
)

// Normal delimiter example:
// $
// Header delimiter example:
// $$
// Body delimiter example:
// $$$$
// File delimiter example:
// $$$$$$
// Ending delimiter example:
// $$$$$$$$

// Banned delimiters include:
// = (equal sign)

var STANDARD_DELIM []byte = []byte("$")
var BANNED_DELIMITERS = []string{
	"=",
	"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
	"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
}

func Base64Encoding(data []byte) []byte {
	var b64_buffer bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &b64_buffer)
	encoder.Write(data)
	encoder.Close()
	return b64_buffer.Bytes()
}

func Base64Decoding(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(data)
	decoder := base64.NewDecoder(base64.StdEncoding, buf)
	return io.ReadAll(decoder)
}

func Base16Encoding(data []byte) []byte {
	// return []byte(hex.EncodeToString(data))
	return data
}

func Base16Decoding(data []byte) ([]byte, error) {
	// return hex.DecodeString(string(data))
	return data, nil
}

type Config struct {
	Delimiter   []byte
	UseEncoding bool
	BufSize     int
	Enc_func    func([]byte) []byte
	Dec_func    func([]byte) ([]byte, error)
}

// NewConfig creates a new Config.
func NewConfig(delimiter []byte, useencoding bool, bufsize int, enc_f func([]byte) []byte, dec_f func([]byte) ([]byte, error)) *Config {
	if delimiter == nil {
		delimiter = STANDARD_DELIM
	}
	for _, d := range BANNED_DELIMITERS {
		if bytes.Contains(delimiter, []byte(d)) {
			panic("Delimiter contains banned characters: " + d)
		}
	}
	return &Config{
		Delimiter:   delimiter,
		UseEncoding: useencoding,
		BufSize:     bufsize,
		Enc_func:    enc_f,
		Dec_func:    dec_f,
	}
}

func (c *Config) NewMessage() *Message {
	return NewMessage(c.Delimiter, c.UseEncoding, c.Enc_func, c.Dec_func)
}

// A Message is a protocol message.
type Message struct {
	Data        []byte
	Delimiter   []byte
	Headers     map[string][]string
	Body        []byte
	Files       map[string]MessageFile
	UseEncoding bool
	Enc_func    func([]byte) []byte
	Dec_func    func([]byte) ([]byte, error)
	// Parsed    bool
	// Generated bool
}

// NewMessage creates a new Message.
func NewMessage(delimiter []byte, useencoding bool, enc_func func([]byte) []byte, dec_func func([]byte) ([]byte, error)) *Message {
	if delimiter == nil {
		delimiter = STANDARD_DELIM
	}
	return &Message{
		Data:        []byte{},
		Delimiter:   delimiter,
		Headers:     make(map[string][]string),
		Body:        []byte{},
		Files:       make(map[string]MessageFile),
		UseEncoding: useencoding,
		Enc_func:    enc_func,
		Dec_func:    dec_func,
		// Parsed:    false,
		// Generated: false,
	}
}

func (m *Message) AddHeader(key string, value string) error {
	if strings.Contains(key, string(m.Delimiter)) {
		return errors.New("header key cannot contain delimiter")
	}
	if strings.Contains(value, string(m.Delimiter)) {
		return errors.New("header value cannot contain delimiter")
	}
	_, ok := m.Headers[key]
	if ok {
		m.Headers[key] = append(m.Headers[key], value)
	} else {
		m.Headers[key] = []string{value}
	}
	return nil
}

func (m *Message) AddContent(content any) error {
	switch content := content.(type) {
	case string:
		m.Body = append(m.Body, []byte(content)...)
	case []byte:
		m.Body = append(m.Body, content...)
	default:
		return errors.New("invalid content type")
	}
	return nil
}

func (m *Message) AddFile(file MessageFile) error {
	// if strings.Contains(file.Name, string(m.Delimiter)) {
	// return errors.New("file name contains file delimiter")
	// }
	m.Files[file.Name] = file
	return nil
}

func (m *Message) AddRawFile(name string, data []byte) error {
	// if strings.Contains(string(name), string(m.Delimiter)) {
	// return errors.New("file name contains file delimiter")
	// }
	m.Files[name] = MessageFile{Name: name, Data: data}
	return nil
}

func (m *Message) HeaderDelimiter() []byte {
	return append(m.Delimiter, m.Delimiter...)
}

func (m *Message) FileDelimiter() []byte {
	return append(m.BodyDelimiter(), m.HeaderDelimiter()...)
}

func (m *Message) BodyDelimiter() []byte {
	return append(m.HeaderDelimiter(), m.HeaderDelimiter()...)
}

func (m *Message) EndingDelimiter() []byte {
	return append(m.BodyDelimiter(), m.BodyDelimiter()...)
}

// parses protocol messages.
// Header is a map of key/value pairs.
// Body is a base64 encoded byte slice.
func (m *Message) Parse() (*Message, error) {
	///////////////////////////////////////////////
	// Splitting order
	// 1. Split body and head
	//	    a. Split head into key/value pairs
	//	    b. Split key/value pairs into key and values
	// 2. Split files from body
	//	    a. Split files into file name and data
	///////////////////////////////////////////////
	header_delimiter := m.HeaderDelimiter()
	file_delimiter := m.FileDelimiter()
	body_delimiter := m.BodyDelimiter()
	ending_delimiter := m.EndingDelimiter()
	// Split data into headers and body
	datalist := bytes.SplitN(m.Data, body_delimiter, 2)
	if len(datalist) != 2 {
		return nil, errors.New("invalid message sent")
	}
	// Get headers from datalist
	// Split headers into key/value pairs
	headers := bytes.Split(datalist[0], header_delimiter)
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, header := range headers {
		wg.Add(1)
		// Start goroutine for each header
		go func(header []byte, wg *sync.WaitGroup, mu *sync.Mutex) {
			defer wg.Done()
			// Split each header into key and values
			head := bytes.Split(header, m.Delimiter)
			str_list := make([]string, 0)
			// Set multiple values for each key
			for _, byt := range head[1:] {
				str_list = append(str_list, string(byt))
			}
			// Set key and values, lock for thread safety
			mu.Lock()
			m.Headers[string(head[0])] = str_list
			mu.Unlock()
		}(header, &wg, &mu)
	}
	// Decode base64 encoded body
	var body []byte
	var err error
	full_body := bytes.Trim(datalist[1], string(ending_delimiter))
	if m.Enc_func != nil && m.Dec_func != nil && m.UseEncoding {
		full_body, err = m.Dec_func(full_body)
		if err != nil {
			return nil, err
		}
	}
	body_data := bytes.Split(full_body, file_delimiter)
	// Extract body from body_data
	body = body_data[len(body_data)-1]
	// Remove body from body_data
	body_data = body_data[:len(body_data)-1]

	// Extract files from body_data
	wg.Add(len(body_data))
	for _, file := range body_data {
		go func(file []byte, wg *sync.WaitGroup, mu *sync.Mutex) {
			defer wg.Done()
			file_data := bytes.Split(file, m.HeaderDelimiter())
			if len(file_data) != 2 {
				return
			}
			file_name := string(file_data[0])
			// File data should always be base64 encoded!
			file_data_bytes, err := Base16Decoding(file_data[1])
			if err != nil {
				return
			}
			mu.Lock()
			m.Files[file_name] = MessageFile{Name: file_name, Data: file_data_bytes}
			mu.Unlock()
		}(file, &wg, &mu)
	}
	wg.Wait()
	m.Body = body
	// m.Parsed = true
	return m, nil
}

// creates a protocol message.
// Header is a map of key/value pairs.
// Body is a base64 encoded byte slice.
func (m *Message) Generate() (*Message, error) {
	var buffer bytes.Buffer
	var wg sync.WaitGroup
	var mu sync.Mutex
	header_delimiter := m.HeaderDelimiter()
	// Create headers
	for key, value := range m.Headers {
		wg.Add(1)
		// Start goroutine for each header
		go func(key string, value []string, buffer *bytes.Buffer, wg *sync.WaitGroup, mu *sync.Mutex) {
			defer wg.Done()
			// Create buffer for length of current header line
			total_len := 0
			for _, str := range value {
				// Append key and value to headerline
				val_len := len(str)
				total_len = total_len + val_len + len(m.Delimiter)
			}
			// Create headerline
			headerline := make([]byte, len(key)+len(m.Delimiter)+total_len+len(m.Delimiter))
			// Copy key to headerline
			copy(headerline, key)
			copy(headerline[len(key):], m.Delimiter)
			// Copy values to headerline
			current_pos := len(key) + len(m.Delimiter)
			for _, str := range value {
				copy(headerline[current_pos:], str)
				copy(headerline[current_pos+len(str):], m.Delimiter)
				current_pos = current_pos + len(str) + len(m.Delimiter)
			}
			// Set last delimiter
			copy(headerline[current_pos:], m.Delimiter)
			// Append headerline to buffer
			// Lock the mutex to prevent datarace
			mu.Lock()
			buffer.Write(headerline)
			mu.Unlock()
		}(key, value, &buffer, &wg, &mu)
	}
	// Wait for all goroutines to finish
	wg.Wait()
	// Create body
	buffer.Write(header_delimiter)
	// Get files
	var bodybuffer bytes.Buffer
	wg.Add(len(m.Files))
	for _, file := range m.Files {
		// Write the file to the body
		go func(file MessageFile, buffer *bytes.Buffer, wg *sync.WaitGroup, mu *sync.Mutex) {
			defer wg.Done()
			// Create buffer for length of current file line
			// Base 64 encode file data
			fdata := Base16Encoding(file.Data)
			// Get size of buffer for all file data and delimiters
			total_len := len(file.Name) + len(m.HeaderDelimiter()) + len(fdata) + len(m.FileDelimiter())
			fileline := make([]byte, total_len)
			// Copy data to fileline
			copy(fileline, file.Name)
			copy(fileline[len(file.Name):], m.HeaderDelimiter())
			copy(fileline[len(file.Name)+len(m.HeaderDelimiter()):], fdata)
			copy(fileline[len(file.Name)+len(m.HeaderDelimiter())+len(fdata):], m.FileDelimiter())
			// Lock the mutex to prevent datarace
			mu.Lock()
			// Append fileline to buffer
			bodybuffer.Write(fileline)
			mu.Unlock()
		}(file, &buffer, &wg, &mu)
	}
	// Wait for all goroutines to finish
	wg.Wait()
	// Append body to buffer
	bodybuffer.Write(m.Body)
	if m.Enc_func != nil && m.Dec_func != nil && m.UseEncoding {
		// If encoding is set, create buffer and encode body
		buffer.Write(m.Enc_func(bodybuffer.Bytes()))
	} else {
		// Else write body to buffer
		buffer.Write(bodybuffer.Bytes())
	}
	// Write ending delimiter to buffer
	buffer.Write(m.EndingDelimiter())
	m.Data = buffer.Bytes()
	// m.Generated = true
	return m, nil
}

// Get content length of the message.
func (m *Message) ContentLength() int {
	return len(m.Data)
}
