package tests

import (
	"testing"

	"github.com/Nigel2392/quickproto"
)

func TestMessageFiles(t *testing.T) {
	// Create a new message
	msg := quickproto.NewMessage([]byte("&"), false, quickproto.Base16Encoding, quickproto.Base16Decoding)
	msg.Headers["key1"] = []string{"value1", "value2", "value3"}
	msg.Headers["key2"] = []string{"value2"}
	msg.Body = []byte("BODYBODYBODY")

	fdata := "!@#$%^&*()$$$$$$$$FILE"
	msg.AddRawFile("file1", []byte(fdata+"1"))
	msg.AddRawFile("file2", []byte(fdata+"2"))
	msg.AddRawFile("file3", []byte(fdata+"3"))

	msg.Generate()

	newmsg := quickproto.NewMessage([]byte("&"), false, quickproto.Base16Encoding, quickproto.Base16Decoding)
	newmsg.Data = msg.Data
	newmsg.Parse()

	if newmsg.Files["file1"].Name == "" {
		t.Error("Expected file1 to be not nil")
	}
	if newmsg.Files["file2"].Name == "" {
		t.Error("Expected file2 to be not nil")
	}
	if newmsg.Files["file3"].Name == "" {
		t.Error("Expected file3 to be not nil")
	}
	if string(newmsg.Files["file1"].Data) != fdata+"1" {
		t.Error("Expected file1 to be FILE1")
	}
	if string(newmsg.Files["file2"].Data) != fdata+"2" {
		t.Error("Expected file2 to be FILE2")
	}
	if string(newmsg.Files["file3"].Data) != fdata+"3" {
		t.Error("Expected file3 to be FILE3")
	}
	if string(newmsg.Body) != "BODYBODYBODY" {
		t.Error("Expected body to be BODYBODYBODY")
	}
	if newmsg.Headers["key1"][0] != "value1" {
		t.Error("Expected key1 to be value1")
	}
	if newmsg.Headers["key1"][1] != "value2" {
		t.Error("Expected key1 to be value2")
	}
	if newmsg.Headers["key1"][2] != "value3" {
		t.Error("Expected key1 to be value3")
	}
	if newmsg.Headers["key2"][0] != "value2" {
		t.Error("Expected key2 to be value2")
	}
}
