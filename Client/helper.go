package main

import (
	"errors"
	"regexp"
)

func VerifyNameCharacters(name string) (bool, error) {
	r, err := regexp.Compile("^[A-Za-z0-9_]+$")
	if err != nil {
		return false, errors.New("error while compiling regexp")
	}
	return r.MatchString(name), nil
}
