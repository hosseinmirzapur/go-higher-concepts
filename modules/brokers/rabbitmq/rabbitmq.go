package rabbitmq

import (
	"context"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	conn *amqp.Connection
)

type chanQ struct {
	ch *amqp.Channel
	q  amqp.Queue
}

func handleErr(err error, msg string) {
	if err != nil {
		log.Panicf("%+v: %s", err, msg)
	}
}

func Connect() {
	conn, err := amqp.Dial(RabbitMQConfig().URI)
	handleErr(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
}

func CreateChanQueue(
	name string,
	durable, deleteWhenUnused, exclusive, noWait bool,
	args amqp.Table,
) chanQ {
	// create a channel for data to be broadcasted on
	ch, err := conn.Channel()
	handleErr(err, "Failed to open a channel")
	defer ch.Close()

	// create a queue for FIFO message processing
	q, err := ch.QueueDeclare(name, durable, deleteWhenUnused, exclusive, noWait, args)
	handleErr(err, "Failed to declare a queue")

	return chanQ{ch, q}
}

func Produce(
	chq chanQ,
	exchange string,
	mandatory, immediate bool,
	body []byte,
	contentType string,
) {
	// define a context to control produce timeout
	ctx, cancel := context.WithTimeout(context.Background(), RabbitMQConfig().Timeout)
	defer cancel()

	// publish the message through the channel and send it to the queue
	err := chq.ch.PublishWithContext(ctx, exchange, chq.q.Name, mandatory, immediate, amqp.Publishing{
		ContentType: contentType,
		Body:        body,
	})
	handleErr(err, "Failed to publish a message")
}

func Consume(
	chq chanQ,
	consumer string,
	autoAck, exclusive, noLocal, noWait bool,
	args amqp.Table,
) {
	// get all messages from the queue of the channel
	msgs, err := chq.ch.Consume(chq.q.Name, consumer, autoAck, exclusive, noLocal, noWait, args)
	handleErr(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s\n", d.Body)
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	// get all the data coming out of the goroutine channel
	<-forever
}
