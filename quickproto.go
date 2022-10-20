package quickproto

import (
	"bytes"
	"encoding/base64"
	"errors"
	"sync"
)

var STANDARD_DELIM []byte = []byte("$")

type Config struct {
	Delimiter []byte
	UseBase64 bool
	BufSize   int
}

// A Message is a protocol message.
type Message struct {
	Data       []byte
	Delimiter  []byte
	Headers    map[string][]string
	Body       []byte
	Use_Base64 bool
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

func (m *Message) HeaderDelimiter() []byte {
	return append(m.Delimiter, m.Delimiter...)
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
	datalist[1] = bytes.Trim(datalist[1], string(ending_delimiter))
	if m.Use_Base64 {
		b64 := datalist[1]
		body, err = base64.StdEncoding.DecodeString(string(b64))
		if err != nil {
			return nil, errors.New("invalid base64 encoded body")
		}
	} else {
		body = datalist[1]
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
	if m.Use_Base64 {
		b64 := base64.StdEncoding.EncodeToString(m.Body)
		buffer.WriteString(b64)
	} else {
		buffer.Write(m.Body)
	}
	buffer.Write(m.EndingDelimiter())
	m.Data = buffer.Bytes()
	// m.Generated = true
	return m, nil
}

// Get content length of the message.
func (m *Message) ContentLength() int {
	return len(m.Data)
}
