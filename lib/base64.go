package lib

import (
	"encoding/base64"
	"fmt"
)

// DecodeBase64 decodes base64 string to bytes
func DecodeBase64(str string) []byte {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		panic(fmt.Sprintf("Decode base64 error: %+v", err))
	}
	return data
}

// EncodeBase64 encodes bytes to base64 string
func EncodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}
