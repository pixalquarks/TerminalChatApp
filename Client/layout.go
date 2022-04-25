package main

import (
	"errors"
	"fmt"
	"github.com/awesome-gocui/gocui"
	"log"
)

type Colors struct {
	msgBoxFG    gocui.Attribute
	msgBoxBG    gocui.Attribute
	msgBoxFGLow gocui.Attribute
}

var msgColor = Colors{
	msgBoxBG:    gocui.ColorGreen,
	msgBoxFG:    gocui.ColorWhite,
	msgBoxFGLow: gocui.ColorRed,
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView(viewArr[0], 0, 0, maxX/5-1, (maxY*4/5)-1, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = "Room : " + client.roomName
		//v.Editable = true
		v.Wrap = true

	}

	if v, err := g.SetView(viewArr[1], maxX/5-1, 0, maxX-1, (maxY*4/5)-1, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = viewArr[1]
		v.Autoscroll = true
	}
	if v, err := g.SetView(viewArr[2], 0, (maxY*4/5)-1, maxX/5-1, maxY-1, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = viewArr[2]
		v.Wrap = true
		v.Autoscroll = true
		if _, err := fmt.Fprint(v, "Helps"); err != nil {
			log.Println(err)
		}
	}
	if v, err := g.SetView(viewArr[3], maxX/5-1, (maxY*4/5)-1, maxX-70, maxY-1, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = client.clientName
		v.Editable = true
		v.BgColor = msgColor.msgBoxBG
		v.FgColor = msgColor.msgBoxFG
		v.EditDelete(true)

		if _, err = setCurrentViewOnTop(g, viewArr[3]); err != nil {
			return err
		}
	}
	if v, err := g.SetView(viewArr[4], maxX-70, (maxY*4/5)-1, maxX-1, maxY-1, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = viewArr[4]
		v.Autoscroll = true
		v.Wrap = true
	}
	return nil
}
