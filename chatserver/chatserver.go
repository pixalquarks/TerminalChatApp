package chatserver

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

type messageUnit struct {
	ClientName        string
	ClientUniqueCode  int
	MessageBody       string
	MessageUniqueCode int
}

type messageHandle struct {
	MQue []messageUnit
	mu   sync.Mutex
}

var messageHandleObject = messageHandle{}

type ChatServer struct {
	Clients map[int]Services_ChatServiceServer
	Mu      sync.RWMutex
}

func (is *ChatServer) addClient(uid int, srv Services_ChatServiceServer) {
	is.Mu.Lock()
	defer is.Mu.Unlock()
	is.Clients[uid] = srv
}

func (is *ChatServer) removeClient(uid int) {
	is.Mu.Lock()
	defer is.Mu.Unlock()
	delete(is.Clients, uid)
}

func (is *ChatServer) getClients() map[int]Services_ChatServiceServer {
	var cs map[int]Services_ChatServiceServer

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
	defer log.Printf("client added successfully")
	log.Println("New chat service created")

	// receive message
	go receiveFromStream(csi, clientUniqueCode, errChannel)
	// send message
	go is.sendToStream(errChannel)

	return <-errChannel
}

// receive message
func receiveFromStream(csi_ Services_ChatServiceServer, clientUniqueCode_ int, errChannel_ chan error) {

	for {
		msg, err := csi_.Recv()
		if err != nil {
			log.Printf("Error receiving message from client :: %v", err)
			errChannel_ <- err
		} else {
			messageHandleObject.mu.Lock()

			messageHandleObject.MQue = append(messageHandleObject.MQue, messageUnit{
				ClientName:        msg.Name,
				ClientUniqueCode:  clientUniqueCode_,
				MessageBody:       msg.Body,
				MessageUniqueCode: rand.Intn(1e8),
			})

			messageHandleObject.mu.Unlock()

			log.Printf("%v", messageHandleObject.MQue[len(messageHandleObject.MQue)-1])
		}
	}
}

// send message
func (is *ChatServer) sendToStream(errChannel_ chan error) {
	for {
		for {
			time.Sleep(500 * time.Millisecond)

			messageHandleObject.mu.Lock()

			if len(messageHandleObject.MQue) == 0 {
				messageHandleObject.mu.Unlock()
				break
			}

			senderUniqueCode := messageHandleObject.MQue[0].ClientUniqueCode
			senderName4Client := messageHandleObject.MQue[0].ClientName
			message4Client := messageHandleObject.MQue[0].MessageBody

			messageHandleObject.mu.Unlock()

			for uid, ss := range is.Clients {
				if senderUniqueCode != uid {
					err := ss.Send(&FromServer{
						Name: senderName4Client,
						Body: message4Client,
					})
					if err != nil {
						errChannel_ <- err
					}
					messageHandleObject.mu.Lock()

					if len(messageHandleObject.MQue) > 1 {
						messageHandleObject.MQue = messageHandleObject.MQue[1:]
					} else {
						messageHandleObject.MQue = []messageUnit{}
					}
					messageHandleObject.mu.Unlock()
				}
			}

		}
		time.Sleep(100 * time.Millisecond)
	}
}
