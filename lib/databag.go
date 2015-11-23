package lib

import (
  "fmt"
  "encoding/json"
  "io/ioutil"
)

type DataBagItem struct {
	Id string
	Entries map[string]*DataBagEntry
}

type DataBagEntry struct {
	Cipher string
	EncryptedData string
	Iv string
	Version float64
}

func NewDataBagItem(path string) *DataBagItem {
  file, e := ioutil.ReadFile(path)
  if e != nil {
    panic(fmt.Sprintf("File error: %v\n", e))
  }

  var kvs map[string]interface{}
  json.Unmarshal(file, &kvs)

  item := new(DataBagItem)
	item.Entries = make(map[string]*DataBagEntry)

  for k, v := range kvs {
		switch k {
		case "id":
			item.Id = v.(string)
		default:
      entry := v.(map[string]interface{})
			item.Entries[k] = &DataBagEntry{
    		Cipher: entry["cipher"].(string),
    		EncryptedData: entry["encrypted_data"].(string),
    		Iv: entry["iv"].(string),
    		Version: entry["version"].(float64),
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
