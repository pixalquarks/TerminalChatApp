package chatserver

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

func (is *ChatServer) SendMessage(_ context.Context, msg *FromClient) (*emptypb.Empty, error) {
	var id string
	var user User
	err := func() error {
		is.Mu.Lock()
		defer is.Mu.Unlock()
		fmt.Println("MU locked in SendMessage chat.go")
		id = msg.Id
		var ok bool
		user, ok = is.Clients[id]
		if !ok {
			return errors.New("couldn't find the name")
		}
		return nil
	}()
	
	if err != nil {
		return &emptypb.Empty{}, err
	}
	fmt.Println("MU unlocked in SendMessage chat.go")
	AppendMessage(user.Name, msg.Id, msg.Body, msg.Timestamp, is.getClientsArray())

	return &emptypb.Empty{}, nil
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
			timeStamp := messageHandleObject.MQue[0].TimeStamp
			SendTo := messageHandleObject.MQue[0].To

			messageHandleObject.mu.Unlock()

			for _, user := range SendTo {
				if senderUniqueCode != user.Uid && user.Server != nil {
					if user.Server.Context().Err() == context.Canceled {
						is.removeClient(user.Uid)
						continue
					}
					err := user.Server.Send(&FromServer{
						Name:      senderName4Client,
						Body:      message4Client,
						Timestamp: timeStamp,
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
