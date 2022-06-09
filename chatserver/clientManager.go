package chatserver

import (
	"context"
	"fmt"
	"log"
	"time"
)

func (is *ChatServer) addClient(uid string, srv Services_ChatServiceServer) {
	is.Mu.Lock()
	log.Println("Mu locked at addClient clientManager.go")
	defer func() {
		is.Mu.Unlock()
		log.Println("MU unlocked at addClient clientManager.go")
	}()
	log.Println("adding new client", uid)
	t := is.Clients[uid]
	t.Server = srv
	is.Clients[uid] = t
	fmt.Println(is.Clients)
}

func (is *ChatServer) removeClient(uid string) {
	is.Mu.Lock()
	log.Println("MU locked at removeClient clientManager.go")
	t := is.Clients[uid]
	name := t.Name
	delete(is.NameToUid, is.Clients[uid].Name)
	delete(is.Clients, uid)
	is.Mu.Unlock()
	log.Println("MU unlocked at removeClient clientManager.go")
	log.Println(name, "left the chat")
	//
	msg := fmt.Sprintf("%v left has the chat", name)
	AppendMessage("server", "", msg, time.Now().Unix(), is.getClientsArray())

}

func (is *ChatServer) getClientsArray() []User {
	arr := make([]User, 0)
	for _, user := range is.getClients() {
		arr = append(arr, user)
	}
	return arr
}

func (is *ChatServer) getClients() map[string]User {
	cs := make(map[string]User)

	is.Mu.RLock()
	//log.Println("MU locked at getClients clientManager.go")
	defer func() {
		is.Mu.RUnlock()
		//log.Println("MU unlocked at getClients clientManager.go")
	}()

	for k, v := range is.Clients {
		cs[k] = v
	}
	return cs
}

func (is *ChatServer) CheckOnClients() {
	for {
		time.Sleep(time.Millisecond * 100)
		for _, user := range is.getClientsArray() {
			if &user == nil {
				continue
			}
			if &user == nil || user.Server.Context().Err() == context.Canceled {
				is.removeClient(user.Uid)
				continue
			}
		}
	}
}
