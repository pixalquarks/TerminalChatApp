package chatserver

import (
	"fmt"
	"log"
	"sync"
)

type messageUnit struct {
	ClientName        string
	ClientUniqueCode  int32
	MessageBody       string
	MessageUniqueCode int
	To                []User
}

type messageHandle struct {
	MQue []messageUnit
	mu   sync.Mutex
}

var messageHandleObject = messageHandle{}

type User struct {
	Name   string
	Uid    int32
	Server Services_ChatServiceServer
}

type ChatServer struct {
	Clients   map[int32]User
	NameToUid map[string]int32
	Mu        sync.RWMutex
}

func (is *ChatServer) mustEmbedUnimplementedServicesServer() {
	//TODO implement me
	panic("implement me")
}

func (is *ChatServer) ChatService(req *StreamRequest, csi Services_ChatServiceServer) error {
	errChannel := make(chan error)
	is.Mu.Lock()
	id := req.Id
	name := is.Clients[id].Name
	is.Mu.Unlock()
	is.addClient(id, csi)
	defer is.removeClient(id)
	defer log.Printf("client removed successfully successfully")
	log.Println("New chat service created")

	message := fmt.Sprintf("%v has entered the chat", name)
	AppendMessage("server", -1, message, is.getClientsArray())

	// send message
	go is.sendToStream(errChannel)

	return <-errChannel
}
