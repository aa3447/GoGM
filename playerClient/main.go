package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/campaign"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/io"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/mapLogic"
	player "home/aa3447/workspace/github.com/aa3447/GoGM/internal/playerLogic"
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

	var currentCampaign *campaign.Campaign

	newPlayer := player.NewPlayer("Hero", "The brave adventurer", "Warrior", "roll")
	//player := player.NewPlayer(io.GetInput()[0], "The brave adventurer", "Warrior")

	commands := io.GetInput()
	for {
		command := commands[0]
		switch command {
			case "manage":
				currentCampaign, newPlayer, err = managementLoop(channel)
				if err != nil {
					fmt.Println("Error in management loop:", err)
					return
				}
			case "play":
				if currentCampaign == nil{
					fmt.Println("Please join a campaign first.")
					continue
				}
				if newPlayer == nil{
					fmt.Println("Please create a player first.")
					continue
				}
				gameLoop(channel, newPlayer, currentCampaign)
			case "quit":
				fmt.Println("Quitting game.")
				return
			default:
				fmt.Println("Unknown command:", command)
		}
		commands = io.GetInput()
	}
}

func managementLoop(channel *ampq.Channel) (*campaign.Campaign, *player.Player, error){
	fmt.Println("Entering management loop. Type 'quit' to exit.")
	var currentPlayer *player.Player

	currentCampaign := campaign.NewCampaign("Player Campaign", "A player campaign for GoGM", player.GM{})

	commands := io.GetInput()
	for {
		command := commands[0]
		switch command {
			case "create":
				currentPlayer = createPlayer()
			case "update":
				//handleUpdate(args)
			case "load":
				//handleLoad(args)
			case "delete":
				//handleDelete(args)
			case "list":
				//handleList(args)
			case "join":
				//handleJoin(args)
			case "quit":
				if currentPlayer == nil{
					fmt.Println("Please create a player first.")
					continue
				}
				if currentCampaign == nil{
					fmt.Println("Please join a campaign first.")
					continue
				}
				
				fmt.Println("Exiting management loop.")
				err := pubsub.QueueDeclareAndBindSetup(channel, currentPlayer)	
				if err != nil{
					fmt.Println("Error declaring and binding queues:", err)
					return nil, nil, err
				}	
				return currentCampaign, currentPlayer, nil
			default:
				fmt.Println("Unknown command:", command)
		}
		commands = io.GetInput()
	}
}

func createPlayer() *player.Player{
	fmt.Println("Creating new player.")
	fmt.Print("Enter player name: ")
	name := io.GetInput()[0]
	fmt.Print("Enter player description: ")
	description := io.GetInput()[0]
	fmt.Print("Enter player background: ")
	background := io.GetInput()[0]
	fmt.Print("Enter stat generation method (roll, buy, assign): ")
	statMethod := io.GetInput()[0]
	var stats []int
	if statMethod == "assign"{
		fmt.Print("Enter stats as space-separated integers (Str Dex Int Con Cha Wis): ")
		statInputs := io.GetInput()
		for _, s := range statInputs{
			var stat int
			fmt.Sscanf(s, "%d", &stat)
			stats = append(stats, stat)
		}
	}
	player := player.NewPlayer(name, description, background, statMethod, stats)
	fmt.Println("Player created:", player.Name)
	return player
}

func gameLoop(channel *ampq.Channel, pl *player.Player, campaign *campaign.Campaign){
	fmt.Println("Entering game loop. Type 'quit' to exit.")
	fmt.Println("Welcome,", pl.Name)
	fmt.Println("You find yourself at the entrance of a mysterious location.")
	fmt.Println("Type 'move <direction>' to move (north, south, east, west), 'map' to view the map, or 'quit' to exit.")

	currentMap := playSubscribers(channel, pl, campaign)

	commands := io.GetInput()
	for {
		command := commands[0]
		args := commands[1:]
		switch command {
			case "move":
				handleMove(currentMap, pl, channel, args)
			case "action":
				//handleAction(args)
			case "map":
				currentMap.PrintMapWithPlayer(pl)
			case "quit":
				fmt.Println("Quitting game.")
				return
			default:
				fmt.Println("Unknown command:", command)
		}
		commands = io.GetInput()
	}
}

func playSubscribers(channel *ampq.Channel, pl *player.Player, campaign *campaign.Campaign) *mapLogic.Map{
	newMaps := mapQueueNewSubscriber(channel, pl)
	cMAp, err := campaign.NewGamestateWithExistingMap(<-newMaps)
	
	if err != nil{
		fmt.Println("Error creating game state:", err)
		return &mapLogic.Map{}
	}

	go addNewMaps(newMaps, campaign)
	go updateMaps(mapQueueUpdateSubscriber(channel, pl), campaign)
	go moveSubscriber(channel, pl, cMAp)

	pl.SetLocation(cMAp.EntranceLocation[0], cMAp.EntranceLocation[1])

	return cMAp
}

