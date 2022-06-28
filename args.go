package main

import (
	"errors"
	"github.com/spf13/viper"
	"regexp"
	"strconv"
)

type ServerVariables struct {
	RoomName string `json:"name"`
	Port     string `json:"port"`
	RoomSize uint   `json:"size"`
	Secret   string `json:"secret"`
	Delay    uint   `json:"delay"`
}

func portValidate(args string) error {
	p, err := strconv.Atoi(args)
	if err != nil {
		return err
	}
	if p < 10 || p > 65535 {
		return errors.New("port out of range, must be between 1000 and 65535")
	}
	return nil
}

func RoomSizeValidate(args uint) error {
	if args < 1 {
		return errors.New("room size must be greater than 0")
	}
	return nil
}

func RoomNameValidate(args string) error {
	r, err := regexp.Compile("^[A-Za-z0-9_,]+")
	if err != nil {
		return err
	}
	if !r.MatchString(args) {
		return errors.New("room name can only contain alphanumeric characters and underscore")
	}
	if len(args) > 15 {
		return errors.New("room name too long")
	}
	return nil
}

func MessageDelayValidate(args uint) error {
	if args < 0 {
		return errors.New("delay must be greater than or equal to 0")
	}
	return nil
}

func GetServerVariables() (*ServerVariables, error) {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	port, ok := viper.Get("PORT").(string)
	if !ok {
		return nil, errors.New("port must be of type string")
	}
	if err := portValidate(port); err != nil {
		return nil, err
	}
	roomName, ok := viper.Get("ROOMNAME").(string)
	if !ok {
		return nil, errors.New("room name must be of type string")
	}
	if err := RoomNameValidate(roomName); err != nil {
		return nil, err
	}
	temp, ok := viper.Get("ROOMSIZE").(string)
	if !ok {
		return nil, errors.New("room size must be a whole number")
	}
	roomSize, err := strconv.Atoi(temp)
	if err != nil {
		return nil, err
	}
	if err := RoomSizeValidate(uint(roomSize)); err != nil {
		return nil, err
	}
	temp, ok = viper.Get("DELAY").(string)
	if !ok {
		return nil, errors.New("delay must be a whole number")
	}
	delay, err := strconv.Atoi(temp)
	if err != nil {
		return nil, err
	}
	if err := MessageDelayValidate(uint(delay)); err != nil {
		return nil, err
	}
	return &ServerVariables{
		RoomName: roomName,
		RoomSize: uint(roomSize),
		Port:     port,
		Delay:    uint(delay),
	}, nil
}
