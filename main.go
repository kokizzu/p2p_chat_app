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
	fmt.Fprint(out, red(bold(name+" > ")), message)

	return nil
}

func nextView(g *gocui.Gui, v *gocui.View) error {
	sending_message_content := v.ViewBuffer()
	if sending_message_content == "" {
		return nil
	}
	if err := add_to_mesasge_box("Hello", sending_message_content); err != nil {
		return err
	}

	v.Clear()
	v.SetCursor(0, 0)
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
		fmt.Fprint(v, yellow(bold("Help: ")), "Use ", yellow("/u <name>"), " to change name, ", yellow("/p <port>"), " to change port, ", yellow("/h <ip:port>"), " to connect to specific ip and ", yellow("CTRL+C"), " to quit!")
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func networkStarter() {
	ip := ""
	if !helper.Local {
		ip = peer.GetHostIPAddress()
	}

	go peer.Start(fmt.Sprintf("%s:%d", ip, helper.Port))
}

func debug_viewer() {
	if helper.Debug {
		for {
			debug_string := <-helper.DebugMessage
			add_to_mesasge_box("debug", debug_string)
		}
	}
}

func main() {
	var err error
	flag.StringVar(&helper.UserName, "name", helper.GetOsHostName(), "Choose Name, Defaults to os host name.")
	flag.IntVar(&helper.Port, "port", 8080, "Port")
	flag.BoolVar(&helper.Local, "local", false, "use localaddress?")
	flag.BoolVar(&helper.Debug, "debug", true, "show debug logs?")

	flag.Parse()
	g, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	go debug_viewer()
	go networkStarter()

	g.Highlight = true
	g.Cursor = true
	g.SelFgColor = gocui.ColorGreen

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("textbox", gocui.KeyEnter, gocui.ModNone, nextView); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
