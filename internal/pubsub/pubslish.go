package pubsub

import (
	"fmt"

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
		"topic",
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
	playerMoveQueueName := PlayerMoveQueue
	GMMoveQueueName := GMMoveQueue + "_" + player.Name
	
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
	playerMoveQueue, err := channel.QueueDeclare(
		playerMoveQueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil{
		return fmt.Errorf("failed to declare a queue: %v", err)
	}
	gmMoveQueue, err := channel.QueueDeclare(
		GMMoveQueueName,
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
		playerMoveQueue.Name,
		PlayerMoveRoutingKey,
		MoveExchange,
		false,
		nil,
	)
	if err != nil{
		return fmt.Errorf("failed to bind a queue: %v", err)
	}
	err = channel.QueueBind(
		gmMoveQueue.Name,
		GMMoveRoutingKey,
		MoveExchange,
		false,
		nil,
	)
	if err != nil{
		return fmt.Errorf("failed to bind a queue: %v", err)
	}

	return nil
}

// PublishToQueueAsJSON publishes the given data of type J to the specified exchange and routing key as a JSON message.
func PublishToQueueAsJSON[J serialization.JSONSafe](channel *ampq.Channel, exchange, routingKey string, data *J) error{
	JSONData, err := serialization.ToJSON(*data)
	if err != nil{
		return fmt.Errorf("failed to serialize to JSON: %v", err)
	}

	err = channel.Publish(
		exchange,
		routingKey,
		false,
		false,
		ampq.Publishing{
			ContentType: "application/json",
			Body:	JSONData,
		},
	)
	if err != nil{
		return fmt.Errorf("failed to publish a message: %v", err)
	}
	return nil
}


