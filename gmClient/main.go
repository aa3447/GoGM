package main

import (
	"fmt"
	"maps"
	"slices"

	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/campaign"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/io"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/mapLogic"
	GM "home/aa3447/workspace/github.com/aa3447/GoGM/internal/playerLogic"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/pubsub"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/serialization"

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

	err = pubsub.SetupExchanges()
	if err != nil{
		fmt.Println("Error starting pubsub:", err)
		return
	}

	campaign := campaign.NewCampaign("Test Campaign", "A test campaign for GoGM", GM.NewGM("GameMaster", "The overseer of the game world"))

	cMap , err := campaign.NewGamestateWithRandomMap(10, 10, 0.2, []float64{0.5, 0.2, 0.2, 0.1},"")
	if err != nil{
		fmt.Println("Error creating game state:", err)
		return
	}

	err = serialization.SaveToFile(*cMap, "gm" ,"map" ,cMap.Name)
	if err != nil{
		fmt.Println("Error creating game state:", err)
		return
	}

	players := []*GM.Player{}
	testPlayer := GM.NewPlayer("Hero", "The brave adventurer", "Warrior", "assign", []int{15,14,13,12,10,8})
	testPlayer.SetLocation(cMap.EntranceLocation[0] ,cMap.EntranceLocation[1])
	players = append(players, testPlayer)
	go playerMoveSubscriber(players, channel, cMap)

	gm := GM.NewGM("GameMaster", "The overseer of the game world")
	
	fmt.Println("Welcome,", gm.Name)
	fmt.Println("You are the Game Master.")
	fmt.Println("Type 'map' to view the map, or 'quit' to exit.")

	commands := io.GetInput()
	for {
		switch commands[0] {
			case "load":
				// Implement load logic here
			case "campaign":
				// Implement campaign logic here
			case "play":
				gameLoop(cMap, channel)
			case "quit":
				fmt.Println("Quitting game.")
				return
			default:
				fmt.Println("Unknown command:", commands[0])
		}
		commands = io.GetInput()			
	}

}

func campaignManagementLoop(channel *ampq.Channel) (*campaign.Campaign, error){
	fmt.Println("Entering campaign management loop. Type 'quit' to exit.")
	currentCampaign := campaign.NewCampaign("GM Campaign", "A GM campaign for GoGM", GM.NewGM("GameMaster", "The overseer of the game world"))

	commands := io.GetInput()
	for {
		command := commands[0]
		args := commands[1:]
		switch command {
			case "create_map":
				switch args[0] {
				case "random":
					newMap, err := mapLogic.GenRandomMap(10, 10, 0.2, []float64{0.5, 0.2, 0.2, 0.1},"")
					if err != nil {
						fmt.Println("Error creating random map:", err)
						continue
					}
					currentCampaign.AddMap(&newMap)
					fmt.Println("Random map created and added to campaign.")
				case "editor":
					// Implement map editor here
				default:
					fmt.Println("Unknown create_map argument:", args)
				}
			case "update_map":
				// Implement map update logic here
			case "load_map":
				// Implement map loading logic here
			case "add_player":
				// Implement player addition logic here
			case "remove_player":
				// Implement player removal logic here
			case "add_npc":
				// Implement NPC addition logic here
			case "remove_npc":
				// Implement NPC removal logic here
			case "quit":
				fmt.Println("Exiting campaign management loop.")
				return currentCampaign, nil
			default:
				fmt.Println("Unknown command:", command)
		}
		commands = io.GetInput()
	}
}

func gameLoop(m *mapLogic.Map, channel *ampq.Channel){
	commands := io.GetInput()
	currentPlayers := slices.Collect(maps.Values(m.GameState.Players))
	for {
		command := commands[0]
		switch command {
			case "map":
				m.PrintMapDebugWithPlayers(currentPlayers)
			case "send":
				if len(commands) < 2{
					fmt.Println("Specify what to send: 'map'")
					commands = io.GetInput()
					continue
				}
				switch commands[1]{
					case "map":
						pubsub.PublishToQueueAsJSON(channel, pubsub.MapExchange, pubsub.MapNewRoutingKey, m)
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
						pubsub.PublishToQueueAsJSON(channel, pubsub.MapExchange, pubsub.MapUpdateRoutingKey, m)
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

func playerMoveSubscriber(players []*GM.Player,channel *ampq.Channel, m *mapLogic.Map) {
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
				_ ,err = m.MovePlayer(player, playerMove.To[0]  - playerMove.From[0], playerMove.To[1] -  playerMove.From[1])
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