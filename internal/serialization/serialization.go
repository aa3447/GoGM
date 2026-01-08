package serialization

import (
	"os"
	"fmt"
	"encoding/json"
)


// ToJSON serializes the given data of type J to JSON.
func ToJSON[J JSONSafe](data J) ([]byte, error){
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil{
		return []byte{}, err
	}
	return jsonData, nil
}

// SaveToFile saves the given data of type J to a JSON file.
func SaveToFile[J JSONSafe](data J, clientType, filetype, filename string) error{
	mapJson, err := ToJSON(data)
	if err != nil{
		return err
	}

	filePath := fmt.Sprintf("./%sClient/%s/%s.json", clientType, filetype, filename)
	err = os.WriteFile(filePath, mapJson, 0644)
	if err != nil{
		return err
	}

	return nil
}
