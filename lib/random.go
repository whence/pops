package lib

import (
	"crypto/rand"
	"fmt"
)

// RandomBytes generates a random iv of size
func RandomBytes(size int) []byte {
	iv := make([]byte, size)
	_, err := rand.Read(iv)
	if err != nil {
		panic(fmt.Sprintf("Random Iv error: %+v", err))
	}
	return iv
}
