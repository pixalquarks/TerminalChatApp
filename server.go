package main

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"pixalquarks.terminalChatServer/chatserver"
	"sync"
)

func main() {
	servVars, err := GetServerVariables()
	if err != nil {
		fmt.Println(err)
		return
	}
	Port := servVars.Port

	if Port == "" {
		Port = "5000" // default port
	}

	// initialize listener
	listen, err := net.Listen("tcp", ":"+Port)
	if err != nil {
		log.Fatalf("Could not listen @ %v :: %v", Port, err)
	}
	log.Println("Listening @ : " + Port)

	// gRPC server instance
	grpcServer := grpc.NewServer()

	// register chat server
	cs := chatserver.ChatServer{
		Name:      servVars.RoomName,
		RoomSize:  servVars.RoomSize,
		Clients:   make(map[string]chatserver.User),
		NameToUid: make(map[string]string),
		Mu:        sync.RWMutex{},
		Delay:     servVars.Delay,
	}
	chatserver.RegisterServicesServer(grpcServer, &cs)
	go cs.CheckOnClients()
	err = grpcServer.Serve(listen)
	if err != nil {
		log.Fatalf("Failed to start gRPC server :: %v", err)
	}
}
