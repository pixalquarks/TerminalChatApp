package chatserver

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

// receive message
func (is *ChatServer) receiveFromStream(csi_ Services_ChatServiceServer, clientUniqueCode_ int, errChannel_ chan error) {

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
				To:                is.getClientsArray(),
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
			SendTo := messageHandleObject.MQue[0].To

			messageHandleObject.mu.Unlock()
			// TODO: when someone enters the chat, the message is sent multiple of time. Fix it later
			if message4Client == "" {
				is.Mu.Lock()
				t := is.Clients[senderUniqueCode]
				t.Name = senderName4Client
				is.NameToUid[t.Name] = t.Uid
				is.Clients[senderUniqueCode] = t
				log.Printf("%v", is.NameToUid)
				is.Mu.Unlock()
				message4Client = fmt.Sprintf("%v has entered the chat", senderName4Client)
				senderName4Client = "server"
			}

			for uid, user := range SendTo {
				if senderUniqueCode != uid {
					err := user.Server.Send(&FromServer{
						Name: senderName4Client,
						Body: message4Client,
					})
					if err != nil {
						errChannel_ <- err
					}
				}
			}

			messageHandleObject.mu.Lock()

			if len(messageHandleObject.MQue) > 1 {
				messageHandleObject.MQue = messageHandleObject.MQue[1:]
			} else {
				messageHandleObject.MQue = []messageUnit{}
			}
			messageHandleObject.mu.Unlock()

		}
		time.Sleep(100 * time.Millisecond)
	}
}
