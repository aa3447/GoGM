package main

import (
	"fmt"

	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/gameLogic"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/io"
	GM "home/aa3447/workspace/github.com/aa3447/GoGM/internal/playerLogic"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/serialization"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/pubsub"
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

	gameState , _ ,err := gameLogic.NewGamestateWithRandomMap(10, 10, 0.2, []float64{0.5, 0.2, 0.2, 0.1})
	if err != nil{
		fmt.Println("Error creating game state:", err)
		return
	}
	err = pubsub.Start()
	if err != nil{
		fmt.Println("Error starting pubsub:", err)
		return
	}
	currentMap := gameState.CurrentMap
	err = serialization.SaveMapToFile(currentMap, currentMap.Name)
	if err != nil{
		fmt.Println("Error creating game state:", err)
		return
	}
	pubsub.PublishMapToQueue(channel, pubsub.MapExchange, pubsub.MapRoutingKey,currentMap)

	gm := GM.NewGM("GameMaster", "The overseer of the game world")
	
	fmt.Println("Welcome,", gm.Name)
	fmt.Println("You are the Game Master.")
	fmt.Println("Type 'map' to view the map, or 'quit' to exit.")

	commands := io.GetInput()
	for {
		command := commands[0]
		switch command {
			case "map":
				currentMap.PrintMapDebug()
			case "quit":
				fmt.Println("Quitting game.")
				return
			default:
				fmt.Println("Unknown command:", command)
		}
		commands = io.GetInput()
	}

}