package chatserver

import (
	"context"
	"errors"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

func (is *ChatServer) SendMessage(_ context.Context, msg *FromClient) (*emptypb.Empty, error) {
	is.Mu.Lock()
	id := msg.Id
	user, ok := is.Clients[id]
	if !ok {
		return &emptypb.Empty{}, errors.New("couldn't find the name")
	}
	is.Mu.Unlock()
	AppendMessage(user.Name, msg.Id, msg.Body, is.getClientsArray())

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
			SendTo := messageHandleObject.MQue[0].To

			messageHandleObject.mu.Unlock()

			for _, user := range SendTo {
				if senderUniqueCode != user.Uid && user.Server != nil {
					if user.Server.Context().Err() == context.Canceled {
						is.removeClient(user.Uid)
						continue
					}
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
