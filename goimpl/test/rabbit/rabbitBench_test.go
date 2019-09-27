package main

import (
	"log"
	"strconv"
	// "sync"
	"testing"
	"time"
	"github.com/streadway/amqp"
	// "encoding/binary"
)

func failOnError(err error, msg string) {
	if err != nil {
	  log.Fatalf("%s: %s", msg, err)
	}
  }

func TestRabbitPubSubThroughputWithBufferedChannel(t *testing.T) {

	nSec := 1
	cc := 1
	cn := make(chan int, cc)
	due := time.Second * time.Duration(nSec)
	
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	ch_sub, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	
	q, err := ch_sub.QueueDeclare(
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
	since := time.Now()
	go func(ch_sub *amqp.Channel) {
		defer func() {
			ch_sub.Close()
		}()
		cnt := 0
		log.Println("ready for subscribe")
		
		msgs, _ := ch_sub.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
		)
		
		for _ = range msgs {
			cnt++
			nowSince := time.Since(since)
			if nowSince > due {
				break
			}
		}
		log.Println("subscribe throughput", cnt/nSec, "msg/sec")
		done <- true
	}(ch_sub)

	log.Println("ready for publish")
	ch_pub, err := conn.Channel()
	q, err = ch_pub.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")
	for i := 0; i < cc; i++ {
		go func() {
			cnt := 0
			defer func() {
				cn <- cnt
			}()
			for {
				err = ch_pub.Publish(
					"",     // exchange
					q.Name, // routing key
					false,  // mandatory
					false,  // immediate
					amqp.Publishing{
						ContentType: "text/plain",
						Body:        []byte(strconv.FormatInt(time.Now().UnixNano(), 10)),
					})
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
	
	nSec := 1
	cc := 1
	cn := make(chan int, cc)
	due := time.Second * time.Duration(nSec)
	done := make(chan bool)
	
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	ch_sub, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	
	q, err := ch_sub.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")
	
	log.Println("start benchmark for Rabbit")
	since := time.Now()
	go func(ch_sub *amqp.Channel) {
		defer func() {
			ch_sub.Close()
		}()
		cnt := 0
		latency := 0
		log.Println("ready for subscribe")
		msgs, _ := ch_sub.Consume(
			q.Name, // queue
			"",     // consumer
			true,   // auto-ack
			false,  // exclusive
			false,  // no-local
			false,  // no-wait
			nil,    // args
			)
		for msg := range msgs{
			cnt++
			nowSince := time.Since(since)
			ns , _ := strconv.Atoi(string(msg.Body))
			now := int(time.Now().UnixNano())
			latency +=  now - ns
			log.Println(latency, cnt, now, ns, now-ns)

			if nowSince > due {
				break
			}
		}
		log.Println("subscribe latency", latency/cnt, "ns", ":cnt", cnt)
		done <- true
	}(ch_sub)

	log.Println("ready for publish")
	ch_pub, err := conn.Channel()
	q, err = ch_pub.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")
	for i := 0; i < cc; i++ {
		go func() {
			cnt := 0
			defer func() {
				cn <- cnt
			}()
			for {
				err = ch_pub.Publish(
					"",     // exchange
					q.Name, // routing key
					false,  // mandatory
					false,  // immediate
					amqp.Publishing{
						ContentType: "text/plain",
						Body:        []byte(strconv.FormatInt(time.Now().UnixNano(), 10)),
					})
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
