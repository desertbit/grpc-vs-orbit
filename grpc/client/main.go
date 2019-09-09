/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a client for Greeter service.
package main

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
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
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return
	}
	defer conn.Close()

	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	start := time.Now()

	for i := 0; i < 300; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
		if err != nil {
			return err
		}

		_ = r
		//log.Printf("Greeting: %s", r.GetMessage())
	}

	elapsed := time.Since(start)
	totalCallsMx.Lock()
	totalCalls += elapsed
	totalCallsMx.Unlock()

	return
}
