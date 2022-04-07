package main

import (
	"errors"
	"fmt"
	"github.com/awesome-gocui/gocui"
	"log"
	"time"
)

var (
	client  *clientHandle
	g       *gocui.Gui
	viewArr = [5]string{"Users", "ChatBox", "Commands", "MessageBox", "Logs"}
	active  = 0
	delay   = 0
)

const messageDelay = 3 // delay between user can input(in seconds)

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}

func nextView(g *gocui.Gui, _ *gocui.View) error {
	nextIndex := (active + 1) % len(viewArr)
	name := viewArr[nextIndex]

	if _, err := setCurrentViewOnTop(g, name); err != nil {
		return err
	}

	if nextIndex == 3 {
		g.Cursor = true
	} else {
		g.Cursor = false
	}

	active = nextIndex
	return nil
}

func OutMessage(g *gocui.Gui, v *gocui.View) error {
	if delay > 0 {
		return nil
	}
	go setChatTimeout()
	out, err := g.View(viewArr[1])
	if err != nil {
		log.Panicln(err)
	}
	msg := v.Buffer()
	if client != nil {
		if err := client.sendMessage(msg); err != nil {
			log.Println(err)
		}
	}
	msg = "--> " + msg
	if _, err := fmt.Fprintln(out, msg); err != nil {
		log.Panicln(err)
	}
	v.Clear()
	return nil
}

func LogOut(g *gocui.Gui, msg string) error {
	out, err := g.View(viewArr[4])
	if err != nil {
		return err
	}
	if _, err := fmt.Fprint(out, msg); err != nil {
		return err
	}
	return nil
}

func quit(_ *gocui.Gui, _ *gocui.View) error {
	return gocui.ErrQuit
}

func main() {
	var err error
	conn, err := GetConnection()
	if err != nil {
		log.Println(err)
	}

	client, err = CreateClient(conn)
	if err != nil {
		log.Println(err)
	}

	g, err = gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		log.Panicln(err)
	}

	g.Highlight = true
	g.Cursor = true
	g.SelFgColor = gocui.ColorGreen

	g.SetManagerFunc(layout)

	if err != nil {
		errMsg := fmt.Sprintf("Failed to connect to gRPC server :: %v", err)
		if err := LogOut(g, errMsg); err != nil {
			log.Panicln(err)
		}
	}

	defer func() {
		if err := LogOut(g, "closing client"); err != nil {
			log.Panicln(err)
		}
		if err := conn.Close(); err != nil {
			log.Printf("Could not close the connection :: %v", err)
		}
		g.Close()
	}()

	go client.receiveMessage(printToConsole)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding(viewArr[3], gocui.KeyEnter, gocui.ModNone, OutMessage); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		log.Panicln(err)
	}
}

// Outputs the message to the console ui chat-box window
func printToConsole(msg string) {
	out, err := g.View(viewArr[1])
	if err != nil {
		if err := LogOut(g, err.Error()); err != nil {
			log.Println(err)
		}
	}
	if msg[len(msg)-1] != '\n' {
		msg = msg + "\n"
	}
	out.WriteString(msg)
	g.UpdateAsync(func(gui *gocui.Gui) error {
		return nil
	})
	return
}

func setChatTimeout() {
	delay = messageDelay
	for {
		time.Sleep(time.Second)
		if delay <= 0 {
			return
		}
		delay--
	}
}
