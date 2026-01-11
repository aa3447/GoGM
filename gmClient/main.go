package main

import (
	"fmt"
	"maps"
	"slices"
	"flag"

	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/campaign"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/io"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/mapLogic"
	GM "home/aa3447/workspace/github.com/aa3447/GoGM/internal/playerLogic"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/pubsub"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/serialization"

	ampq "github.com/rabbitmq/amqp091-go"
)

func main(){

	onlinePtr := flag.Bool("online", false, "Run the GM client in online mode")
	flag.Parse()
	var channel *ampq.Channel
	var conn *ampq.Connection
	var err error

	if *onlinePtr{
		conn, err = ampq.Dial("amqp://guest:guest@localhost:5672/")
		if err != nil{
			fmt.Println("failed to connect to RabbitMQ:", err)
			return
		}
		defer conn.Close()

		channel , err = conn.Channel()
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
	}

	var gameCampaign *campaign.Campaign

	gm := GM.NewGM("GameMaster", "The overseer of the game world")
	
	fmt.Println("Welcome,", gm.Name)
	fmt.Println("You are the Game Master.")
	fmt.Println("Type 'map' to view the map, or 'quit' to exit.")

	commands := io.GetInput()
	for {
		command := commands[0]
		args := commands[1:]
		switch command {
			case "load":
				gameCampaign, err = serialization.LoadFromJSONFile("gm", "campaigns", args[0], campaign.Campaign{})
				if err != nil{
					fmt.Println("Error loading campaign:", err)
					continue
				}
				fmt.Println("Campaign loaded:", gameCampaign.Name)
			case "campaign":
				gameCampaign, err = campaignManagementLoop(channel)
				if err != nil{
					fmt.Println("Error in campaign management loop:", err)
					return
				}
			case "play":
				if onlinePtr == nil || !*onlinePtr || channel == nil{
					fmt.Println("Please run the GM client with -online flag to enter play mode.")
					commands = io.GetInput()
					continue
				}
				players := []*GM.Player{}
				// change to connected players for gamestate :todo:
				for _, player := range gameCampaign.Players{
					players = append(players, player)
				}
				cMap := gameCampaign.CurrentMap

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
	var loadedMap mapLogic.Map

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
			case "edit_map":
				// Implement map editor here
			case "load_map": // loads a map from file
				if len(args) == 0 {
					fmt.Println("Please specify the map name to load.")
					continue
				}
				loadedMapPtr, err := serialization.LoadFromJSONFile("gm", "map", args[0], mapLogic.Map{})
				if err != nil{
					fmt.Println("Error loading map:", err)
					continue
				}
				loadedMap = *loadedMapPtr
				fmt.Println("Map loaded:", args[0])
			case "save_map": // saves a map to file
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
			case "add_map": // adds the loaded map to the campaign
				if loadedMap.Name == "" {
					fmt.Println("No map loaded to add. Please load a map first.")
					continue
				}
				currentCampaign.AddMap(&loadedMap)
				fmt.Println("Map added to campaign:", loadedMap.Name)
			case "remove_map": // removes a map from the campaign
				if len(args) == 0 {
					fmt.Println("Please specify the map name to remove.")
					continue
				}
				currentCampaign.RemoveMap(args[0])
				fmt.Println("Map removed from campaign:", args[0])
			case "current_map": // sets the map that is loaded when the campaign is started
				if len(args) == 0 {
					fmt.Println("Please specify the map name to set as the current map.")
					continue
				}
				err := currentCampaign.SetCurrentMap(args[0])
				if err != nil{
					fmt.Println("Error setting current map:", err)
					continue
				}
				fmt.Println("Current map set to:", args[0])
			case "load_campaign": // loads a campaign from file
				loadedCampaign, err := serialization.LoadFromJSONFile("gm", "campaigns", args[0], campaign.Campaign{})
				if err != nil{
					fmt.Println("Error loading campaign:", err)
					continue
				}
				currentCampaign = loadedCampaign
				fmt.Println("Campaign loaded:", currentCampaign.Name)
			case "save_campaign": // saves the current campaign to file
				err := serialization.SaveToFile(*currentCampaign, "gm", "campaigns", currentCampaign.Name)
				if err != nil{
					fmt.Println("Error saving campaign:", err)
					continue
				}
				fmt.Println("Campaign saved.")
			case "add_player":
				// Implement player addition logic here
			case "edit_player":
				// Implement player editing logic here
			case "remove_player":
				// Implement player removal logic here
			case "add_npc":
				// Implement NPC addition logic here
			case "edit_npc":
				// Implement NPC editing logic here
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
			case "map": // view the current map
				m.PrintMapDebugWithPlayers(currentPlayers)
			case "send": // send info to players
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
			case "update": // send updates to players
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