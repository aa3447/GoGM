package pubsub

import (
	"fmt"
	ampq "github.com/rabbitmq/amqp091-go"
)

func SubscribeToMapQueue(channel *ampq.Channel, exchangeName, routingKey string) (<-chan ampq.Delivery, error){
	
	// Ensure the exchange exists
	mapQueue, err := channel.QueueDeclare(
		MapQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil{
		return nil,fmt.Errorf("failed to declare map queue: %v", err)
	}

	msgs, err := channel.Consume(
		mapQueue.Name,
		"map_consumer",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil{
		return nil, err
	}

	return msgs, nil
}