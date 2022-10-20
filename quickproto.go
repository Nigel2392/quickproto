package quickproto

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"sync"
)

var STANDARD_DELIM []byte = []byte("$")

type Config struct {
	Delimiter []byte
	UseBase64 bool
	BufSize   int
}

type SendConstraint interface {
	string | []byte
}

// A Message is a protocol message.
type Message struct {
	Data       []byte
	Delimiter  []byte
	Headers    map[string][]string
	Body       []byte
	Use_Base64 bool
	Files      []MessageFile
	// Parsed    bool
	// Generated bool
}

// NewMessage creates a new Message.
func NewMessage(delimiter []byte, use_b64 bool) *Message {
	return &Message{
		Data:       []byte{},
		Delimiter:  delimiter,
		Headers:    make(map[string][]string),
		Body:       []byte{},
		Use_Base64: use_b64,
		Files:      make([]MessageFile, 0),
		// Parsed:    false,
		// Generated: false,
	}
}

func (m *Message) AddHeader(key string, value string) {
	_, ok := m.Headers[key]
	if ok {
		m.Headers[key] = append(m.Headers[key], value)
	} else {
		m.Headers[key] = []string{value}
	}
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

func (m *Message) AddFile(file MessageFile) {
	m.Files = append(m.Files, file)
}

func (m *Message) AddRawFile(Name string, Data []byte) {
	m.Files = append(m.Files, MessageFile{Name: Name, Data: Data})
}

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

// Splitting order
// 1. Split body and head
//	    a. Split head into key/value pairs
//	    b. Split key/value pairs into key and values
// 2. Split files from body
//	    a. Split files into file name and data

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
	// Check if the message has already been parsed
	// if m.Parsed || m.Generated {
	// 	return m, errors.New("message has already been parsed")
	// }
	header_delimiter := m.HeaderDelimiter()
	body_delimiter := m.BodyDelimiter()
	file_delimiter := m.FileDelimiter()
	ending_delimiter := m.EndingDelimiter()
	// Split data into headers and body
	datalist := bytes.SplitN(m.Data, body_delimiter, 2)
	if len(datalist) != 2 {
		return nil, errors.New("invalid message sent")
	}
	// Get headers from m.Data
	headers := bytes.Split(datalist[0], header_delimiter)
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, header := range headers {
		wg.Add(1)
		// Start goroutine for each header
		go func(header []byte, wg *sync.WaitGroup, mu *sync.Mutex) {
			defer wg.Done()
			// Split header into key and values
			head := bytes.Split(header, m.Delimiter)
			str_list := make([]string, 0)
			for _, byt := range head[1:] {
				str_list = append(str_list, string(byt))
			}
			mu.Lock()
			m.Headers[string(head[0])] = str_list
			mu.Unlock()
		}(header, &wg, &mu)
	}
	// Get body from m.Data
	// Decode base64 encoded body
	var body []byte
	var err error
	full_body := bytes.Trim(datalist[1], string(ending_delimiter))
	if m.Use_Base64 {
		buf := bytes.NewBuffer(full_body)
		decoder := base64.NewDecoder(base64.StdEncoding, buf)
		full_body, err = io.ReadAll(decoder)
		if err != nil {
			return nil, err
		}
		// full_body, err = base64.StdEncoding.DecodeString(string(full_body))
		// if err != nil {
		// return nil, err
		// }
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
			file_data := bytes.Split(file, m.Delimiter)
			if len(file_data) != 2 {
				return
			}
			mu.Lock()
			m.Files = append(m.Files, MessageFile{Name: string(file_data[0]), Data: file_data[1]})
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
	// Check if the message has already been generated
	// if m.Generated || m.Parsed {
	// 	return m, errors.New("message has already been generated")
	// }
	var buffer bytes.Buffer
	var wg sync.WaitGroup
	var mu sync.Mutex

	header_delimiter := append(m.Delimiter, m.Delimiter...)
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
		go func(file MessageFile, buffer *bytes.Buffer, wg *sync.WaitGroup, mu *sync.Mutex) {
			defer wg.Done()
			// Create buffer for length of current file line
			total_len := len(file.Name) + len(m.HeaderDelimiter()) + len(file.Data) + len(m.FileDelimiter())
			// Create fileline
			fileline := make([]byte, total_len)
			// Copy name to fileline
			copy(fileline, file.Name)
			copy(fileline[len(file.Name):], m.HeaderDelimiter())
			// Copy data to fileline
			copy(fileline[len(file.Name)+len(m.HeaderDelimiter()):], file.Data)
			copy(fileline[len(file.Name)+len(m.HeaderDelimiter())+len(file.Data):], m.FileDelimiter())
			// Append fileline to buffer
			// Lock the mutex to prevent datarace
			mu.Lock()
			bodybuffer.Write(fileline)
			mu.Unlock()
		}(file, &buffer, &wg, &mu)
	}
	// Wait for all goroutines to finish
	wg.Wait()
	// Append body to buffer
	bodybuffer.Write(m.Body)
	// Encode body
	if m.Use_Base64 {
		buf := make([]byte, base64.StdEncoding.EncodedLen(bodybuffer.Len()))
		base64.StdEncoding.Encode(buf, bodybuffer.Bytes())
		buffer.Write(buf)
	} else {
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
