package serialization

import (
	"fmt"
	"os"
	"encoding/json"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/mapLogic"
)

func mapToJSON(tileMap mapLogic.Map) ([]byte,error){
	mapJson, err := json.MarshalIndent(tileMap, "", "  ")
	if err != nil{
		return []byte{},err
	}

	return mapJson,nil
}

func SaveMapToFile(tileMap mapLogic.Map, filename string) error{
	mapJson, err := mapToJSON(tileMap)
	if err != nil{
		return err
	}

	filePath := fmt.Sprintf("./client/map/%s.json", filename)
	err = os.WriteFile(filePath, mapJson, 0644)
	if err != nil{
		return err
	}

	return nil
}