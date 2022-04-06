package chatserver

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// TODO: parsing is unable to separate names from message

func GetNamesFromCommandString(command string) ([]string, string, error) {
	r, err := regexp.Compile("^[[A-Za-z0-9_,]+]")
	if err != nil {
		return make([]string, 0), "", errors.New("could not compile regexp")
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

func VerifyName(name string) (bool, error) {
	r, err := regexp.Compile("^[A-Za-z0-9_]+$")
	if err != nil {
		return false, errors.New("error while compiling regexp")
	}
	return r.MatchString(name), nil
}

func GetNamesFromClientStruct(clientStructs map[int]User) []string {
	names := make([]string, len(clientStructs))
	i := 0
	for _, user := range clientStructs {
		names[i] = user.Name
		i++
	}
	return names
}

type Set struct {
	store map[uint32]bool
}

func (s *Set) Add(key uint32) {
	s.store[key] = true
}

func (s *Set) Remove(key uint32) {
	delete(s.store, key)
}

func (s *Set) Contains(key uint32) bool {
	return s.store[key] == true
}

func (s *Set) Size() int {
	return len(s.store)
}

func (s *Set) Iter() map[uint32]bool {
	return s.store
}
