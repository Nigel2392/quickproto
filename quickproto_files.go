package quickproto

import (
	"os"
	"path/filepath"
)

// Message file used for storing files sent inside the Message struct
type messageFile struct {
	Name string
	Data []byte
}

// NewmessageFile creates a new messageFile.
func NewmessageFile(name string, data []byte) messageFile {
	// Get file mime type
	return messageFile{
		Name: name,
		Data: data,
	}
}

// Size returns the size of the file.
func (f *messageFile) Size() int {
	return len(f.Data)
}

// String returns the name of the file.
func (f *messageFile) String() string {
	return f.Name
}

// Save the file to a path.
func (f *messageFile) Save(path string) error {
	return os.WriteFile(filepath.Join(path, f.Name), f.Data, 0644)
}
