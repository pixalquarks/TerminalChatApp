package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

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
		return make([]string, 0), "", errors.New("no names found in the command, make sure the command is entered correctly example \"[name1,name2]\"")
	}
	names = names[1 : len(names)-1]
	namesArr := strings.Split(names, ",")
	return namesArr, command[t[1]:], nil
}

func VerifyNameCharacters(name string) (bool, error) {
	r, err := regexp.Compile("^[A-Za-z0-9_]+$")
	if err != nil {
		return false, errors.New("error while compiling regexp")
	}
	return r.MatchString(name), nil
}
