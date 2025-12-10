package pubsub

// RabbitMQ Exchanges
const (
	MapExchange    = "mapExchange"
	MoveExchange   = "moveExchange"
)

// RabbitMQ Routing Keys
const (
	MapNewRoutingKey  = "mapNewRoutingKey"
	MapUpdateRoutingKey = "mapUpdateRoutingKey"
	MoveRoutingKey = "moveRoutingKey"
)

// RabbitMQ Queues
const (
	MapQueueNew       = "mapQueueNew"
	MapQueueUpdate    = "mapQueueUpdate"
	MoveQueue      = "moveQueue"
)
