package main

import (
	"fmt"
	"log"
	"net"
	"bytes"
	"testing"
	"time"
	"errors"
)

const delim = byte('\n')

func TestCloaPubSubThroughputWithBufferedChannel(t *testing.T) {
	nSec := 1
	cc := 100
	cn := make(chan int, cc)
	done := make(chan bool)
	since := time.Now()
	due := time.Second * time.Duration(nSec)
	log.Println("start benchmark for cloa")
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
		for {
			_, err = conn.Read(buf)
			if err != nil {
				t.Error(err)
			}
			cnt ++
			nowSince := time.Since(since)
			if nowSince> due {
				break
			}
		}
		log.Println("subscribe throughput", cnt / nSec, "msg/sec")
		conn.Close()
		done <- true
	}()
	log.Println("ready for publish")
	buf := make([]byte, 1024)
	resp := []byte("OK")
	for i := 0; i < cc; i++ {
		conn, err := net.Dial("tcp", "localhost:4242")
		if err != nil {
			t.Error(err)
		}
		go func(i int, conn net.Conn) {
			cnt := 0
			defer func() {
				cn <- cnt
			}()
			for {
				_, err = conn.Write([]byte(fmt.Sprintf("PUB 23482093 PP %d\r\n", time.Now().UnixNano())))
				if err != nil {
					t.Error(err)
				}
				_, err = conn.Read(buf)
				if !bytes.Equal(buf[:2], resp) {
					t.Error(errors.New("pub failed"))
				}
				cnt++
				nowSince := time.Since(since)
				if nowSince > due {
					break
				}
			}
			conn.Close()
		}(i, conn)
	}
	sum := 0
	for i := 0; i < cc; i++ {
		sum += <-cn
	}
	close(cn)
	log.Println("publish throughput", sum/nSec, "msg/sec")
	<-done
}

func TestCloaPubSubThroughputWithMultiSubs(t *testing.T) {
	nSec := 1
	cc := 100
	cn := make(chan int, cc)
	done := make(chan bool)
	since := time.Now()
	due := time.Second * time.Duration(nSec)
	log.Println("start benchmark for cloa")
	go func() {
		log.Println("ready for subscribe")
		sn := make(chan int, cc)
		for i:=0; i<cc; i++ {
			go func(i int) {
				conn, err := net.Dial("tcp", "localhost:4242")
				if err != nil {
					t.Error(err)
				}
				defer conn.Close()
				cnt := 0
				_, err = conn.Write([]byte("SUB 23482093 PP %d\r\n"))
				if err != nil {
					t.Error(err)
				}
				buf := make([]byte, 1024)
				for {
					_, err = conn.Read(buf)
					if err != nil {
						t.Error(err)
					}
					cnt ++
					nowSince := time.Since(since)
					if nowSince> due {
						break
					}
				}
				sn <- cnt
			}(i)
		}
		sum := 0
		for i:=0; i<cc; i++ {
			sum += <- sn
		}
		log.Println("subscribe throughput", sum / nSec, "msg/sec")
		close(sn)
		done <- true
	}()

	log.Println("ready for publish")
	buf := make([]byte, 1024)
	resp := []byte("OK")
	for i := 0; i < cc; i++ {
		conn, err := net.Dial("tcp", "localhost:4242")
		if err != nil {
			t.Error(err)
		}
		go func(i int, conn net.Conn) {
			cnt := 0
			defer func() {
				cn <- cnt
			}()
			for {
				_, err = conn.Write([]byte(fmt.Sprintf("PUB 23482093 PP %d\r\n", time.Now().UnixNano())))
				if err != nil {
					t.Error(err)
				}
				_, err = conn.Read(buf)
				if !bytes.Equal(buf[:2], resp) {
					t.Error(errors.New("pub failed"))
				}
				cnt++
				nowSince := time.Since(since)
				if nowSince > due {
					break
				}
			}
			conn.Close()
		}(i, conn)
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
