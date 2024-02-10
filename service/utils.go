package service

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

const JsonIndent = "  "

// PrettyPrintBytes  returns a pretty-printed version of the input byte slice. If an error occurs during JSON parsing
func PrettyPrintBytes(data interface{}) bytes.Buffer {
	jsonObj, err := json.Marshal(data)
	if err != nil {
		return bytes.Buffer{}
	}

	var out bytes.Buffer
	err = json.Indent(&out, jsonObj, "", JsonIndent)
	if err != nil {
		return bytes.Buffer{}
	}

	return out
}

// PpObj  is a convenience function for pretty printing an object to stdout.
//
//	It's intended to be used in development and debugging.
func PpObj(data interface{}) string {
	b := PrettyPrintBytes(data)
	return b.String()
}

// FileType  returns the type of file referenced by path.
func FileType(filename string) (string, error) {
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		return "", err
	}
	buf := make([]byte, 512)
	_, err = f.Read(buf)
	if err != nil {
		return "", err
	}

	contentType := http.DetectContentType(buf)

	return contentType, nil
}
