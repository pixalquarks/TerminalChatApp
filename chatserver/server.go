package chatserver

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
)

type messageUnit struct {
	ClientName        string
	ClientUniqueCode  int
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
	Uid    int
	Server Services_ChatServiceServer
}

type ChatServer struct {
	Clients   map[int]User
	NameToUid map[string]int
	Mu        sync.RWMutex
}

func (is *ChatServer) addClient(uid int, srv Services_ChatServiceServer) {
	is.Mu.Lock()
	defer is.Mu.Unlock()
	log.Println("adding new client", uid)
	is.Clients[uid] = User{
		Name:   "",
		Uid:    uid,
		Server: srv,
	}
}

func (is *ChatServer) removeClient(uid int) {
	is.Mu.Lock()
	t := is.Clients[uid]
	name := t.Name
	delete(is.NameToUid, is.Clients[uid].Name)
	delete(is.Clients, uid)
	is.Mu.Unlock()
	log.Println(name)
	//
	messageHandleObject.mu.Lock()
	messageHandleObject.MQue = append(messageHandleObject.MQue, messageUnit{
		ClientName:        "server",
		ClientUniqueCode:  uid,
		MessageBody:       fmt.Sprintf("%v left has the chat", name),
		MessageUniqueCode: rand.Intn(1e8),
		To:                is.getClientsArray(),
	})
	//
	messageHandleObject.mu.Unlock()

}

func (is *ChatServer) getClientsArray() []User {
	arr := make([]User, 0)
	for _, user := range is.getClients() {
		arr = append(arr, user)
	}
	return arr
}

func (is *ChatServer) getClients() map[int]User {
	cs := make(map[int]User)

	is.Mu.RLock()
	defer is.Mu.RUnlock()
	for k, v := range is.Clients {
		cs[k] = v
	}
	return cs
}

func (is *ChatServer) mustEmbedUnimplementedServicesServer() {
	//TODO implement me
	panic("implement me")
}

func (is *ChatServer) ChatService(csi Services_ChatServiceServer) error {
	clientUniqueCode := rand.Intn(1e6)
	log.Printf("New user with uniqueID :: %v", clientUniqueCode)
	errChannel := make(chan error)

	is.addClient(clientUniqueCode, csi)
	defer is.removeClient(clientUniqueCode)
	defer log.Printf("client removed successfully successfully")
	log.Println("New chat service created")

	// receive message
	go is.receiveFromStream(csi, clientUniqueCode, errChannel)
	// send message
	go is.sendToStream(errChannel)

	return <-errChannel
}
