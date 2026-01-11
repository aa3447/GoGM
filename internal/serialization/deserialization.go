package serialization

import (
	"encoding/json"
	"fmt"
	"os"

	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/campaign"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/mapLogic"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/equipment"
)

// JSONSafe is a constraint that includes types that can be serialized/deserialized to/from JSON.
type JSONSafe interface{
	mapLogic.Map | mapLogic.PlayerMove | campaign.Campaign | equipment.Weapon | equipment.Armor | equipment.Potion
}

// JSONTo deserializes JSON data into the specified type J.
func JSONTo[J JSONSafe](data []byte, dataType J) (*J, error){
	var unmarshaledData J
	err := json.Unmarshal(data, &unmarshaledData)
	if err != nil{
		return nil, err
	}
	return &unmarshaledData, nil
}

// LoadFromJSONFile loads JSON data from a file and deserializes it into the specified type J.
func LoadFromJSONFile[J JSONSafe](clientType, filetype, filename string, dataType J) (*J, error){
	filePath := fmt.Sprintf("./%sClient/%s/%s.json", clientType, filetype, filename)
	fileData, err := os.ReadFile(filePath)
	if err != nil{
		return nil, err
	}

	json, err := JSONTo(fileData, dataType)
	if err != nil{
		return nil, err
	}

	return json, nil
}