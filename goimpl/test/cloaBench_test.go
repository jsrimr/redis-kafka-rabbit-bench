package main

import (
	// "bufio"
	"errors"
	"fmt"
	"log"
	"net"
	// "strconv"
	"testing"
	"time"
)

const delim = byte('\n')

func TestCloaPubSubThroughputWithBufferedChannel(t *testing.T) {
	nSec := 10
	cc := 5
	cn := make(chan int, cc)
	done := make(chan bool)
	since := time.Now()
	due := time.Second * time.Duration(nSec)
	log.Println("start benchmark for cloa")
	/*
		go func() {
			conn, err := net.Dial("tcp", "localhost:4242")
			if err != nil {
				t.Error(err)
			}
			defer conn.Close()
			cnt := 0
			log.Println("ready for subscribe")
			_, err = conn.Write([]byte("SUB 23482093 PP %d\r\n"))
			if err != nil {
				t.Error(err)
			}
			buf := make([]byte, 1024)
			sl := bufio.NewReader(conn)
			for {
				buf, err = sl.ReadBytes(delim)

				// l, err := conn.
				cnt ++
				nowSince := time.Since(since)
				if nowSince> due {
					break
				}
			}
			log.Println("subscribe throughput", cnt / nSec, "msg/sec")
			done <- true
		}()
	*/
	log.Println("ready for publish")
	conn, err := net.Dial("tcp", "localhost:4242")
	if err != nil {
		t.Error(err)
	}
	buf := make([]byte, 1024)
	for i := 0; i < cc; i++ {
		go func() {
			cnt := 0
			defer func() {
				cn <- cnt
			}()
			for {
				_, err = conn.Write([]byte(fmt.Sprintf("PUB 23482093 %d", time.Now().UnixNano())))
				if err != nil {
					t.Error(err)
				}
				_, err = conn.Read(buf)
				// if buf[:2] == "OK" {
				// 	t.Error(errors.New("pub failed"))
				// }
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

/*
func TestRedisPubSubLatemncyWithBufferedChannel(t *testing.T) {
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
		defer func() {
			sub.Close()
		}()
		cnt := 0
		latency := 0
		log.Println("ready for subscribe")
		for {
			msg, _ := sub.ReceiveMessage()
			cnt ++
			nowSince := time.Since(since)
			ns, _ := strconv.Atoi(msg.Payload)
			latency += int(time.Now().UnixNano())-ns
			if nowSince> due {
				break
			}
		}
		log.Println("subscribe throughput", latency / cnt, "ns", ":cnt", cnt)
		done <- true
	}(sub)
	log.Println("ready for publish")
	for i:=0; i<cc; i++ {
		go func() {
			cnt:=0
			defer func() {
				cn<-cnt
			}()
			for {
				client.Publish("channel", strconv.FormatInt(time.Now().UnixNano(),10))
				cnt++
				nowSince := time.Since(since)
				if nowSince> due {
					break
				}
			}
		}()
	}
	sum:=0
	for i:=0; i<cc; i++ {
		sum += <-cn
	}
	close(cn)
	log.Println("publish throughput", sum / nSec, "msg/sec")
	<- done
}

func TestRedisSet(t *testing.T) {
	client := redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: 0})
	defer client.Close()
	nSec := 10
	cc := 5
	cn := make(chan int, cc)

	since := time.Now()
	due := time.Second * time.Duration(nSec)
	log.Println("start benchmark for redis")
	log.Println("ready for set")
	for i:=0; i<cc; i++ {
		go func() {
			cnt:=0
			defer func() {
				cn<-cnt
			}()
			for {
				client.Set("k", "123456789012345678901234567989012", 0)
				cnt++
				nowSince := time.Since(since)
				if nowSince> due {
					break
				}
			}
		}()
	}
	sum:=0
	for i:=0; i<cc; i++ {
		sum += <-cn
	}
	close(cn)
	log.Println("publish throughput", sum / nSec, "msg/sec")
}
*/
