package rabbitmq

// Configuration holds RabbitMQ connection settings.
type Configuration struct {
	URL string `env:"RABBITMQ_URL" envDefault:"amqp://guest:guest@localhost:5672/"`
}
