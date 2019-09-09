package main

import (
	"fmt"
	"log"
	"net"
	"runtime"

	"github.com/desertbit/orbit"
	"github.com/desertbit/orbit-vs-grpc/orbit/api"
	"github.com/desertbit/orbit/control"
)

const (
	listenAddr = ":40150"
)

func main() {
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalln(err)
		return
	}

	s := orbit.NewServer(ln, &orbit.ServerConfig{
		NewConnNumberWorkers: runtime.NumCPU(),
	})
	defer s.Close_()

	go handleNewSessionRoutine(s)

	err = s.Listen()
	if err != nil {
		log.Fatalln(err)
		return
	}
	return
}

func handleNewSessionRoutine(s *orbit.Server) {
	defer s.Close()

	var (
		closingChan    = s.ClosingChan()
		newSessionChan = s.NewSessionChan()
	)

	for {
		select {
		case <-closingChan:
			return

		case session := <-newSessionChan:
			err := newSession(session)
			if err != nil {
				fmt.Printf("handleNewSessionRoutine: %v\n", err)
			}
		}
	}
}

func newSession(s *orbit.Session) (err error) {
	// Always close the session on error.
	defer func() {
		if err != nil {
			s.Close()
		}
	}()

	ctrl, _, err := s.Init(&orbit.Init{
		Control: orbit.InitControl{
			Funcs: control.Funcs{
				api.SayHello: sayHello,
			},
		},
	})
	if err != nil {
		return
	}

	ctrl.Ready()
	return
}

func sayHello(ctx *control.Context) (v interface{}, err error) {
	var args api.HelloRequest
	err = ctx.Decode(&args)
	if err != nil {
		return
	}

	//log.Printf("Received: %v", args.Name)

	return &api.HelloReply{
		Message: "Hello " + args.Name,
	}, nil
}
