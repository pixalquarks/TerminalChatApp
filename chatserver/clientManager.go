package chatserver

import (
	"context"
	"fmt"
	"log"
	"time"
)

func (is *ChatServer) addClient(uid int32, srv Services_ChatServiceServer) {
	is.Mu.Lock()
	defer is.Mu.Unlock()
	log.Println("adding new client", uid)
	t := is.Clients[uid]
	t.Server = srv
	is.Clients[uid] = t
	fmt.Println(is.Clients)
}

func (is *ChatServer) removeClient(uid int32) {
	is.Mu.Lock()
	t := is.Clients[uid]
	name := t.Name
	delete(is.NameToUid, is.Clients[uid].Name)
	delete(is.Clients, uid)
	is.Mu.Unlock()
	log.Println(name, "left the chat")
	//
	msg := fmt.Sprintf("%v left has the chat", name)
	AppendMessage("server", -1, msg, is.getClientsArray())

}

func (is *ChatServer) getClientsArray() []User {
	arr := make([]User, 0)
	for _, user := range is.getClients() {
		arr = append(arr, user)
	}
	return arr
}

func (is *ChatServer) getClients() map[int32]User {
	cs := make(map[int32]User)

	is.Mu.RLock()
	defer is.Mu.RUnlock()
	for k, v := range is.Clients {
		cs[k] = v
	}
	return cs
}

func (is *ChatServer) CheckOnClients() {
	for {
		time.Sleep(time.Millisecond * 100)
		for _, user := range is.getClientsArray() {
			if user.Server.Context().Err() == context.Canceled {
				is.removeClient(user.Uid)
				continue
			}
		}
	}
}
