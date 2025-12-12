package pubsub

import (
	"fmt"
	ampq "github.com/rabbitmq/amqp091-go"
)

// SubscribeToQueue subscribes to the specified queue with given arguments and returns a channel of deliveries.
func SubscribeToQueue(channel *ampq.Channel, queueName string, queueArgs []bool, consumerArgs []bool) (<-chan ampq.Delivery, error){
	
	if len(queueArgs) != 4 || len(consumerArgs) != 4{
		return nil, fmt.Errorf("invalid args length")
	}

	// Ensure the queue exists
	mapQueue, err := channel.QueueDeclare(
		queueName,
		queueArgs[0],
		queueArgs[1],
		queueArgs[2],
		queueArgs[3],
		nil,
	)
	if err != nil{
		return nil, fmt.Errorf("failed to declare a queue: %v", err)
	}

	msgs, err := channel.Consume(
		mapQueue.Name,
		"",
		consumerArgs[0],
		consumerArgs[1],
		consumerArgs[2],
		consumerArgs[3],
		nil,
	)
	if err != nil{
		return nil, fmt.Errorf("failed to link consumer: %v", err)
	}

	return msgs, nil
}