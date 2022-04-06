package chatserver

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"math/rand"
)

func (is *ChatServer) CommandService(ctx context.Context, command *Command) (*emptypb.Empty, error) {
	cmd := rune(command.Type)
	switch cmd {
	case 'p', 'P':
		log.Println(command.Value)
		names, msg, err := GetNamesFromCommandString(command.Value)
		if err != nil {
			return &emptypb.Empty{}, err
		}
		writtenNames := make([]User, 0)
		errNames := make([]string, 0)
		is.Mu.Lock()
		for _, name := range names {
			if usrId, ok := is.NameToUid[name]; ok {
				t := is.Clients[usrId]
				writtenNames = append(writtenNames, t)
			} else {
				errNames = append(errNames, name)
			}
		}
		is.Mu.Unlock()
		messageHandleObject.mu.Lock()
		clientUid, ok := is.NameToUid[command.Client]
		if !ok {
			return &emptypb.Empty{}, errors.New("error finding the sender name")
		}
		messageHandleObject.MQue = append(messageHandleObject.MQue, messageUnit{
			ClientName:        command.Client,
			ClientUniqueCode:  clientUid,
			MessageBody:       msg,
			MessageUniqueCode: rand.Intn(1e8),
			To:                writtenNames,
		})
		if len(errNames) > 0 {
			messageHandleObject.MQue = append(messageHandleObject.MQue, messageUnit{
				ClientName:        "server",
				ClientUniqueCode:  clientUid,
				MessageBody:       fmt.Sprintf("Could not find %v names in chat room", errNames),
				MessageUniqueCode: rand.Intn(1e8),
				To: []User{
					is.Clients[clientUid],
				},
			})
		}
		messageHandleObject.mu.Unlock()
	default:
		return &emptypb.Empty{}, errors.New("no such command exists")
	}

	return &emptypb.Empty{}, nil
}

func (is *ChatServer) GetClients(context.Context, *emptypb.Empty) (*Clients, error) {
	log.Printf("GetClients command \n")
	clientsMap := is.getClients()
	log.Println(clientsMap)
	l := len(clientsMap)
	clientArr := make([]*Client, l)
	count := 0

	for _, v := range clientsMap {
		clientArr[count] = &Client{
			Name: v.Name,
			Id:   int32(v.Uid),
		}
		count++
	}
	log.Println(clientArr)
	clients := &Clients{
		Client: clientArr,
		Count:  uint32(count),
	}

	return clients, nil
}

func (is *ChatServer) VerifyName(ctx context.Context, clientName *ClientName) (*ClientNameResponse, error) {
	log.Println(is.NameToUid)
	_, ok := is.NameToUid[clientName.Name]
	log.Println(ok)
	defer fmt.Println("verified")
	return &ClientNameResponse{
		Exists: ok,
	}, nil
}
