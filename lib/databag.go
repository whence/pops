// Implement encryption/decrytion of Chef encrypted data bag V1
// Super thanks to https://github.com/dgryski/dkeyczar for showing me how to use Go's crypto functions.

package lib

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

const (
	dataBagCipherAlgorithm = "aes-256-cbc"
	dataBagVersion         = 1
)

type encryptedDataBagItem struct {
	ID      string
	Entries map[string]*encryptedDataBagEntry
}

type encryptedDataBagEntry struct {
	Cipher        string  `json:"cipher"`
	EncryptedData string  `json:"encrypted_data"`
	Iv            string  `json:"iv"`
	Version       float64 `json:"version"`
}

type dataBagItem struct {
	ID      string
	Entries map[string]string
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

// Encrypt a databag item file
func Encrypt(itemPath, secretPath string) (string, error) {
	item := newDataBagItem(readFile(itemPath))

	itemFileName := filepath.Base(itemPath)
	itemExtName := filepath.Ext(itemPath)
	itemBaseName := itemFileName[:len(itemFileName)-len(itemExtName)]
	if item.ID != itemBaseName {
		return "", fmt.Errorf("Filename of %s does not match the ID: %s", itemFileName, item.ID)
	}

	secretData := readFile(secretPath)
	entries := item.encrypt(secretData)
	bytes, e := json.MarshalIndent(entries, "", "  ")
	if e != nil {
		panic("Failed to marshal data bag item")
	}
	return string(bytes), nil
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

func newDataBagItem(raw []byte) *dataBagItem {
	var kvs map[string]interface{}
	if json.Unmarshal(raw, &kvs) != nil {
		panic("Failed to unmarshal data bag item")
	}

	item := new(dataBagItem)
	item.Entries = make(map[string]string)

	for k, v := range kvs {
		switch k {
		case "id":
			item.ID = v.(string)
		default:
			item.Entries[k] = v.(string)
		}
	}

	return item
}

func (encryptedItem *encryptedDataBagItem) decrypt(secretData []byte) map[string]string {
	entries := make(map[string]string, len(encryptedItem.Entries)+1)
	entries["id"] = encryptedItem.ID

	for key, entry := range encryptedItem.Entries {
		if entry.Version != dataBagVersion {
			panic(fmt.Sprintf("Not implemented for encrypted bag version %f", entry.Version))
		}
		if entry.Cipher != dataBagCipherAlgorithm {
			panic(fmt.Sprintf("Not implemented for encrypted bag cipher %s", entry.Cipher))
		}

		entries[key] = entry.decrypt(secretData)
	}

	return entries
}

func (item *dataBagItem) encrypt(secretData []byte) map[string]interface{} {
	entries := make(map[string]interface{}, len(item.Entries)+1)
	entries["id"] = item.ID

	for key, entry := range item.Entries {
		iv := RandomBytes(aes.BlockSize)
		entries[key] = &encryptedDataBagEntry{
			Cipher:        dataBagCipherAlgorithm,
			Version:       dataBagVersion,
			Iv:            EncodeBase64(iv) + "\n",
			EncryptedData: encryptData(entry, iv, secretData) + "\n",
		}
	}

	return entries
}

func (entry *encryptedDataBagEntry) decrypt(secretData []byte) string {
	ciphertext := DecodeBase64(entry.EncryptedData)
	if len(ciphertext)%aes.BlockSize != 0 {
		panic("Ciphertext is not a multiple of the block size")
	}

	iv := DecodeBase64(entry.Iv)
	block := newCipher(secretData)
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)
	ciphertext = unPKCS7Padding(ciphertext)

	var wrapper version1Wrapper
	if err := json.Unmarshal(ciphertext, &wrapper); err != nil {
		panic(fmt.Sprintf("Failed to unmarshal data bag content. %+v", err))
	}

	return wrapper.Content
}

func encryptData(data string, iv, secretData []byte) string {
	block := newCipher(secretData)
	mode := cipher.NewCBCEncrypter(block, iv)

	wrapper := version1Wrapper{
		Content: data,
	}
	ciphertext, err := json.Marshal(wrapper)
	if err != nil {
		panic(fmt.Sprintf("Failed to marshal data bag content. %+v", err))
	}

	ciphertext = pkcs7padding(ciphertext, aes.BlockSize)
	cipherOut := make([]byte, len(ciphertext))
	mode.CryptBlocks(cipherOut, ciphertext)

	return EncodeBase64(cipherOut)
}

func newCipher(secretData []byte) cipher.Block {
	key := sha256.Sum256(secretData)
	block, err := aes.NewCipher(key[:])
	if err != nil {
		panic(fmt.Sprintf("Failed to create cipher. %+v", err))
	}
	return block
}

func unPKCS7Padding(data []byte) []byte {
	dataLen := len(data)
	endIndex := int(data[dataLen-1])
	// no need to protect 16 byte. http://stackoverflow.com/questions/7447242/how-does-pkcs7-not-lose-data
	return data[:dataLen-endIndex]
}

func pkcs7padding(data []byte, blocksize int) []byte {
	pad := blocksize - len(data)%blocksize
	b := make([]byte, pad, pad)
	for i := 0; i < pad; i++ {
		b[i] = uint8(pad)
	}
	return append(data, b...)
}
