package main

import (
	"fmt"
	"log"
	"net"
	"bytes"
	"testing"
	"time"
	"errors"
	"strings"
	"strconv"
)

const delim = byte('\n')
var resp = []byte("OK")

func TestCloaPubSubThroughputWithBufferedChannel(t *testing.T) {
	nSec := 1
	cc := 100
	cn := make(chan int, cc)
	done := make(chan bool)
	since := time.Now()
	due := time.Second * time.Duration(nSec)
	log.Println("start benchmark for cloa")
	log.Println("start benchmark for cloa")
	log.Println("> publisher: ", cc)
	log.Println("> subscribers: ", 1)
	log.Println("> duration: ", nSec, "seconds")
	
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
		buf := make([]byte, 65536)
		_, err = conn.Read(buf)
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(buf[:2], resp) {
			t.Error(errors.New("sub failed"))
		}
		READSUBS:
		for {
			n, err := conn.Read(buf)
			if err != nil {
				t.Error(err)
			}
			msgs := strings.Split(string(buf[:n]), "\r\n")
			for _, msg := range msgs {
				if msg == "" {
					break
				}
				cnt ++
				nowSince := time.Since(since)
				if nowSince> due {
					break READSUBS
				}
			}
		}
		log.Println("subscribe throughput", cnt / nSec, "msg/sec")
		conn.Close()
		done <- true
	}()
	log.Println("ready for publish")
	buf := make([]byte, 1024)
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
	log.Println("> publisher: ", cc)
	log.Println("> subscribers: ", cc)
	log.Println("> duration: ", nSec, "seconds")
	
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
				buf := make([]byte, 65536)
				_, err = conn.Read(buf)
				if err != nil {
					t.Error(err)
				}
				if !bytes.Equal(buf[:2], resp) {
					t.Error(errors.New("sub failed"))
				}
				READSUBS:
				for {
					n, err := conn.Read(buf)
					if err != nil {
						t.Error(err)
					}
					msgs := strings.Split(string(buf[:n]), "\r\n")
					for _, msg := range msgs {
						if msg == "" {
							break
						}
						cnt ++
						nowSince := time.Since(since)
						if nowSince> due {
							break READSUBS
						}
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
	buf := make([]byte, 65536)
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

func TestCloaPubSubThroughputWithPipelined(t *testing.T) {
	nSec := 1
	cc := 100
	pc := 10
	cn := make(chan int, cc)
	done := make(chan bool)
	since := time.Now()
	due := time.Second * time.Duration(nSec)
	log.Println("start benchmark for cloa")
	log.Println("> publishers: ", cc)
	log.Println("> subscribers: ", cc)
	log.Println("> duration: ", nSec, "seconds")

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
				buf := make([]byte, 65536)
				_, err = conn.Read(buf)
				if err != nil {
					t.Error(err)
				}
				if !bytes.Equal(buf[:2], resp) {
					t.Error(errors.New("sub failed"))
				}
				READSUBS:
				for {
					n, err := conn.Read(buf)
					if err != nil {
						t.Error(err)
					}
					msgs := strings.Split(string(buf[:n]), "\r\n")
					for _, msg := range msgs {
						if msg == "" {
							break
						}
						cnt ++
						nowSince := time.Since(since)
						if nowSince> due {
							break READSUBS
						}
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
	buf := make([]byte, 65536)
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
				/* simulates pipelining by repeated commands */
				_, err = conn.Write([]byte(strings.Repeat(fmt.Sprintf("PUB 23482093 PP %d\r\n", time.Now().UnixNano()),pc)))
				if err != nil {
					t.Error(err)
				}
				n, err := conn.Read(buf)
				if err != nil {
					t.Error(err)
				}
				if !bytes.Equal(buf[:2], resp) {
					log.Println("pub RESPONSE", string(buf[:n]), "-------")
					t.Error(errors.New("pub failed"))
				}
				// log.Println(string(buf[:n]))
				cnt += pc
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

func TestCloaPubSubLatencyWithBufferedChannel(t *testing.T) {
	nSec := 1
	cc := 100
	cn := make(chan int, cc)
	done := make(chan bool)
	since := time.Now()
	due := time.Second * time.Duration(nSec)
	log.Println("start benchmark for cloa")
	log.Println("> publishers: ", cc)
	log.Println("> subscribers: ", cc)
	log.Println("> duration: ", nSec, "seconds")

	go func() {
		conn, err := net.Dial("tcp", "localhost:4242")
		defer conn.Close()
		if err != nil {
			t.Error(err)
		}
		defer conn.Close()
		cnt := 0
		latency := 0
		log.Println("ready for subscribe")
		_, err = conn.Write([]byte("SUB 23482093 PP %d\r\n"))
		if err != nil {
			t.Error(err)
		}
		buf := make([]byte, 65536)
		_, err = conn.Read(buf)
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(buf[:2], resp) {
			t.Error(errors.New("sub failed"))
		}
		READSUBS:
		for {
			n, err := conn.Read(buf)
			if err != nil {
				t.Error(err)
			}
			msgs := strings.Split(string(buf[:n]), "\r\n")
			for _, msg := range msgs {
				if msg == "" {
					break
				}
				cnt ++
				nowSince := time.Since(since)
				
				ns, _ := strconv.Atoi(strings.Split(msg, " ")[3])
				latency += int(time.Now().UnixNano())-ns
				if nowSince> due {
					break READSUBS
				}
			}
		}
		// TODO : calculate cuncurrent latency
		log.Println("subscribe throughput", latency / cnt, "ns", ":cnt", cnt)
		done <- true
	}()
	log.Println("ready for publish")
	buf := make([]byte, 1024)
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