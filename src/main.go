package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/streadway/amqp"
)

type jsonDockerhub struct {
	Repository RepositoryType `json:"repository"`
	PushData   PushDataType   `json:"push_data"`
}

type RepositoryType struct {
	RepoName string `json:"repo_name"`
}

type PushDataType struct {
	Tag    string `json:"tag"`
	Pusher string `json:"pusher"`
	// Images []string `json:"images"`
}

var (
	uri          = flag.String("uri", "amqp://guest:guest@localhost:5672/", "AMQP URI")
	exchange     = flag.String("exchange", "test-exchange", "Durable, non-auto-deleted AMQP exchange name")
	exchangeType = flag.String("exchange-type", "direct", "Exchange type - direct|fanout|topic|x-custom")
	queue        = flag.String("queue", "test-queue", "Ephemeral AMQP queue name")
	bindingKey   = flag.String("key", "test-key", "AMQP binding key")
	consumerTag  = flag.String("consumer-tag", "simple-consumer", "AMQP consumer tag (should not be blank)")
	lifetime     = flag.Duration("lifetime", 5*time.Second, "lifetime of process before shutdown (0s=infinite)")
	jobsPath     = flag.String("jobs-path", "/usr/local/opt/rhbuilder", "Path to RabbitHook jobs")
	logPath      = flag.String("log-path", "/var/log/rhbuilder.log", "Path to log file")
)

var Info *log.Logger
var Error *log.Logger

func init() {
	flag.Parse()
}

func createLogger() *os.File {
	logFile, err := os.OpenFile(*logPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)

	var multi io.Writer

	if err != nil {
		log.Println("Uh oh! Could not open log file.")
		log.Printf("Ensure there is a file at %q\n", *logPath)
		log.Println("Skeeping log file output altogether.")
		multi = io.MultiWriter(os.Stdout)
	} else {
		multi = io.MultiWriter(logFile, os.Stdout)
	}

	Info = log.New(multi,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(multi,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	return logFile
}

func main() {

	logFile := createLogger()
	defer logFile.Close()

	c, err := NewConsumer(*uri, *exchange, *exchangeType, *queue, *bindingKey, *consumerTag)
	if err != nil {
		Error.Fatalf("%s", err)
	}

	if *lifetime > 0 {
		Info.Printf("running for %s", *lifetime)
		time.Sleep(*lifetime)
	} else {
		Info.Printf("running forever")
		select {}
	}

	Info.Printf("shutting down")

	if err := c.Shutdown(); err != nil {
		Error.Fatalf("error during shutdown: %s", err)
	}
}

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	tag     string
	done    chan error
}

func NewConsumer(amqpURI, exchange, exchangeType, queueName, key, ctag string) (*Consumer, error) {
	c := &Consumer{
		conn:    nil,
		channel: nil,
		tag:     ctag,
		done:    make(chan error),
	}

	var err error

	Info.Printf("dialing %q", amqpURI)
	c.conn, err = amqp.Dial(amqpURI)
	if err != nil {
		return nil, fmt.Errorf("Dial: %s", err)
	}

	go func() {
		fmt.Printf("closing: %s", <-c.conn.NotifyClose(make(chan *amqp.Error)))
	}()

	Info.Printf("got Connection, getting Channel")

	c.channel, err = c.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("Channel: %s", err)
	}

	Info.Printf("got Channel, declaring Exchange (%q)", exchange)

	if err = c.channel.ExchangeDeclare(
		exchange,     // name of the exchange
		exchangeType, // type
		false,        // durable CHANGED FROM TRUE
		true,         // delete when complete CHANGED FROM FALSE
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return nil, fmt.Errorf("Exchange Declare: %s", err)
	}

	Info.Printf("declared Exchange, declaring Queue %q", queueName)
	queue, err := c.channel.QueueDeclare(
		queueName, // name of the queue
		false,     // durable CHANGED from true
		true,      // delete when usused CHANGED FROM FALSE
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("Queue Declare: %s", err)
	}

	Info.Printf("declared Queue (%q %d messages, %d consumers), binding to Exchange (key %q)",
		queue.Name, queue.Messages, queue.Consumers, key)

	if err = c.channel.QueueBind(
		queue.Name, // name of the queue
		key,        // bindingKey
		exchange,   // sourceExchange
		false,      // noWait
		nil,        // arguments
	); err != nil {
		return nil, fmt.Errorf("Queue Bind: %s", err)
	}

	Info.Printf("Queue bound to Exchange, starting Consume (consumer tag %q)", c.tag)

	deliveries, err := c.channel.Consume(
		queue.Name, // name
		c.tag,      // consumerTag,
		false,      // noAck
		false,      // exclusive
		false,      // noLocal
		false,      // noWait
		nil,        // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("Queue Consume: %s", err)
	}

	go handle(deliveries, c.done)

	return c, nil
}

func (c *Consumer) Shutdown() error {
	// will close() the deliveries channel
	if err := c.channel.Cancel(c.tag, true); err != nil {
		return fmt.Errorf("Consumer cancel failed: %s", err)
	}

	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}

	defer Info.Printf("AMQP shutdown OK")

	// wait for handle() to exit
	return <-c.done
}

func handle(deliveries <-chan amqp.Delivery, done chan error) {
	for d := range deliveries {
		Info.Printf(
			"got %dB delivery: [%v], key: (%q), \n%q",
			len(d.Body),
			d.DeliveryTag,
			d.RoutingKey,
			d.Body,
		)

		processUpdate(d.Body)

		d.Ack(false)
	}
	Info.Printf("handle: deliveries channel closed")
	done <- nil
}

func processUpdate(msg []byte) {
	var payload jsonDockerhub

	err := json.Unmarshal(msg, &payload)

	if err != nil {
		Error.Fatalf("%s", err)
	}

	if len(payload.Repository.RepoName) == 0 {
		Info.Printf("=> %q", payload)
		return
	}
	data := payload.PushData
	repo := payload.Repository

	Info.Printf("Updated repo: %q, Tag: %q\n", repo.RepoName, data.Tag)

	cmdPath := path.Join(*jobsPath, repo.RepoName)

	Info.Print("Built CMD: %q\n", cmdPath)

	if _, err := os.Stat(cmdPath); err == nil {

		out, err := exec.Command(cmdPath, "--id", repo.RepoName, "--tag", data.Tag).Output()

		if err != nil {
			Info.Println("Error Executing Command %q", err)
		}

		Info.Printf("Job executed. Output: %q\n", out)
	} else {
		Info.Println("Skeeping Job. File not found.")
	}
}
