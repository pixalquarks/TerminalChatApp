package main

import (
	"bufio"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"os"
	"pixalquarks.terminalChatServer/chatserver"
	"strings"
)

func main() {

	fmt.Println("Enter the server IP:Port ::: ")
	reader := bufio.NewReader(os.Stdin)
	serverID, err := reader.ReadString('\n')
	//
	if err != nil {
		log.Printf("Failed to read from console :: %v", err)
	}
	serverID = strings.Trim(serverID, "\r\n")
	//
	log.Println("Connecting to : " + serverID)
	//
	//// call chatService to create a stream
	conn, err := grpc.Dial(serverID, grpc.WithInsecure())
	//
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server :: %v", err)
	}
	//
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("Could not close the connection :: %v", err)
		}
	}()
	//
	client := chatserver.NewServicesClient(conn)
	stream, err := client.ChatService(context.Background())

	if err != nil {
		log.Fatalf("Failed to call ChatService :: %v", err)
	}

	// implement communication with gRPC server
	ch := clientHandle{stream: stream, client: client}
	ch.clientConfig()
	go ch.sendMessage()
	go ch.receiveMessage()
	//
	//// blocker
	bl := make(chan bool)
	<-bl
}

type clientHandle struct {
	client     chatserver.ServicesClient
	stream     chatserver.Services_ChatServiceClient
	clientName string
}

func (ch *clientHandle) clientConfig() {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("Your Name : ")
		name, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Failed to read from console :: %v", err)
		}
		name = strings.Trim(name, "\r\n")
		if t, err := VerifyNameCharacters(name); err != nil {
			fmt.Println("Error while verifying name, please try again")
		} else {
			if t {
				resp, err := ch.client.VerifyName(context.Background(), &chatserver.ClientName{
					Name: name,
				})
				if err != nil {
					fmt.Println("Error while verifying name, please try again")
				}
				if !resp.Exists {
					ch.clientName = strings.Trim(name, "\r\n")
					break
				} else {
					fmt.Println("This name is already taken")
				}
			} else {
				fmt.Println("Name should only contain alphanumeric characters and underscore")
			}
		}

	}

	clientMessage := ""
	clientMessageBox := &chatserver.FromClient{
		Name: ch.clientName,
		Body: clientMessage,
	}
	err := ch.stream.Send(clientMessageBox)

	if err != nil {
		log.Printf("Error while sending message :: %v", err)
	}
}

// send message
func (ch *clientHandle) sendMessage() {

	for {
		reader := bufio.NewReader(os.Stdin)
		clientMessage, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Failed to read from console :: %v, message: %v", err, clientMessage)
		}
		msg := strings.Trim(clientMessage, "\r\n")
		fmt.Println(msg)

		if msg[0] == '!' {
			command := msg[1]
			if command == 'L' || command == 'l' {
				fmt.Println("command")
				if t, err := ch.client.GetClients(context.Background(), &emptypb.Empty{}); err != nil {
					log.Printf("Error while executing command :: %v", err)
				} else {
					fmt.Println("received message")
					res := ""
					for _, v := range t.Client {
						res += v.Name + "\t"
					}
					fmt.Println(res)
				}
			} else if command == 'P' || command == 'p' {
				if _, err := ch.client.CommandService(context.Background(), &chatserver.Command{
					Type:   uint32(command),
					Value:  msg[2:],
					Client: ch.clientName,
				}); err != nil {
					log.Printf("Error while executing commnad")
				} else {
					fmt.Println("Command executed successfully")
				}
			}
		} else {

			clientMessageBox := &chatserver.FromClient{
				Name: ch.clientName,
				Body: msg,
			}
			err = ch.stream.Send(clientMessageBox)

			if err != nil {
				log.Printf("Error while sending message :: %v", err)
			}
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
		if msg.Name == "server" {
			fmt.Printf("***** System Message : %s *****\n", msg.Body)
		} else {
			fmt.Printf("%s :: %s \n", msg.Name, msg.Body)
		}
	}
}
