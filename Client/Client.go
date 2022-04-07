package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"os"
	"pixalquarks.terminalChatServer/chatserver"
	"strings"
)

type clientHandle struct {
	client     chatserver.ServicesClient
	stream     chatserver.Services_ChatServiceClient
	clientName string
	uid        int32
}

func GetConnection() (*grpc.ClientConn, error) {
	fmt.Println("Enter the server IP:Port ::: ")
	reader := bufio.NewReader(os.Stdin)
	serverID, err := reader.ReadString('\n')
	//
	if err != nil {
		log.Printf("Failed to read from console :: %v", err)
		return nil, err
	}
	serverID = strings.Trim(serverID, "\r\n")
	//
	log.Println("Connecting to : " + serverID)

	conn, err := grpc.Dial(serverID, grpc.WithInsecure())
	return conn, err
}

func CreateClient(conn *grpc.ClientConn) (*clientHandle, error) {
	client := chatserver.NewServicesClient(conn)
	// implement communication with gRPC server
	ch := clientHandle{client: client}
	name, err := ch.GetName()
	if err != nil {
		return nil, err
	}
	if err := ch.Config(name); err != nil {
		return nil, err
	}
	stream, err := client.ChatService(context.Background(), &chatserver.StreamRequest{
		Id: ch.uid,
	})
	if err != nil {
		return nil, err
	}
	ch.stream = stream
	return &ch, nil
}

func (ch *clientHandle) GetName() (string, error) {
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
					return name, nil
				} else {
					fmt.Println("This name is already taken")
				}
			} else {
				fmt.Println("Name should only contain alphanumeric characters and underscore")
			}
		}
	}
}

func (ch *clientHandle) Config(name string) error {
	if t, err := VerifyNameCharacters(name); err != nil {
		return errors.New("error while verifying name, please try again")
	} else {
		if t {
			resp, err := ch.client.CreateClient(context.Background(), &chatserver.ClientName{
				Name: name,
			})
			if err != nil {
				return errors.New("error while verifying name, please try again")
			}
			if !resp.Exists {
				ch.clientName = name
				ch.uid = resp.Id
			} else {
				fmt.Println("This name is already taken")
			}
		} else {
			fmt.Println("Name should only contain alphanumeric characters and underscore")
		}
	}
	return nil
}

func (ch *clientHandle) quit() error {
	_, err := ch.client.RemoveClient(context.Background(), &chatserver.Client{
		Name: ch.clientName,
		Id:   int32(ch.uid),
	})
	if err != nil {
		return err
	}
	return nil
}

// send message
func (ch *clientHandle) sendMessage(msg string) error {
	if msg[0] == '!' {
		command := msg[1]
		if command == 'L' || command == 'l' {
			fmt.Println("command")
			if t, err := ch.client.GetClients(context.Background(), &emptypb.Empty{}); err != nil {
				return errors.New(fmt.Sprintf("Error while executing command :: %v", err))
			} else {
				res := ""
				for _, v := range t.Client {
					res += v.Name + "\t"
				}
				fmt.Println(res)
			}
		} else if command == 'P' || command == 'p' {
			if _, err := ch.client.CommandService(context.Background(), &chatserver.Command{
				Type:  uint32(command),
				Value: msg[2:],
				Id:    int32(ch.uid),
			}); err != nil {
				log.Printf("Error while executing commnad")
			} else {
				fmt.Println("Command executed successfully")
			}
		}
	} else {

		clientMessageBox := &chatserver.FromClient{
			Id:   int32(ch.uid),
			Body: msg,
		}
		_, err := ch.client.SendMessage(context.Background(), clientMessageBox)

		if err != nil {
			return errors.New(fmt.Sprintf("Error while sending message :: %v", err))
		}
	}
	return nil
}

// receive message
func (ch *clientHandle) receiveMessage(out func(message string)) {
	for {
		msg, err := ch.stream.Recv()
		if err != nil {
			log.Printf("Error while receiving message :: %v", err)
			return
		}
		if msg != nil {
			if msg.Name == "server" {
				m := fmt.Sprintf("***** System Message : %s *****\n", msg.Body)
				out(m)
			} else {
				m := fmt.Sprintf("%s :: %s \n", msg.Name, msg.Body)
				out(m)
			}
		}
	}
}
