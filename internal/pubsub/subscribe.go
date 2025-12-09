package pubsub

import (
	"fmt"
	ampq "github.com/rabbitmq/amqp091-go"
)

func SubscribeToMapQueue(channel *ampq.Channel, queueName string) (<-chan ampq.Delivery, error){
	
	// Ensure the queue exists
	mapQueue, err := channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil{
		return nil, fmt.Errorf("failed to declare a queue: %v", err)
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