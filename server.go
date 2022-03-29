package main

import (
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"pixalquarks.terminalChatServer/chatserver"
	"sync"
)

func main() {

	Port := os.Getenv("PORT")

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
		Clients: make(map[int]chatserver.Services_ChatServiceServer),
		Mu:      sync.RWMutex{},
	}
	chatserver.RegisterServicesServer(grpcServer, &cs)
	err = grpcServer.Serve(listen)
	if err != nil {
		log.Fatalf("Failed to start gRPC server :: %v", err)
	}
}
