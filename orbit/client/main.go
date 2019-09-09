package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/desertbit/orbit"
	"github.com/desertbit/orbit-vs-grpc/orbit/api"
)

const (
	remoteAddr = "127.0.0.1:40150"
)

var (
	totalCallsMx sync.Mutex
	totalCalls   time.Duration
)

func main() {
	start := time.Now()

	var wg sync.WaitGroup
	wg.Add(500)
	for i := 0; i < 500; i++ {
		go func() {
			err := do()
			if err != nil {
				log.Fatalln(err)
				return
			}
			wg.Done()
		}()
	}

	wg.Wait()

	elapsed := time.Since(start)
	log.Printf("calls took %s", totalCalls)
	log.Printf("total took %s", elapsed)
}

func do() (err error) {
	conn, err := net.Dial("tcp", remoteAddr)
	if err != nil {
		return
	}
	defer conn.Close()

	s, err := orbit.ClientSession(conn, &orbit.Config{})
	if err != nil {
		return
	}
	defer s.Close()

	ctrl, _, err := s.Init(&orbit.Init{})
	if err != nil {
		return fmt.Errorf("0 %v", err)
	}

	ctrl.Ready()

	start := time.Now()

	for i := 0; i < 300; i++ {
		ctx, err := ctrl.Call(api.SayHello, &api.HelloRequest{
			Name: "world",
		})
		if err != nil {
			return fmt.Errorf("1 %v", err)
		}

		var ret api.HelloReply
		err = ctx.Decode(&ret)
		if err != nil {
			return fmt.Errorf("2 %v", err)
		}

		//log.Printf("Greeting: %s", ret.Message)
	}

	elapsed := time.Since(start)
	totalCallsMx.Lock()
	totalCalls += elapsed
	totalCallsMx.Unlock()

	return
}
