package pubsub

import (
	"fmt"

	ampq "github.com/rabbitmq/amqp091-go"
)

func init(){
	conn, err := ampq.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil{
		fmt.Println("Failed to connect to RabbitMQ:", err)
		return
	}
	defer conn.Close()

	channel , err := conn.Channel()
	if err != nil{
		fmt.Println("Failed to open a channel:", err)
		return
	}
	defer channel.Close()

	err = channel.ExchangeDeclare(
		"mapExchange",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil{
		fmt.Println("Failed to declare an exchange:", err)
		return
	}

	mapQueue, err := channel.QueueDeclare(
		"mapQueue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil{
		fmt.Println("Failed to declare a queue:", err)
		return
	}

	err = channel.QueueBind(
		mapQueue.Name,
		"mapRoutingKey",
		"mapExchange",
		false,
		nil,
	)
	if err != nil{
		fmt.Println("Failed to bind a queue:", err)
		return
	}

}

