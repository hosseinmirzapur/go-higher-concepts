package memphis

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/memphisdev/memphis.go"
)

var conn *memphis.Conn

const FETCH_SIZE int = 10

func Connect() {
	conn, err := memphis.Connect(
		MemphisConfig().Hostname,
		MemphisConfig().ApplicationUser,
		MemphisConfig().Password,
	)
	if err != nil {
		log.Fatal("cannot connect to memphis broker")
	}

	defer conn.Close()

}

func Produce(msg []byte, stationName, producerName string) error {
	// create a producer with a producer-name and a station-name
	p, err := conn.CreateProducer(stationName, producerName)
	if err != nil {
		return err
	}

	// instantiate memphis headers
	hdrs := memphis.Headers{}
	hdrs.New()
	// can set headers via this:
	// hdrs.Add("key", "value")

	return p.Produce(msg, memphis.MsgHeaders(hdrs))
}

func Consume(stationName, consumerName string) {
	// create a single consumer to connect to memphis broker and wait for data created by producer
	consumer, err := conn.CreateConsumer(
		stationName,
		consumerName,
		memphis.PullInterval(10*time.Second),
	)
	if err != nil {
		log.Fatal(err)
	}

	consumer.SetContext(context.Background())

	// fetch messages with a predefined batch_size
	messages, err := consumer.Fetch(FETCH_SIZE, true)
	if err != nil {
		log.Fatal(err)
	}

	var msg_map map[string]interface{}
	for _, msg := range messages {
		err = json.Unmarshal(msg.Data(), &msg_map)
		if err != nil {
			log.Println(err)
			continue
		}

		// Do whatever desired with the message

		// sending an ACK packet to increase connection/network reliablity
		err = msg.Ack()
		if err != nil {
			log.Println(err)
			continue
		}

	}
	time.Sleep(time.Second * 30)

	// there will be 3 fetches due to the 30 seconds of sleep and the 10 seconds of pull interval
	// then the function returns(breaks) and the connection is closed
}
