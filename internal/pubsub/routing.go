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

	PlayerMoveRoutingKey = "playerMoveRoutingKey"
	GMMoveRoutingKey     = "gmMoveRoutingKey"
)

// RabbitMQ Queues
const (
	MapQueueNew       = "mapQueueNew"
	MapQueueUpdate    = "mapQueueUpdate"

	PlayerMoveQueue = "playerMoveQueue"
	GMMoveQueue     = "gmMoveQueue"
)
