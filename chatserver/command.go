package chatserver

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"math/rand"
)

func (is *ChatServer) CommandService(_ context.Context, command *Command) (*emptypb.Empty, error) {
	cmd := rune(command.Type)
	switch cmd {
	case 'p', 'P':
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
		id := command.Id
		user, ok := is.Clients[id]
		if !ok {
			return &emptypb.Empty{}, errors.New("error finding the sender name")
		}
		messageHandleObject.MQue = append(messageHandleObject.MQue, messageUnit{
			ClientName:        user.Name,
			ClientUniqueCode:  id,
			MessageBody:       msg,
			MessageUniqueCode: rand.Intn(1e8),
			To:                writtenNames,
		})
		if len(errNames) > 0 {
			messageHandleObject.MQue = append(messageHandleObject.MQue, messageUnit{
				ClientName:        "server",
				ClientUniqueCode:  id,
				MessageBody:       fmt.Sprintf("Could not find %v names in chat room", errNames),
				MessageUniqueCode: rand.Intn(1e8),
				To: []User{
					user,
				},
			})
		}
		messageHandleObject.mu.Unlock()
	default:
		return &emptypb.Empty{}, errors.New("no such command exists")
	}

	return &emptypb.Empty{}, nil
}

func (is *ChatServer) GetClients(_ context.Context, _ *emptypb.Empty) (*Clients, error) {
	clientsMap := is.getClients()
	l := len(clientsMap)
	clientArr := make([]*Client, l)
	count := 0

	for _, v := range clientsMap {
		clientArr[count] = &Client{
			Name: v.Name,
			Id:   v.Uid,
		}
		count++
	}
	clients := &Clients{
		Client: clientArr,
		Count:  uint32(count),
	}

	return clients, nil
}

func (is *ChatServer) CreateClient(_ context.Context, clientName *ClientName) (*ClientNameResponse, error) {
	_, ok := is.NameToUid[clientName.Name]
	if ok {
		fmt.Printf("from create client : client with name %s already exists\n", clientName.Name)
		return &ClientNameResponse{
			Exists: ok,
			Id:     -1,
		}, nil
	}
	clientUniqueCode := rand.Int31n(1e6)
	log.Printf("New user with uniqueID :: %v", clientUniqueCode)
	is.Mu.Lock()
	is.Clients[clientUniqueCode] = User{
		Name: clientName.Name,
		Uid:  clientUniqueCode,
	}
	is.NameToUid[clientName.Name] = clientUniqueCode
	is.Mu.Unlock()
	defer fmt.Println("verified")
	return &ClientNameResponse{
		Exists: ok,
		Id:     clientUniqueCode,
	}, nil
}

func (is *ChatServer) VerifyName(_ context.Context, name *ClientName) (*Exists, error) {
	nameOk, err := VerifyNameCharacters(name.Name)
	if err != nil {
		return &Exists{
			Exists: true,
		}, err
	}
	if !nameOk {
		return &Exists{
			Exists: false,
		}, errors.New("name should only contain alphanumeric characters and underscore")
	}
	_, ok := is.NameToUid[name.Name]
	return &Exists{
		Exists: ok,
	}, nil

}

func (is *ChatServer) RemoveClient(_ context.Context, client *Client) (*emptypb.Empty, error) {
	fmt.Println("removing a client")
	is.removeClient(client.Id)
	return &emptypb.Empty{}, nil
}
