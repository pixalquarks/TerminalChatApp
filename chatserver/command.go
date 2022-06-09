package chatserver

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"time"
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
		log.Println("MU locked at commandService commnad.go")
		for _, name := range names {
			if usrId, ok := is.NameToUid[name]; ok {
				t := is.Clients[usrId]
				writtenNames = append(writtenNames, t)
			} else {
				errNames = append(errNames, name)
			}
		}
		is.Mu.Unlock()
		log.Println("MU unlocked at commandService command.go")
		id := command.Id
		user, ok := is.Clients[id]
		if !ok {
			return &emptypb.Empty{}, errors.New("error finding the sender name")
		}

		AppendMessage(user.Name, id, msg, time.Now().Unix(), writtenNames)

		if len(errNames) > 0 {
			msg := fmt.Sprintf("Could not find %v names in chat room", errNames)
			AppendMessage("server", "", msg, time.Now().Unix(), []User{user})
		}
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

func (is *ChatServer) CreateClient(_ context.Context, clientName *ClientName) (*CreateClientResponse, error) {
	_, ok := is.NameToUid[clientName.Name]
	log.Println("creating client", clientName.Name)
	if IsRoomFull(is) {
		log.Println("room is full")
		return &CreateClientResponse{
			Created: false,
		}, errors.New("room if full")
	}
	if ok {
		fmt.Printf("from create client : client with name %s already exists\n", clientName.Name)
		return &CreateClientResponse{
			Created: false,
		}, errors.New("client with same name already exists")
	}
	clientUniqueCode := uuid.NewString()
	log.Printf("New user with uniqueID :: %v", clientUniqueCode)
	is.Mu.Lock()
	log.Println("MU locked at CreateClient command.go")
	is.Clients[clientUniqueCode] = User{
		Name: clientName.Name,
		Uid:  clientUniqueCode,
	}
	is.NameToUid[clientName.Name] = clientUniqueCode
	is.Mu.Unlock()
	log.Println("MU unlocked at CreateClient command.go")
	log.Printf("Client added")
	defer fmt.Println("verified")
	return &CreateClientResponse{
		Created:  !ok,
		Id:       clientUniqueCode,
		RoomName: is.Name,
		Delay:    uint32(is.Delay),
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
	log.Println("name verified")
	return &Exists{
		Exists: ok,
	}, nil

}

func (is *ChatServer) RemoveClient(_ context.Context, client *Client) (*emptypb.Empty, error) {
	fmt.Println("removing a client")
	is.removeClient(client.Id)
	return &emptypb.Empty{}, nil
}
