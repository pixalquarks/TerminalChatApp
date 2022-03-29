package main

import (
	"bufio"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"os"
	"pixalquarks.terminalChatServer/chatserver"
	"strings"
)

func main() {

	fmt.Println("Enter the server IP:Port ::: ")
	reader := bufio.NewReader(os.Stdin)
	serverID, err := reader.ReadString('\n')

	if err != nil {
		log.Printf("Failed to read from console :: %v", err)

	}
	serverID = strings.Trim(serverID, "\r\n")

	log.Println("Connecting to : " + serverID)

	// call chatService to create a stream
	conn, err := grpc.Dial(serverID, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Failed to connect to gRPC server :: %v", err)
	}

	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("Could not close the connection :: %v", err)
		}
	}()

	client := chatserver.NewServicesClient(conn)

	stream, err := client.ChatService(context.Background())

	if err != nil {
		log.Fatalf("Failed to call ChatService :: %v", err)
	}

	// implement communication with gRPC server
	ch := clientHandle{stream: stream}
	ch.clientConfig()
	go ch.sendMessage()
	go ch.receiveMessage()

	// blocker
	bl := make(chan bool)
	<-bl
}

type clientHandle struct {
	stream     chatserver.Services_ChatServiceClient
	clientName string
}

func (ch *clientHandle) clientConfig() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Your Name : ")
	name, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Failed to read from console :: %v", err)
	}
	ch.clientName = strings.Trim(name, "\r\n")
}

// send message
func (ch *clientHandle) sendMessage() {

	for {
		reader := bufio.NewReader(os.Stdin)
		clientMessage, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Failed to read from console :: %v, message: %v", err, clientMessage)
		}
		clientMessage = strings.Trim(clientMessage, "\r\n")

		clientMessageBox := &chatserver.FromClient{
			Name: ch.clientName,
			Body: clientMessage,
		}
		err = ch.stream.Send(clientMessageBox)

		if err != nil {
			log.Printf("Error while sending message :: %v", err)
		}
	}
}

// receive message
func (ch *clientHandle) receiveMessage() {
	for {
		msg, err := ch.stream.Recv()
		if err != nil {
			log.Printf("Error while receiving message :: %v", err)
		}

		fmt.Printf("%s :: %s \n", msg.Name, msg.Body)
	}
}
