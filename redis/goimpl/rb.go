package main
import (
	"log"
	"time"

	"github.com/go-redis/redis"
)

func main() {
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
	log.Println("set throughput", sum / nSec, "msg/sec")

}