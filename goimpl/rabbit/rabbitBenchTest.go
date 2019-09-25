package main

import (
	"log"
	"strconv"
	"sync"
	"testing"
	"time"
	"github.com/streadway/amqp"
)
func failOnError(err error, msg string) {
	if err != nil {
	  log.Fatalf("%s: %s", msg, err)
	}
  }

func TestRabbitPubSubThroughputWithBufferedChannel(t *testing.T) {

	nSec := 10
	cc := 5
	cn := make(chan int, cc)
	since := time.Now()
	due := time.Second * time.Duration(nSec)

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch_sub, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	done := make(chan bool)
	log.Println("start benchmark for Rabbit")
	go func() {
		defer func() {
			ch_sub.close()
		}()
		cnt := 0
		log.Println("ready for subscribe")
		for {
			sub.ReceiveMessage()
			cnt++
			nowSince := time.Since(since)
			if nowSince > due {
				break
			}
		}
		log.Println("subscribe throughput", cnt/nSec, "msg/sec")
		done <- true
	}()

	log.Println("ready for publish")
	for i := 0; i < cc; i++ {
		go func() {
			cnt := 0
			defer func() {
				cn <- cnt
			}()
			for {
				client.Publish("channel", strconv.FormatInt(time.Now().UnixNano(), 10))
				cnt++
				nowSince := time.Since(since)
				if nowSince > due {
					break
				}
			}
		}()
	}
	sum := 0
	for i := 0; i < cc; i++ {
		sum += <-cn
	}
	close(cn)
	log.Println("publish throughput", sum/nSec, "msg/sec")
	<-done
}

func TestRabbitPubSubLatemncyWithBufferedChannel(t *testing.T) {
	client := Rabbit.NewClient(&Rabbit.Options{Addr: "localhost:6379", Password: "", DB: 0})
	defer client.Close()
	nSec := 10
	cc := 2
	cn := make(chan int, cc)
	done := make(chan bool)
	sub := client.Subscribe("channel")

	since := time.Now()
	due := time.Second * time.Duration(nSec)
	log.Println("start benchmark for Rabbit")
	go func(sub *Rabbit.PubSub) {
		defer func() {
			sub.Close()
		}()
		cnt := 0
		latency := 0
		log.Println("ready for subscribe")
		for {
			msg, _ := sub.ReceiveMessage()
			cnt++
			nowSince := time.Since(since)
			ns, _ := strconv.Atoi(msg.Payload)
			latency += int(time.Now().UnixNano()) - ns
			if nowSince > due {
				break
			}
		}
		log.Println("subscribe throughput", latency/cnt, "ns", ":cnt", cnt)
		done <- true
	}(sub)
	log.Println("ready for publish")
	for i := 0; i < cc; i++ {
		go func() {
			cnt := 0
			defer func() {
				cn <- cnt
			}()
			for {
				client.Publish("channel", strconv.FormatInt(time.Now().UnixNano(), 10))
				cnt++
				nowSince := time.Since(since)
				if nowSince > due {
					break
				}
			}
		}()
	}
	sum := 0
	for i := 0; i < cc; i++ {
		sum += <-cn
	}
	close(cn)
	log.Println("publish throughput", sum/nSec, "msg/sec")
	<-done
}
