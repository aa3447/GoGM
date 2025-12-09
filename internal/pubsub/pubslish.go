package pubsub

import (
	"fmt"

	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/mapLogic"
	player "home/aa3447/workspace/github.com/aa3447/GoGM/internal/playerLogic"
	"home/aa3447/workspace/github.com/aa3447/GoGM/internal/serialization"

	ampq "github.com/rabbitmq/amqp091-go"
)

func Start() error{
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
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil{
		return fmt.Errorf("failed to declare an exchange: %v", err)
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
		return fmt.Errorf("failed to declare an exchange: %v", err)
	}

	return nil
}

func QueueDeclareAndBindInit(channel *ampq.Channel, player *player.Player) error{
	mapQueueName := MapQueue + "_" + player.Name
	moveQueueName := MoveQueue + "_" + player.Name
	mapQueue, err := channel.QueueDeclare(
		mapQueueName,
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
		mapQueue.Name,
		MapRoutingKey,
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


