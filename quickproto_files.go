package quickproto

import (
	"os"
	"path/filepath"
)

// Message file used for storing files sent inside the Message struct
type MessageFile struct {
	Name string
	Data []byte
}

// NewMessageFile creates a new MessageFile.
func NewMessageFile(name string, data []byte) MessageFile {
	return MessageFile{
		Name: name,
		Data: data,
	}
}

// Size returns the size of the file.
func (f *MessageFile) Size() int {
	return len(f.Data)
}

// String returns the name of the file.
func (f *MessageFile) String() string {
	return f.Name
}

// Save the file to a path.
func (f *MessageFile) Save(path string) error {
	return os.WriteFile(filepath.Join(path, f.Name), f.Data, 0644)
}
