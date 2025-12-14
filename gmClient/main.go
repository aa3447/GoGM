package main

import (
	"fmt"

	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/gameLogic"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/io"
	GM "home/aa3447/workspace/github.com/aa3447/GoGM/internal/playerLogic"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/serialization"
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

	gameState , err := gameLogic.NewGamestateWithRandomMap(10, 10, 0.2, []float64{0.5, 0.2, 0.2, 0.1},"")
	if err != nil{
		fmt.Println("Error creating game state:", err)
		return
	}
	err = pubsub.SetupExchanges()
	if err != nil{
		fmt.Println("Error starting pubsub:", err)
		return
	}
	currentMap := gameState.CurrentMap
	
	err = serialization.SaveToFile(*currentMap, "gm" ,"map" ,currentMap.Name)
	if err != nil{
		fmt.Println("Error creating game state:", err)
		return
	}

	players := []*GM.Player{}
	testPlayer := GM.NewPlayer("Hero", "The brave adventurer", "Warrior", "assign", []int{15,14,13,12,10,8})
	testPlayer.SetLocation(gameState.CurrentMap.EntranceLocation[0] ,gameState.CurrentMap.EntranceLocation[1])
	players = append(players, testPlayer)
	go playerMoveSubscriber(players, channel, gameState)

	gm := GM.NewGM("GameMaster", "The overseer of the game world")
	
	fmt.Println("Welcome,", gm.Name)
	fmt.Println("You are the Game Master.")
	fmt.Println("Type 'map' to view the map, or 'quit' to exit.")

	commands := io.GetInput()
	for {
		command := commands[0]
		switch command {
			case "map":
				currentMap.PrintMapDebugWithPlayers(players)
			case "send":
				if len(commands) < 2{
					fmt.Println("Specify what to send: 'map'")
					commands = io.GetInput()
					continue
				}
				switch commands[1]{
					case "map":
						pubsub.PublishToQueueAsJSON(channel, pubsub.MapExchange, pubsub.MapNewRoutingKey, currentMap)
						fmt.Println("Map sent to players.")
					default:
						fmt.Println("Unknown send command:", commands[0])
						commands = io.GetInput()
						continue
				}
			case "update":
				if len(commands) < 2{
					fmt.Println("Specify what to update: 'map'")
					commands = io.GetInput()
					continue
				}
				switch commands[1]{
					case "map":
						pubsub.PublishToQueueAsJSON(channel, pubsub.MapExchange, pubsub.MapUpdateRoutingKey, currentMap)
						fmt.Println("Map update sent to players.")
					default:
						fmt.Println("Unknown update command:", commands[0])
						commands = io.GetInput()
						continue
				}
			case "quit":
				fmt.Println("Quitting game.")
				return
			default:
				fmt.Println("Unknown command:", command)
		}
		commands = io.GetInput()
	}

}

func playerMoveSubscriber(players []*GM.Player,channel *ampq.Channel, gameState *gameLogic.Gamestate) {
	subscribeArgs := []bool{true, false, false, false}
	msgs, err := pubsub.SubscribeToQueue(channel, pubsub.PlayerMoveQueue, subscribeArgs, subscribeArgs)
	if err != nil{
		fmt.Println("Error subscribing to player move queue:", err)
		return
	}

	for d := range msgs{
		playerMove ,err := serialization.JSONTo(d.Body, mapLogic.PlayerMove{})
		if err != nil{
			fmt.Println("Error deserializing player move:", err)
			continue
		}
		for _, player := range players{
			if player.Name == playerMove.PlayerName{
				_ ,err = gameState.MovePlayer(player, playerMove.To[0]  - playerMove.From[0], playerMove.To[1] -  playerMove.From[1])
				if err != nil{
					fmt.Println("Error applying player move:", err)
					continue
				}
				fmt.Printf("Processed move for player %s to (%d, %d)\n", playerMove.PlayerName, playerMove.To[0], playerMove.To[1])
				pubsub.PublishToQueueAsJSON(channel, pubsub.MoveExchange, pubsub.GMMoveRoutingKey, playerMove)
			}
		}
	}
}