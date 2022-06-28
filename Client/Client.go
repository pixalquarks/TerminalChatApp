package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"pixalquarks.terminalChatServer/chatserver"
)

type clientHandle struct {
	client     chatserver.ServicesClient
	stream     chatserver.Services_ChatServiceClient
	clientName string
	roomName   string
	delay      int
	uid        string
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
	if client == nil {
		return nil, errors.New("could not create client")
	}
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

func (ch *clientHandle) GetUserNames() ([]string, error) {
	clientsArr, err := ch.client.GetClients(context.Background(), &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	arr := make([]string, 0)
	for _, v := range clientsArr.Client {
		arr = append(arr, v.Name)
	}
	return arr, nil
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
			fmt.Println(err.Error())
		} else {
			if t {
				resp, err := ch.client.VerifyName(context.Background(), &chatserver.ClientName{
					Name: name,
				})
				if err != nil {
					return "", err
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
				return err
			}
			if resp.Created {
				ch.clientName = name
				ch.uid = resp.Id
				ch.roomName = resp.RoomName
				ch.delay = int(resp.Delay)
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
		Id:   ch.uid,
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
				Id:    ch.uid,
			}); err != nil {
				if er := LogOut(err.Error()); er != nil {
					log.Println(er)
				}
			} else {
				if err := LogOut("command executed successfully"); err != nil {
					log.Println(err)
				}
			}
		}
	} else {

		clientMessageBox := &chatserver.FromClient{
			Id:        ch.uid,
			Body:      msg,
			Timestamp: time.Now().Unix(),
		}
		_, err := ch.client.SendMessage(context.Background(), clientMessageBox)

		if err != nil {
			return errors.New(fmt.Sprintf("Error while sending message :: %v", err))
		}
	}
	return nil
}

// receive message
func (ch *clientHandle) receiveMessage(out func(sender string, message string, timeStamp int64)) {
	for {
		msg, err := ch.stream.Recv()
		if err != nil {
			log.Printf("Error while receiving message :: %v", err)
			return
		}
		if msg != nil {
			if msg.Name == "server" {
				UpdateUserList()
				//m := fmt.Sprintf("***** System Message : %s *****\n", msg.Body)
				out(msg.Name, msg.Body, msg.Timestamp)
			} else {
				//m := fmt.Sprintf("%s :: %s \n", msg.Name, msg.Body)
				out(msg.Name, msg.Body, msg.Timestamp)
			}
		}
	}
}
