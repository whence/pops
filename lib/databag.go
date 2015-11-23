package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type dataBagItem struct {
	ID      string
	Entries map[string]*dataBagEntry
}

type dataBagEntry struct {
	Cipher        string
	EncryptedData string
	Iv            string
	Version       float64
}

// Decrypt a databag item file
func Decrypt(path string) string {
	item := newDataBagItem(path)
	return fmt.Sprintf("Results: %+v\n", item)
}

func newDataBagItem(path string) *dataBagItem {
	file, e := ioutil.ReadFile(path)
	if e != nil {
		panic(fmt.Sprintf("File error: %v\n", e))
	}

	var kvs map[string]interface{}
	json.Unmarshal(file, &kvs)

	item := new(dataBagItem)
	item.Entries = make(map[string]*dataBagEntry)

	for k, v := range kvs {
		switch k {
		case "id":
			item.ID = v.(string)
		default:
			entry := v.(map[string]interface{})
			item.Entries[k] = &dataBagEntry{
				Cipher:        entry["cipher"].(string),
				EncryptedData: entry["encrypted_data"].(string),
				Iv:            entry["iv"].(string),
				Version:       entry["version"].(float64),
			}
		}
	}

	return item
}

// func (e *EncryptedDataBagItem) DecryptKey(keyName, secret string) interface{} {
// 	key := e.Keys[keyName]
//
// 	switch key.Version {
// 	case 1:
// 		return version1Decoder([]byte(secret), key.Iv, key.EncryptedData)
// 	default:
// 		panic(fmt.Sprintf("not implemented for encrypted bag version %d!", key.Version))
// 	}
// }
