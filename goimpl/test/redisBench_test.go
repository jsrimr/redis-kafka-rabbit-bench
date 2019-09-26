package main

import (
	"log"
	"strconv"
	"sync"
	"testing"
	"time"
	
	"github.com/go-redis/redis"
)

type counter struct {
	i  int
	mu sync.Mutex
}

func (c *counter) increment() {
	c.mu.Lock()
	c.i++
	c.mu.Unlock()
}

func TestBuff(t *testing.T) {
	cn := make(chan int, 4)
	cn <- 1
	cn <- 2
	cn <- 3
	cn <- 4

	sum := 0
	for i := 0; i < 4; i++ {
		sum += <-cn
	}
	log.Println(sum)
}

func TestRedisBenchWithWaitGroup(t *testing.T) {
	client := redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: 0})
	defer client.Close()
	cnt := counter{i: 0}

	nSec := 3
	cc := 100

	since := time.Now()
	log.Println("start")
	wg := sync.WaitGroup{}
	wg.Add(cc)

	sub := client.Subscribe("channel")
	go func() {
		for {
			sub.ReceiveMessage()
			due := time.Since(since)
			if due > time.Second*time.Duration(nSec) {
				break
			}
		}
	}()

	for i := 0; i < cc; i++ {
		go func() {
			defer wg.Done()
			for {
				client.Publish("channel", "123")
				cnt.increment()
				due := time.Since(since)
				if due > time.Second*time.Duration(nSec) {
					break
				}
			}
		}()
	}
	wg.Wait()
	log.Println("set throughput", cnt.i/nSec)
	sub.Close()
}

func TestRedisPubSubThroughputWithBufferedChannel(t *testing.T) {
	client := redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: 0})
	defer client.Close()
	nSec := 1
	cc := 100
	cn := make(chan int, cc)
	done := make(chan bool)
	sub := client.Subscribe("channel")

	since := time.Now()
	due := time.Second * time.Duration(nSec)
	log.Println("start benchmark for redis")
	go func(sub *redis.PubSub) {
		defer func() {
			sub.Close()
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

func TestRedisPubSubThroughputWithMutliConnection(t *testing.T) {
	client := redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: 0})
	defer client.Close()
	nSec := 1
	cc := 100
	cn := make(chan int, cc)
	done := make(chan bool)
	sub := client.Subscribe("channel")

	since := time.Now()
	due := time.Second * time.Duration(nSec)
	log.Println("start benchmark for redis (multiple client)")
	go func(sub *redis.PubSub) {
		defer func() {
			sub.Close()
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
	}(sub)
	log.Println("ready for publish")
	for i := 0; i < cc; i++ {
		client := redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: 0})
		defer client.Close()
		go func(client *redis.Client) {
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
		}(client)
	}
	sum := 0
	for i := 0; i < cc; i++ {
		sum += <-cn
	}
	close(cn)
	log.Println("publish throughput", sum/nSec, "msg/sec")
	<-done
}

func TestRedisPubSubThroughputPipelined(t *testing.T) {
	client := redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: 0})
	defer client.Close()
	nSec := 1
	cc := 100
	pc := 10
	cn := make(chan int, cc)
	done := make(chan bool)
	sub := client.Subscribe("channel")

	since := time.Now()
	due := time.Second * time.Duration(nSec)
	log.Println("start benchmark for redis (multiple client + pipelined)")
	go func(sub *redis.PubSub) {
		defer func() {
			sub.Close()
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
	}(sub)
	log.Println("ready for publish")
	for i := 0; i < cc; i++ {
		client := redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: 0})
		defer client.Close()
		go func(client *redis.Client) {
			cnt := 0
			defer func() {
				cn <- cnt
			}()
			for {
				pipe := client.Pipeline()
				for j:=0; j<pc; j++ {
					pipe.Publish("channel", strconv.FormatInt(time.Now().UnixNano(), 10))
				}
				_, err := pipe.Exec()
				if err != nil {
					t.Error(err)
				}
				cnt+=pc
				nowSince := time.Since(since)
				if nowSince > due {
					break
				}
			}
		}(client)
	}
	sum := 0
	for i := 0; i < cc; i++ {
		sum += <-cn
	}
	close(cn)
	log.Println("publish throughput", sum/nSec, "msg/sec")
	<-done
}

func TestRedisPubSubThroughputMultiSubscribe(t *testing.T) {
	nSec := 1
	cc := 100
	sc := 2
	pc := 10
	cn := make(chan int, cc)
	done := make(chan bool)

	log.Println("start benchmark for redis (multiple client + pipelined + multisubscribe)")
	go func() {
		log.Println("ready for subscribe")
		sn := make(chan int, sc)
		since := time.Now()
		due := time.Second * time.Duration(nSec)
		for i:=0; i<sc; i++ {
			go func() {
				client := redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: 0})
				sub := client.Subscribe("channel")
				defer client.Close()
				defer sub.Close()
				cnt := 0
				for {
					sub.ReceiveMessage()
					cnt++
					nowSince := time.Since(since)
					if nowSince > due {
						break
					}
				}
				sn <- cnt
			}()
		}
		sum := 0
		for i:=0; i<sc; i++ {
			sum += <- sn
		}
		log.Println("subscribe throughput", sum / nSec, "msg/sec")
		close(sn)
		done <- true
	}()
	log.Println("ready for publish")
	since := time.Now()
	due := time.Second * time.Duration(nSec)
	for i := 0; i < cc; i++ {
		client := redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: 0})
		defer client.Close()
		go func(client *redis.Client) {
			cnt := 0
			defer func() {
				cn <- cnt
			}()
			for {
				pipe := client.Pipeline()
				for j:=0; j<pc; j++ {
					pipe.Publish("channel", strconv.FormatInt(time.Now().UnixNano(), 10))
				}
				_, err := pipe.Exec()
				if err != nil {
					t.Error(err)
				}
				cnt+=pc
				nowSince := time.Since(since)
				if nowSince > due {
					break
				}
			}
		}(client)
	}
	sum := 0
	for i := 0; i < cc; i++ {
		sum += <-cn
	}
	close(cn)
	log.Println("publish throughput", sum/nSec, "msg/sec")
	<-done
}

func TestRedisPubSubLatencyWithBufferedChannel(t *testing.T) {
	client := redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: 0})
	defer client.Close()
	nSec := 10
	cc := 2
	cn := make(chan int, cc)
	done := make(chan bool)
	sub := client.Subscribe("channel")

	since := time.Now()
	due := time.Second * time.Duration(nSec)
	log.Println("start benchmark for redis")
	go func(sub *redis.PubSub) {
		defer sub.Close()
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

func TestRedisSet(t *testing.T) {
	client := redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: 0})
	defer client.Close()
	nSec := 1
	cc := 100
	cn := make(chan int, cc)

	since := time.Now()
	due := time.Second * time.Duration(nSec)
	log.Println("start benchmark for redis")
	log.Println("ready for set")
	for i := 0; i < cc; i++ {
		go func() {
			cnt := 0
			defer func() {
				cn <- cnt
			}()
			for {
				client.Set("k", "123", 0)
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
}
