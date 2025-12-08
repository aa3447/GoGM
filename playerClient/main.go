package main

import (
	"fmt"
	"os"
	"errors"

	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/gameLogic"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/playerLogic"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/serialization"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/io"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/pubsub"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/mapLogic"


	ampq "github.com/rabbitmq/amqp091-go"
)

func main(){
	conn, err := ampq.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil{
		fmt.Println("failed to connect to RabbitMQ:", err)
		return
	}
	defer conn.Close()
	
	channel , err := conn.Channel()
	if err != nil{
		fmt.Println("failed to open a channel:", err)
		return
	}
	defer channel.Close()
	
	maps := mapQueueSubscriber(channel)

	gameState := gameLogic.NewGamestateWithExistingMap(<-maps)

	player := player.NewPlayer("Hero", "The brave adventurer", "Warrior")
	player.SetLocation(gameState.CurrentMap.EntranceLocation[0], gameState.CurrentMap.EntranceLocation[1])
	
	fmt.Println("Welcome,", player.Name)
	fmt.Println("You find yourself at the entrance of a mysterious location.")
	fmt.Println("Type 'move <direction>' to move (north, south, east, west), 'map' to view the map, or 'quit' to exit.")
	
	commands := io.GetInput()
	for {
		command := commands[0]
		args := commands[1:]
		switch command {
			case "move":
				handleMove(gameState,player,args)
			case "action":
				//handleAction(args)
			case "map":
				gameState.CurrentMap.PrintMapWithPlayer(player)
			case "quit":
				fmt.Println("Quitting game.")
				return
			default:
				fmt.Println("Unknown command:", command)
		}
		commands = io.GetInput()
	}
}


func handleMove(gs *gameLogic.Gamestate , player *player.Player ,args []string){
	if len(args) < 1{
		fmt.Println("Move where?")
		return
	}
	direction := args[0]
	fmt.Println("Moving", direction)

	switch direction{
		case "north", "up":
			gs.MovePlayer(player,-1, 0)
		case "south", "down":
			gs.MovePlayer(player,1, 0)
		case "east", "right":
			gs.MovePlayer(player,0, 1)
		case "west", "left":
			gs.MovePlayer(player,0, -1)
		default:
			fmt.Println("Unknown direction:", direction)
	}
}

func mapQueueSubscriber(channel *ampq.Channel) chan *mapLogic.Map{
	maps := make(chan *mapLogic.Map)

	msgs, err := pubsub.SubscribeToMapQueue(channel, pubsub.MapExchange, pubsub.MapRoutingKey)
	if err != nil{
		fmt.Println("Error subscribing to map queue:", err)
		return nil
	}

	go func(){
		for d := range msgs{
			currentMap, err := serialization.JSONToMap(d.Body)
			if err != nil{
				fmt.Println("Error deserializing map:", err)
				continue
			}
			_,err = os.Stat("./playerClient/map/" + currentMap.Name + ".json")
			if errors.Is(err, os.ErrNotExist) {
				err = serialization.SaveMapToFile(currentMap, "player", currentMap.Name)
				if err != nil{
					fmt.Println("Error saving map to file:", err)
					continue
				}
				currentMap.FileLocation = "./playerClient/" + currentMap.Name
			}
			maps <- currentMap
		}
	}()
	return maps
}