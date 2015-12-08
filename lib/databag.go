package lib

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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
func Encrypt(itemPath, secretPath string) string {
	item := newDataBagItem(readFile(itemPath))
	secretData := readFile(secretPath)
	entries := item.encrypt(secretData)
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
		iv := randomIv()
		entries[key] = &encryptedDataBagEntry{
			Cipher:        dataBagCipherAlgorithm,
			Version:       dataBagVersion,
			Iv:            encodeBase64(iv) + "\n",
			EncryptedData: encryptData(entry, iv, secretData) + "\n",
		}
	}

	return entries
}

func (entry *encryptedDataBagEntry) decrypt(secretData []byte) string {
	ciphertext := decodeBase64(entry.EncryptedData)
	if len(ciphertext)%aes.BlockSize != 0 {
		panic("Ciphertext is not a multiple of the block size")
	}

	iv := decodeBase64(entry.Iv)
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

	return encodeBase64(cipherOut)
}

func newCipher(secretData []byte) cipher.Block {
	key := sha256.Sum256(secretData)
	block, err := aes.NewCipher(key[:])
	if err != nil {
		panic(fmt.Sprintf("Failed to create cipher. %+v", err))
	}
	return block
}

func randomIv() []byte {
	iv := make([]byte, aes.BlockSize)
	io.ReadFull(rand.Reader, iv)
	return iv
}

func decodeBase64(str string) []byte {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		panic(fmt.Sprintf("Decode base64 error: %+v", err))
	}
	return data
}

func encodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
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
