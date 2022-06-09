package chatserver

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type messageUnit struct {
	ClientName        string
	ClientUniqueCode  string
	MessageBody       string
	MessageUniqueCode int
	TimeStamp         int64
	To                []User
}

type messageHandle struct {
	MQue []messageUnit
	mu   sync.Mutex
}

var messageHandleObject = messageHandle{}

type User struct {
	Name    string
	Uid     string
	Server  Services_ChatServiceServer
	CanSend bool
}

type ChatServer struct {
	Name      string
	RoomSize  uint
	Clients   map[string]User
	NameToUid map[string]string
	Mu        sync.RWMutex
	Delay     uint
}

func (is *ChatServer) mustEmbedUnimplementedServicesServer() {
	//TODO implement me
	panic("implement me")
}

func (is *ChatServer) ChatService(req *StreamRequest, csi Services_ChatServiceServer) error {
	errChannel := make(chan error)
	is.Mu.Lock()
	log.Println("MU locked at ChatService server.go")
	id := req.Id
	name := is.Clients[id].Name
	is.Mu.Unlock()
	log.Println("MU locked at ChatService server.go")
	is.addClient(id, csi)
	defer is.removeClient(id)
	defer log.Printf("client removed successfully successfully")
	log.Println("New chat service created")

	message := fmt.Sprintf("%v has entered the chat", name)
	AppendMessage("server", "", message, time.Now().Unix(), is.getClientsArray())

	// send message
	go is.sendToStream(errChannel)

	return <-errChannel
}
