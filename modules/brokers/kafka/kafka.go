package kafka

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

var connURL = fmt.Sprintf("%s:%s", KafkaConfig().Host, KafkaConfig().Port)

// partition can usually be zero.
// produce is an "atomic" function so either all messages are written or none are.
func Produce(msgs []kafka.Message, topic string, partition int) error {

	// establish a connection to the kafka broker
	conn, err := kafka.DialLeader(
		context.Background(),
		"tcp",
		connURL,
		topic,
		partition,
	)
	if err != nil {
		return err
	}

	// set the write deadline for the connection
	conn.SetWriteDeadline(time.Now().Add(6 * time.Second))

	// write messages of type "byte" to the broker channel
	_, err = conn.WriteMessages(
		msgs...,
	)
	if err != nil {
		return err
	}

	// close the connection and if error found log the error
	if err := conn.Close(); err != nil {
		log.Fatal(err)
	}

	return nil
}

func Consume(topic string, partition int) error {
	// establish a connection to the kafka broker
	conn, err := kafka.DialLeader(
		context.Background(),
		"tcp",
		connURL,
		topic,
		partition,
	)
	if err != nil {
		return err
	}

	conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	batch := conn.ReadBatch(10e3, 1e6) // fetch 10KB min, 1MB max
	b := make([]byte, 10e3)            // 10KB max per message
	for {
		n, err := batch.Read(b)
		if err != nil {
			break
		}
		fmt.Println(string(b[:n])) // we can change the way we consume the data
	}

	if err := batch.Close(); err != nil {
		return err
	}

	if err := conn.Close(); err != nil {
		return err
	}

	return nil
}
