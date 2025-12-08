package serialization

import (
	"os"
	"encoding/json"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/mapLogic"
)

func JSONToMap(data []byte) (*mapLogic.Map, error){
	var tileMap mapLogic.Map
	err := json.Unmarshal(data, &tileMap)
	if err != nil{
		return nil, err
	}
	return &tileMap, nil
}

func LoadMapFromJSONFile(clientType, filename string) (*mapLogic.Map, error){
	filePath := "./" + clientType + "Client/map/" + filename
	data, err := os.ReadFile(filePath)
	if err != nil{
		return nil, err
	}

	tileMap, err := JSONToMap(data)
	if err != nil{
		return nil, err
	}

	return tileMap, nil
}