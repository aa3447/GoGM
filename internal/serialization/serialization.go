package serialization

import (
	"fmt"
	"os"
	"encoding/json"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/mapLogic"
)

func MapToJSON(tileMap *mapLogic.Map) ([]byte,error){
	mapJson, err := json.MarshalIndent(*tileMap, "", "  ")
	if err != nil{
		return []byte{},err
	}

	return mapJson,nil
}

func SaveMapToFile(tileMap *mapLogic.Map, clientName, filename string) error{
	mapJson, err := MapToJSON(tileMap)
	if err != nil{
		return err
	}

	filePath := fmt.Sprintf("./" + clientName + "Client/map/%s.json", filename)
	err = os.WriteFile(filePath, mapJson, 0644)
	if err != nil{
		return err
	}

	return nil
}