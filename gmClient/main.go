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

	var campaign *campaign.Campaign

	gm := GM.NewGM("GameMaster", "The overseer of the game world")
	
	fmt.Println("Welcome,", gm.Name)
	fmt.Println("You are the Game Master.")
	fmt.Println("Type 'map' to view the map, or 'quit' to exit.")

	commands := io.GetInput()
	for {
		switch commands[0] {
			case "load":
				
			case "campaign":
				campaign, err = campaignManagementLoop(channel)
				if err != nil{
					fmt.Println("Error in campaign management loop:", err)
					return
				}
			case "play":
				players := []*GM.Player{}
				// change to connected players for gamestate :todo:
				for _, player := range campaign.Players{
					players = append(players, player)
				}
				cMap := campaign.CurrentMap

				go playerMoveSubscriber(players, channel, cMap)
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

// campaignManagementLoop handles the campaign management commands for the GM.
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
			case "save_map":
				if len(args) == 0 {
					fmt.Println("Please specify the map name to save.")
					continue
				}
				if _, exists := currentCampaign.Maps[args[0]]; !exists {
					fmt.Println("Map not found in campaign:", args[0])
					continue
				}
				if currentCampaign.CurrentMap == nil  {
					fmt.Println("No current map set for the campaign.")
					continue
				}
				err := serialization.SaveToFile(*currentCampaign.Maps[args[0]], "gm", "map", args[0])
				if err != nil{
					fmt.Println("Error saving map:", err)
					continue
				}
				fmt.Println("Map saved:", args[0])
			case "current_map":
				if len(args) == 0 {
					fmt.Println("Please specify the map name to set as the current map.")
					continue
				}
				currentCampaign.CurrentMap = currentCampaign.Maps[args[0]]
			case "load_campaign":
				loadedCampaign, err := serialization.LoadFromJSONFile("gm", "campaigns", args[0], campaign.Campaign{})
				if err != nil{
					fmt.Println("Error loading campaign:", err)
					continue
				}
				currentCampaign = loadedCampaign
				fmt.Println("Campaign loaded:", currentCampaign.Name)
			case "save_campaign":
				err := serialization.SaveToFile(*currentCampaign, "gm", "campaigns", currentCampaign.Name)
				if err != nil{
					fmt.Println("Error saving campaign:", err)
					continue
				}
				fmt.Println("Campaign saved.")
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

// gameLoop handles the main game loop for the GM.
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

// playerMoveSubscriber listens for player move messages and updates the map accordingly.
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