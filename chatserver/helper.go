package chatserver

import (
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
)

func GetNamesFromCommandString(command string) ([]string, string, error) {
	r, err := regexp.Compile("^[[A-Za-z0-9_,]+]")
	if err != nil {
		return make([]string, 0), "", err
	}
	match := r.MatchString(command)
	if !match {
		return make([]string, 0), "", errors.New("no names found in the command")
	}
	t := r.FindStringIndex(command)

	if t[0] < 0 || t[1] >= len(command) || t[0] == t[1] {
		return make([]string, 0), "", errors.New("no names found in the command")
	}
	names := command[t[0]:t[1]]
	fmt.Println(t)

	if names == "" || names == " " {
		return make([]string, 0), "", errors.New("no names found in the command")
	}
	names = names[1 : len(names)-1]
	namesArr := strings.Split(names, ",")
	return namesArr, command[t[1]:], nil
}

func VerifyNameCharacters(name string) (bool, error) {
	r, err := regexp.Compile("^[A-Za-z0-9_]+$")
	if err != nil {
		return false, err
	}
	return r.MatchString(name), nil
}

func AppendMessage(clientName string, cUniqCode int32, message string, to []User) {
	messageHandleObject.mu.Lock()
	defer messageHandleObject.mu.Unlock()

	messageHandleObject.MQue = append(messageHandleObject.MQue, messageUnit{
		ClientName:        clientName,
		ClientUniqueCode:  cUniqCode,
		MessageBody:       message,
		MessageUniqueCode: rand.Intn(1e8),
		To:                to,
	})
}

func IsRoomFull(serv *ChatServer) bool {
	t := len(serv.Clients)
	fmt.Println(t, serv.RoomSize)
	b := t >= int(serv.RoomSize)
	fmt.Println(b)
	return len(serv.Clients) >= int(serv.RoomSize)
}
