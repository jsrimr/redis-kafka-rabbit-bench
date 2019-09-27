package main
// refer : https://ednsquare.com/story/golang-implementing-kafka-consumers-using-sarama------J1JS3J
// producer : https://medium.com/rahasak/kafka-producer-with-golang-fab7348a5f9a
import (
	"log"
	"strconv"
	"testing"
	"time"
	"github.com/Shopify/sarama"
)
func TestKafkaPubSubThroughput(t *testing.T) {

	nSec := 1
	cc := 100
	cn := make(chan int, cc)
	due := time.Second * time.Duration(nSec)
	log.Println("start benchmark for kafka")
	since := time.Now()

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	brokers := []string{"localhost:9092"}

	// Create new consumer
	master, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := master.Close(); err != nil {
			panic(err)
		}
	}()

	topic := "test"
	
	consumer, err := master.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}
	// Count how many message processed
	// Get signnal for finish
	done := make(chan bool)
	msgCount := 0
	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				log.Println(err)
			case _ = <-consumer.Messages():
				msgCount++
				// log.Println("Received messages", string(msg.Key), string(msg.Value))
			}
			nowSince := time.Since(since)
			if nowSince > due {
				break
			}
		}
		done <- true
	}()
	
	log.Println("ready for publish")
	//--publish
	config = sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokers, config) // async 도 있는듯. 차이가 뭘까
    if err != nil {
		// Should not reach here
		panic(err)
	}

	for i := 0; i < cc; i++ {
		go func() {
			cnt := 0
			defer func() {
				cn <- cnt
			}()
			for {
				msg := &sarama.ProducerMessage{
					Topic: topic,
					Value: sarama.StringEncoder(strconv.FormatInt(time.Now().UnixNano(), 10)),
				}
				producer.SendMessage(msg)
				// log.Println(topic, partition, offset)
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
	log.Println("sub throughput", msgCount/nSec, "msg/sec")
	<-done
}


func TestKafkaPubSubLatency(t *testing.T) {
	
	nSec := 1
	cc := 1
	cn := make(chan int, cc)
	due := time.Second * time.Duration(nSec)
	done := make(chan bool)
	topic := "test"
	
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	brokers := []string{"localhost:9092"}
	// Create new consumer
	master, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := master.Close(); err != nil {
			panic(err)
		}
	}()
	consumer, err := master.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}

	log.Println("start benchmark for Kafka")
	since := time.Now()
	latency := 0.0
	cnt:=0
	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				log.Println(err)
			case msg := <-consumer.Messages():
				cnt++
				// log.Println("receved msg: " , string(msg.Value))
				ns ,_ := strconv.ParseFloat(string(msg.Value),32)
				// log.Println("converted value: " , ns)
				latency += float64(time.Now().UnixNano()) - ns
				log.Println("latency", latency,)
			}
			nowSince := time.Since(since)
			if nowSince > due {
				break
			}
		}
		done <- true
	}()

	log.Println("ready for publish")
	config = sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokers, config) // async 도 있는듯. 차이가 뭘까 -> 별 차이없다고 한다
    if err != nil {
		// Should not reach here
		panic(err)
	}
	for i := 0; i < cc; i++ {
		go func() {
			cnt := 0
			defer func() {
				cn <- cnt
			}()
			for {
				msg := &sarama.ProducerMessage{
					Topic: topic,
					Value: sarama.StringEncoder(strconv.FormatInt(time.Now().UnixNano(), 10)),
				}
				// log.Println("Sent : ", msg.Value)
				producer.SendMessage(msg)
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
	log.Println("sub throughput", cnt/nSec, "msg/sec")
	log.Println("pub throughput", sum/nSec, "msg/sec")
	log.Println("latency avg", latency/float64(cnt), "ns")
	<-done
}