func handleMove(m *mapLogic.Map, player *player.Player, channel *ampq.Channel ,args []string){
	if len(args) < 1{
		fmt.Println("Move where?")
		return
	}
	var playerMove mapLogic.PlayerMove
	var err error

	direction := args[0]
	fmt.Println("Moving", direction)

	switch direction{
		case "north", "up":
			playerMove, err = m.MovePlayer(player,-1, 0)
		case "south", "down":
			playerMove, err = m.MovePlayer(player,1, 0)
		case "east", "right":
			playerMove, err = m.MovePlayer(player,0, 1)
		case "west", "left":
			playerMove, err = m.MovePlayer(player,0, -1)
		default:
			fmt.Println("Unknown direction:", direction)
	}

	if (err == nil) {
		err = pubsub.PublishToQueueAsJSON(channel, pubsub.MoveExchange, pubsub.PlayerMoveRoutingKey, &playerMove)
		if err != nil{
			log.Printf("Failed to publish move: %v\n", err)
			return
		}
		fmt.Println("Moved", direction)
		return
	}

	fmt.Printf("Failed to move: %v\n", err)
}

func addNewMaps(newMapsChan chan *mapLogic.Map, c *campaign.Campaign){
	for newMaps := range newMapsChan{
		newMap := newMaps
		c.AddMap(newMap)
		fmt.Println("New map added:", newMap.Name)
	}
}

func updateMaps(updateMapsChan chan *mapLogic.Map, c *campaign.Campaign){
	for updateMaps := range updateMapsChan{
		updatedMap := updateMaps
		for i, existingMap := range c.Maps{
			if existingMap.Name == updatedMap.Name{
				c.Maps[i] = updatedMap
				break
			}
		}
		fmt.Println("Map updated:", updatedMap.Name)
	}
}

func mapQueueNewSubscriber(channel *ampq.Channel, player *player.Player) chan *mapLogic.Map{
	newMaps := make(chan *mapLogic.Map)
	queueName := pubsub.MapQueueNew + "_" + player.Name

	subscribeArgs := []bool{true, false, false, false}
	msgs, err := pubsub.SubscribeToQueue(channel, queueName, subscribeArgs, subscribeArgs)
	if err != nil{
		fmt.Println("Error subscribing to map queue:", err)
		return nil
	}

	go func(){
		for d := range msgs{
			currentMap, err := serialization.JSONTo(d.Body, mapLogic.Map{})
			if err != nil{
				fmt.Println("Error deserializing map:", err)
				continue
			}
			_,err = os.Stat("./playerClient/map/" + currentMap.Name + ".json")

			if errors.Is(err, os.ErrNotExist) {
				err = serialization.SaveToFile(*currentMap, "player" ,"map" ,currentMap.Name)
				//err = serialization.SaveMapToFile(currentMap, "player", player.Name)
				if err != nil{
					fmt.Println("Error saving map to file:", err)
					continue
				}
				currentMap.FileLocation = "./playerClient/map/" + currentMap.Name
				newMaps <- currentMap
			}
		}
	}()
	return newMaps
}

func mapQueueUpdateSubscriber(channel *ampq.Channel, player *player.Player) chan *mapLogic.Map{
	updateMaps := make(chan *mapLogic.Map)
	queueName := pubsub.MapQueueUpdate + "_" + player.Name

	subscribeArgs := []bool{true, false, false, false}
	msgs, err := pubsub.SubscribeToQueue(channel, queueName, subscribeArgs, subscribeArgs)
	if err != nil{
		fmt.Println("Error subscribing to map queue:", err)
		return nil
	}

	go func(){
		for d := range msgs{
			currentMap, err := serialization.JSONTo(d.Body, mapLogic.Map{})
			if err != nil{
				fmt.Println("Error deserializing map:", err)
				continue
			}
			_,err = os.Stat("./playerClient/map/" + currentMap.Name + ".json")
			if errors.Is(err, os.ErrNotExist) {
				fmt.Println("Received update for unknown map:", currentMap.Name)
				continue
			}
			err = serialization.SaveToFile(*currentMap, "player", "map", currentMap.Name)
			if err != nil{
				fmt.Println("Error saving map to file:", err)
				continue
			}
			currentMap.FileLocation = "./playerClient/map/" + currentMap.Name
			updateMaps <- currentMap
		}
	}()
	return updateMaps
}

func moveSubscriber(channel *ampq.Channel, player *player.Player, m *mapLogic.Map) {
	subscribeArgs := []bool{true, false, false, false}
	msgs, err := pubsub.SubscribeToQueue(channel, pubsub.GMMoveQueue, subscribeArgs, subscribeArgs)
	if err != nil{
		log.Println("Error subscribing to GM move queue:", err)
		return
	}

	for d := range msgs{
		playerMove ,err := serialization.JSONTo(d.Body, mapLogic.PlayerMove{})
		if err != nil{
			log.Println("Error deserializing player move:", err)
			continue
		}
		for _, p := range m.GameState.Players{
			if p.Name != player.Name && p.Name == playerMove.PlayerName{
				m.MovePlayer(p, playerMove.To[0]-playerMove.From[0], playerMove.To[1]-playerMove.From[1])
			}
		}
	}
}
