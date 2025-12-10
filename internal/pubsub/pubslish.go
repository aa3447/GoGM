package pubsub

import (
	"fmt"

	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/mapLogic"
	player "home/aa3447/workspace/github.com/aa3447/GoGM/internal/playerLogic"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/serialization"

	ampq "github.com/rabbitmq/amqp091-go"
)

func SetupExchanges() error{
	conn, err := ampq.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil{
		return fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	channel , err := conn.Channel()
	if err != nil{
		return fmt.Errorf("failed to open a channel: %v", err)
	}
	defer channel.Close()

	err = channel.ExchangeDeclare(
		MapExchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil{
		return fmt.Errorf("failed to declare map exchange: %v", err)
	}
	err = channel.ExchangeDeclare(
		MoveExchange,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil{
		return fmt.Errorf("failed to declare move exchange: %v", err)
	}

	return nil
}


func QueueDeclareAndBindSetup(channel *ampq.Channel, player *player.Player) error{
	mapQueueNewName := MapQueueNew + "_" + player.Name
	mapQueueUpdateName := MapQueueUpdate + "_" + player.Name
	moveQueueName := MoveQueue + "_" + player.Name
	
	mapQueueNew, err := channel.QueueDeclare(
		mapQueueNewName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil{
		return fmt.Errorf("failed to declare a queue: %v", err)
	}
	mapQueueUpdate, err := channel.QueueDeclare(
		mapQueueUpdateName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil{
		return fmt.Errorf("failed to declare a queue: %v", err)
	}
	moveQueue, err := channel.QueueDeclare(
		moveQueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil{
		return fmt.Errorf("failed to declare a queue: %v", err)
	}

	err = channel.QueueBind(
		mapQueueNew.Name,
		MapNewRoutingKey,
		MapExchange,
		false,
		nil,
	)
	if err != nil{
		return fmt.Errorf("failed to bind a queue: %v", err)
	}
	err = channel.QueueBind(
		mapQueueUpdate.Name,
		MapUpdateRoutingKey,
		MapExchange,
		false,
		nil,
	)
	if err != nil{
		return fmt.Errorf("failed to bind a queue: %v", err)
	}
	err = channel.QueueBind(
		moveQueue.Name,
		MoveRoutingKey,
		MoveExchange,
		false,
		nil,
	)
	if err != nil{
		return fmt.Errorf("failed to bind a queue: %v", err)
	}

	return nil
}

func PublishMapToQueue(channel *ampq.Channel, exchange, routingKey string, tileMap *mapLogic.Map) error{
	mapJSON, err := serialization.MapToJSON(tileMap)
	if err != nil{
		return fmt.Errorf("failed to serialize map to JSON: %v", err)
	}

	err = channel.Publish(
		exchange,
		routingKey,
		false,
		false,
		ampq.Publishing{
			ContentType: "application/json",
			Body:	mapJSON,
		},
	)
	if err != nil{
		return fmt.Errorf("failed to publish a message: %v", err)
	}
	return nil
}


