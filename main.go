package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/gookit/color"
	"github.com/jroimartin/gocui"
	"github.com/sairash/p2p_chat_app/helper"
	"github.com/sairash/p2p_chat_app/peer"
)

var (
	red    = color.FgRed.Render
	bold   = color.Bold.Render
	blue   = color.FgBlue.Render
	yellow = color.FgYellow.Render
	green  = color.FgGreen.Render
	g      *gocui.Gui
)

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}

func add_to_mesasge_box(name string, message string) error {
	out, err := g.View("messages")
	if err != nil {
		return err

	}
	fmt.Fprint(out, name, message)

	return nil
}

func commandSection(command string) {

	switch command[1] {
	case 'h':

		// output := strings.TrimPrefix(command, "/h ")

		// peer.ConnnectHost(output)
		// chewang do;,a tamang is a gir; amd i and kme wow ownfocommand
		// connection to host
		return
	case 'l':
		reply := "List of known hosts. \n"

		for k := range helper.ConnectedHosts {
			reply += "- " + k + "\n"
		}

		helper.MessageChan <- helper.DisplayMessage{
			Message: helper.Message{
				Text: reply,
				Name: "commandSection",
			},
			TypeOfMessage: helper.ImportantDebug,
		}

		return
	case 'u':
		// username
		return
	default:
		// Error command message
		helper.MessageChan <- helper.DebugMessage("Command not available.", "commandSection")
	}
}

func sendMessage(g *gocui.Gui, v *gocui.View) error {
	defer v.Clear()
	defer v.SetCursor(0, 0)

	sending_message_content := v.ViewBuffer()
	if sending_message_content == "" {
		return nil
	}

	if sending_message_content[0] == '/' {
		commandSection(sending_message_content)
		return nil
	}
	peer.Send(sending_message_content)

	if err := add_to_mesasge_box(green(bold(helper.UserName+" > ")), sending_message_content); err != nil {
		return err
	}

	return nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("messages", 0, 0, maxX-1, maxY-5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Messages"
		v.Wrap = true
		v.Autoscroll = true
	}

	if v, err := g.SetView("textbox", 0, maxY-4, maxX-1, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = helper.UserName
		v.Editable = true

		if _, err = setCurrentViewOnTop(g, "textbox"); err != nil {
			return err
		}
	}

	if v, err := g.SetView("help", 0, maxY-2, maxX-1, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		fmt.Fprint(v, yellow(bold("Help: ")), "Use ", yellow("/l"), " to list the connected peers.", yellow("CTRL+C"), " to quit!")
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func networkStarter() {
	if !helper.Local {
		helper.IP = peer.GetHostIPAddressV4()
	}

	helper.IPPORT = fmt.Sprintf("%s:%d", helper.IP, helper.Port)

	go peer.Start(helper.IPPORT)
	go peer.StartDiscovery()
}

func message_value_addder() {
	for {
		string_message := <-helper.MessageChan
		if string_message.TypeOfMessage == helper.Peer {
			add_to_mesasge_box(blue(bold(string_message.Name+" > ")), string_message.Message.Text)
		} else if helper.Debug || string_message.TypeOfMessage == helper.ImportantDebug {
			add_to_mesasge_box(red(bold("Debug > ")), string_message.Message.Text+"\n")
		}

	}
}

func main() {
	var err error
	flag.StringVar(&helper.UserName, "name", helper.GetOsHostName(), "Choose Name, Defaults to os host name.")
	flag.IntVar(&helper.Port, "port", 8080, "Port")
	flag.BoolVar(&helper.Local, "local", false, "use local address?")
	flag.BoolVar(&helper.Debug, "debug", false, "show debug logs?")

	flag.Parse()
	g, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	go message_value_addder()
	go networkStarter()

	g.Highlight = true
	g.Cursor = true
	g.SelFgColor = gocui.ColorGreen

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("textbox", gocui.KeyEnter, gocui.ModNone, sendMessage); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
