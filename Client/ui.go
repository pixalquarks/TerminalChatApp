package main

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/awesome-gocui/gocui"
	"github.com/fatih/color"
)

var (
	client  *clientHandle
	g       *gocui.Gui
	viewArr = [5]string{"Users", "ChatBox", "Commands", "MessageBox", "Logs"}
	active  = 0
	canSend = true
	greenF  = color.New(color.FgGreen)
	redF    = color.New(color.FgRed)
	blue    = color.New(color.FgBlue)
	magenta = color.New(color.FgMagenta)
)

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
	if !canSend {
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
	x, _ := out.Size()
	t := len(msg)
	var sb strings.Builder
	if t < x {
		dif := x - t
		for dif > 0 {
			sb.WriteString(" ")
			dif--
		}
		sb.WriteString(msg)
	}

	if _, err := fmt.Fprintln(out, sb.String()); err != nil {
		log.Panicln(err)
	}
	v.Clear()
	go setChatTimeout()
	return nil
}

func LogOut(msg string) error {
	out, err := g.View(viewArr[4])
	if err != nil {
		return err
	}
	if _, err := fmt.Fprint(out, msg); err != nil {
		return err
	}
	return nil
}

func UpdateUserList() {
	out, err := g.View(viewArr[0])
	if err != nil {
		if err := LogOut(err.Error()); err != nil {
			log.Println(err)
		}
	}
	users, err := client.GetUserNames()
	if err != nil {
		if err := LogOut(err.Error()); err != nil {
			log.Println(err)
		}
	}
	out.Clear()
	for _, user := range users {
		if user == client.clientName {
			continue
		}
		if _, err := magenta.Fprintln(out, user); err != nil {
			if err := LogOut(err.Error()); err != nil {
				log.Println(err)
			}
		}
	}
	if _, err := greenF.Fprintln(out, client.clientName); err != nil {
		if err := LogOut(err.Error()); err != nil {
			log.Println(err)
		}
	}
}

func quit(_ *gocui.Gui, _ *gocui.View) error {
	return gocui.ErrQuit
}

func main() {
	var err error
	conn, err := GetConnection()
	if err != nil {
		log.Panicln(err)
	}

	client, err = CreateClient(conn)
	if err != nil {
		log.Panicln(err)
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
		if err := LogOut(errMsg); err != nil {
			log.Panicln(err)
		}
	}

	defer func() {
		if err := LogOut("closing client"); err != nil {
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
func printToConsole(sender string, msg string, timeStamp int64) {
	out, err := g.View(viewArr[1])
	x, _ := out.Size()
	if err != nil {
		if err := LogOut(err.Error()); err != nil {
			log.Println(err)
		}
	}
	m := formatMessage(sender, msg, timeStamp, x)
	out.WriteString(m)
	g.UpdateAsync(func(gui *gocui.Gui) error {
		return nil
	})
	return
}

func setChatTimeout() {
	v, err := g.View(viewArr[3])
	if err != nil {
		if err := LogOut(err.Error()); err != nil {
			log.Panicln(err)
		}
	}
	v.FgColor = msgColor.msgBoxFGLow
	canSend = false
	time.Sleep(time.Second * time.Duration(client.delay))
	canSend = true
	msg := v.Buffer()
	v.FgColor = msgColor.msgBoxFG
	v.Clear()
	v.WriteString(msg)
	if err := v.SetCursor(len(msg), 0); err != nil {
		log.Panicln(err)
	}
	g.UpdateAsync(func(gui *gocui.Gui) error {
		return nil
	})
}
