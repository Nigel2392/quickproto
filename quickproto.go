package quickproto

import (
	"bytes"
	"errors"
	"strings"
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

// Standard delimiter.
var STANDARD_DELIM []byte = []byte("$")

// These are tested not to work.
// Byte "\x00" is used as a body, when body is empty.
var BANNED_DELIMITERS = []string{
	"=", "_", "\x08", "\x1e", "\x00", "(", ")",
	"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
	"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
}

// A Message is a protocol message.
type Message struct {
	Data        []byte
	Delimiter   []byte
	Headers     map[string][]string
	Body        []byte
	Files       map[string]*messageFile
	UseEncoding bool
	Encode_func func([]byte) []byte
	Decode_func func([]byte) ([]byte, error)
	F_Encoder   func([]byte) []byte
	F_Decoder   func([]byte) ([]byte, error)
}

// NewMessage creates a new Message.
func NewMessage(delimiter []byte, useencoding bool, encode_func func([]byte) []byte, decode_func func([]byte) ([]byte, error)) *Message {
	if delimiter == nil {
		delimiter = STANDARD_DELIM
	}
	return &Message{
		Data:        []byte{},
		Delimiter:   delimiter,
		Headers:     make(map[string][]string),
		Body:        []byte{},
		Files:       make(map[string]*messageFile),
		UseEncoding: useencoding,
		Encode_func: encode_func,
		Decode_func: decode_func,
		F_Encoder:   Base64Encoding,
		F_Decoder:   Base64Decoding,
	}
}

// Add a header to the message.
func (m *Message) AddHeader(key string, value string) error {
	if strings.Contains(key, string(m.Delimiter)) {
		return errors.New("header key cannot contain delimiter")
	}
	if strings.Contains(value, string(m.Delimiter)) {
		return errors.New("header value cannot contain delimiter")
	}
	m.Headers[key] = append(m.Headers[key], value)
	return nil
}

// Add content to the message.
// Either add []bytes, or a string.
func (m *Message) AddContent(content any) error {
	switch content := content.(type) {
	case string:
		m.Body = append(m.Body, content...)
	case []byte:
		m.Body = append(m.Body, content...)
	case *messageFile:
		m.Files[content.Name] = content
	default:
		return errors.New("invalid content type")
	}
	return nil
}

// Add a MessageFile to the message.
func (m *Message) AddFile(file *messageFile) {
	m.Files[file.Name] = file
}

// Create a MessageFile, and add it to the message.
func (m *Message) AddRawFile(name string, data []byte) {
	m.Files[name] = &messageFile{Name: name, Data: data}
}

// Header delimiter, returns DELIMITER + DELIMITER
func (m *Message) HeaderDelimiter() []byte {
	return append(m.Delimiter, m.Delimiter...)
}

// Body delimiter, returns HEADER_DELIMITER + HEADER_DELIMITER
func (m *Message) BodyDelimiter() []byte {
	return append(m.HeaderDelimiter(), m.HeaderDelimiter()...)
}

// File delimiter, returns BODY_DELIMITER + HEADER_DELIMITER
func (m *Message) FileDelimiter() []byte {
	return append(m.BodyDelimiter(), m.HeaderDelimiter()...)
}

// End delimiter, returns BODY_DELIMITER + BODY_DELIMITER
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
	// Split headers into key/value pairs
	headers := bytes.Split(datalist[0], header_delimiter)
	for _, header := range headers {
		head := bytes.Split(header, m.Delimiter)
		str_list := make([]string, 0)
		// Set multiple values for each key
		for _, byt := range head[1:] {
			str_list = append(str_list, string(byt))
		}
		// Set key and values, lock for thread safety
		m.Headers[string(head[0])] = str_list
	}
	// Decode base64 encoded body
	var body []byte
	var err error
	full_body := bytes.Trim(datalist[1], string(ending_delimiter))
	if m.Encode_func != nil && m.Decode_func != nil && m.UseEncoding {
		full_body, err = m.Decode_func(full_body)
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
	for _, file := range body_data {
		file_data_list := bytes.SplitN(file, header_delimiter, 2)
		file_name := string(file_data_list[0])
		// File data should always be encoded!
		file_data, err := m.F_Decoder(file_data_list[1])
		if err != nil {
			return nil, err
		}
		mf := NewmessageFile(file_name, file_data)
		m.Files[file_name] = &mf
	}
	if len(body) != 1 && body[0] != 0x00 {
		m.Body = body
	}
	// m.Parsed = true
	return m, nil
}

// creates a protocol message.
// Header is a map of key/value pairs.
// Body is a base64 encoded byte slice.
func (m *Message) Generate() (*Message, error) {
	var buffer bytes.Buffer
	// Predefine delimiters so we dont have to calculate them every time.
	var (
		header_delimiter = m.HeaderDelimiter()
		file_delimiter   = m.FileDelimiter()
		LenDelim         = len(m.Delimiter)
		lenHDelim        = LenDelim * 2
		lenFDelim        = lenHDelim * 3
	)
	// Create headers
	for key, value := range m.Headers {
		// Start goroutine for each header
		// Create buffer for length of current header line
		total_len := 0
		for _, str := range value {
			// Append key and value to headerline
			val_len := len(str)
			total_len = total_len + val_len + LenDelim
		}
		// Create headerline
		headerline := make([]byte, len(key)+LenDelim+total_len+LenDelim)
		// Copy key to headerline
		copy(headerline, key)
		copy(headerline[len(key):], m.Delimiter)
		// Copy values to headerline
		current_pos := len(key) + LenDelim
		for _, str := range value {
			copy(headerline[current_pos:], str)
			copy(headerline[current_pos+len(str):], m.Delimiter)
			current_pos = current_pos + len(str) + LenDelim
		}
		// Set last delimiter
		copy(headerline[current_pos:], m.Delimiter)
		// Append headerline to buffer
		// Lock the mutex to prevent datarace
		buffer.Write(headerline)
	}
	buffer.Write(header_delimiter)
	// Get files
	var bodybuffer bytes.Buffer
	for _, file := range m.Files {
		// Write the file to the body
		// Create buffer for length of current file line
		// Encode file data
		fdata := m.F_Encoder(file.Data)
		// Get size of buffer for all file data and delimiters
		total_len := len(file.Name) + lenHDelim + len(fdata) + lenFDelim
		fileline := make([]byte, total_len)
		// Copy data to fileline
		copy(fileline, file.Name)
		copy(fileline[len(file.Name):], header_delimiter)
		copy(fileline[len(file.Name)+lenHDelim:], fdata)
		copy(fileline[len(file.Name)+lenHDelim+len(fdata):], file_delimiter)
		// Append fileline to buffer
		bodybuffer.Write(fileline)
	}
	// Wait for all goroutines to finish
	// Create body

	// Append body to buffer
	bodybuffer.Write(m.Body)
	// Write a NULL byte if body is empty.
	// This is to prevent one of the files ending up as the body, when no body is provided.
	if len(m.Body) == 0 {
		bodybuffer.Write([]byte{0x00})
	}
	if m.Encode_func != nil && m.Decode_func != nil && m.UseEncoding {
		// If encoding is set, create buffer and encode body
		buffer.Write(m.Encode_func(bodybuffer.Bytes()))
	} else {
		// Else write body to buffer
		buffer.Write(bodybuffer.Bytes())
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

func (m *Message) FileSizes() map[string]int {
	files := make(map[string]int)
	for _, file := range m.Files {
		files[file.Name] = file.Size()
	}
	return files
}
