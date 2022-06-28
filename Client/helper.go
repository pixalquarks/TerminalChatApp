package main

import (
	"errors"
	"github.com/fatih/color"
	"regexp"
	"strings"
)

func VerifyNameCharacters(name string) (bool, error) {
	r, err := regexp.Compile("^[A-Za-z0-9_]+$")
	if err != nil {
		return false, errors.New("error while compiling regexp")
	}
	return r.MatchString(name), nil
}

func formatMessage(sender string, msg string, timeStamp int64, maxWidth int) string {
	var sb strings.Builder
	serv := "SERVER"
	sep := " :: "
	endWhiteSpace := "    "
	restSpace := maxWidth - 10 - len(serv) - len(sep)
	if sender == "server" {
		l := len(msg)
		hl := 0
		if l < restSpace {
			hl = (restSpace - l) / 2
		}
		sb.WriteString(endWhiteSpace)
		temp := hl
		for temp > 0 {
			sb.WriteString("*")
			temp--
		}
		sb.WriteString(serv)
		sb.WriteString(sep)
		sb.WriteString(msg)
		temp = hl
		for temp > 0 {
			sb.WriteString("*")
			temp--
		}
		sb.WriteString(endWhiteSpace)
		sb.WriteString("\n")
		return redF.Sprintf(sb.String())
	}
	sb.WriteString("  ")
	sb.WriteString(blue.Add(color.Italic).Sprintf(sender))
	sb.WriteString("-->")
	sb.WriteString(msg)
	sb.WriteString("\n")
	return sb.String()
} //FIXME:random string errors refer to : C:\Users\kashy\Pictures\Screenshots\Screenshot (10) for more info.
