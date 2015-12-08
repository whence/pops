package lib

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type encryptedDataBagItem struct {
	ID      string
	Entries map[string]*encryptedDataBagEntry
}

type encryptedDataBagEntry struct {
	Cipher        string
	EncryptedData string
	Iv            string
	Version       float64
}

type version1Wrapper struct {
	Content string `json:"json_wrapper"`
}

// Decrypt a databag item file
func Decrypt(itemPath, secretPath string) string {
	encryptedItem := newEncryptedDataBagItem(readFile(itemPath))
	secretData := readFile(secretPath)
	entries := encryptedItem.decrypt(secretData)
	bytes, e := json.MarshalIndent(entries, "", "  ")
	if e != nil {
		panic("Failed to marshal data bag item")
	}
	return string(bytes)
}

func readFile(path string) []byte {
	content, e := ioutil.ReadFile(path)
	if e != nil {
		panic(fmt.Sprintf("File error: %v\n", e))
	}
	return content
}

func newEncryptedDataBagItem(raw []byte) *encryptedDataBagItem {
	var kvs map[string]interface{}
	if json.Unmarshal(raw, &kvs) != nil {
		panic("Failed to unmarshal data bag item")
	}

	item := new(encryptedDataBagItem)
	item.Entries = make(map[string]*encryptedDataBagEntry)

	for k, v := range kvs {
		switch k {
		case "id":
			item.ID = v.(string)
		default:
			entry := v.(map[string]interface{})
			item.Entries[k] = &encryptedDataBagEntry{
				Version:       entry["version"].(float64),
				Cipher:        entry["cipher"].(string),
				EncryptedData: entry["encrypted_data"].(string),
				Iv:            entry["iv"].(string),
			}
		}
	}

	return item
}

func (encryptedItem *encryptedDataBagItem) decrypt(secretData []byte) map[string]string {
	entries := make(map[string]string, len(encryptedItem.Entries)+1)
	entries["id"] = encryptedItem.ID

	for key, entry := range encryptedItem.Entries {
		if entry.Version != 1 {
			panic(fmt.Sprintf("Not implemented for encrypted bag version %f", entry.Version))
		}
		if entry.Cipher != "aes-256-cbc" {
			panic(fmt.Sprintf("Not implemented for encrypted bag cipher %s", entry.Cipher))
		}

		entries[key] = entry.decrypt(secretData)
	}

	return entries
}

// proudly stolen from https://github.com/go-chef/cryptobag/blob/master/decrypter_v1.go
func (entry *encryptedDataBagEntry) decrypt(secretData []byte) string {
	ciphertext := decodeBase64(entry.EncryptedData)
	initVector := decodeBase64(entry.Iv)
	keySha := sha256.Sum256(secretData)

	block, err := aes.NewCipher(keySha[:])
	if err != nil {
		fmt.Println(err)
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		panic("Ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, initVector)
	mode.CryptBlocks(ciphertext, ciphertext)

	ciphertext = unPKCS7Padding(ciphertext)

	var wrapper version1Wrapper
	if err := json.Unmarshal(ciphertext, &wrapper); err != nil {
		panic(fmt.Sprintf("Failed to unmarshal data bag content. %+v", err))
	}

	return wrapper.Content
}

func decodeBase64(str string) []byte {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		fmt.Println("error:", err)
	}
	return data
}

func unPKCS7Padding(data []byte) []byte {
	dataLen := len(data)
	endIndex := int(data[dataLen-1])
	// no need to protect 16 byte. http://stackoverflow.com/questions/7447242/how-does-pkcs7-not-lose-data
	return data[:dataLen-endIndex]
}
