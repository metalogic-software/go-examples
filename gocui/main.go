// Copyright 2016 The Metalogic Software Corporation Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"

	"github.com/jroimartin/gocui"
)

const (
	global        = ""
	vTop          = "top"
	vDestinations = "destinations"
	vMessage      = "message"
	vNote         = "note"
	vPlatforms    = "platforms"
	vSubmit       = "submit"
	vTranscripts  = "transcripts"
	vText         = "text"
)

func top(g *gocui.Gui, text string) {
	if v, err := g.View(vTop); err == nil {
		v.Clear()
		fmt.Fprintf(v, "%s", text)
	}
}

func nextView(g *gocui.Gui, v *gocui.View) error {

	if v == nil {
		return g.SetCurrentView(vPlatforms)
	}
	getLine(g, v)

	switch v.Name() {
	case vPlatforms:
		return g.SetCurrentView(vTranscripts)
	case vTranscripts:
		return g.SetCurrentView(vDestinations)
	case vDestinations:
		if x, err := g.View(vTop); err == nil {
			x.Clear()
			fmt.Fprintf(x, "Hit enter to submit selections")
		}
		return g.SetCurrentView(vSubmit)
	case vSubmit:
		return g.SetCurrentView(vText)
	case vText:
		return g.SetCurrentView(vPlatforms)
	}
	return g.SetCurrentView(vPlatforms)
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

func menuDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		ox, oy := v.Origin()

		if err := v.SetCursor(cx, cy+1); err != nil {
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func menuUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if cy == oy {
			if err := v.SetOrigin(ox, oy); err != nil {
				return err
			}
			return nil
		}
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

func submit(g *gocui.Gui, v *gocui.View) error {
	return message(g, fmt.Sprintf("Submitted %s", opt))
}

func message(g *gocui.Gui, text string) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView(vMessage, maxX/2-30, maxY/2, maxX/2+30, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, text)
		if err := g.SetCurrentView(vMessage); err != nil {
			return err
		}
	}
	return nil
}

func delMsg(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView(vMessage); err != nil {
		return err
	}
	top(g, "Select target platform, source transcript and destination.")

	if err := g.SetCurrentView(vPlatforms); err != nil {
		return err
	}
	return nil
}

func getLine(g *gocui.Gui, v *gocui.View) error {
	var line string
	var err error

	_, cy := v.Cursor()
	if line, err = v.Line(cy); err != nil {
		line = ""
	}
	switch v.Name() {
	case vPlatforms:
		opt.platform = line
	case vTranscripts:
		opt.transcript = line
	case vDestinations:
		opt.destination = line
	}
	return err
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding(vPlatforms, gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding(vPlatforms, gocui.KeyArrowRight, gocui.ModNone, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding(vTranscripts, gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding(vTranscripts, gocui.KeyArrowRight, gocui.ModNone, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding(vDestinations, gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding(vDestinations, gocui.KeyArrowRight, gocui.ModNone, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding(vSubmit, gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding(vText, gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding(vPlatforms, gocui.KeyArrowDown, gocui.ModNone, menuDown); err != nil {
		return err
	}
	if err := g.SetKeybinding(vPlatforms, gocui.KeyArrowUp, gocui.ModNone, menuUp); err != nil {
		return err
	}
	if err := g.SetKeybinding(vTranscripts, gocui.KeyArrowDown, gocui.ModNone, menuDown); err != nil {
		return err
	}
	if err := g.SetKeybinding(vTranscripts, gocui.KeyArrowUp, gocui.ModNone, menuUp); err != nil {
		return err
	}
	if err := g.SetKeybinding(vDestinations, gocui.KeyArrowDown, gocui.ModNone, menuDown); err != nil {
		return err
	}
	if err := g.SetKeybinding(vDestinations, gocui.KeyArrowUp, gocui.ModNone, menuUp); err != nil {
		return err
	}
	if err := g.SetKeybinding(vPlatforms, gocui.KeyEnter, gocui.ModNone, getLine); err != nil {
		return err
	}
	if err := g.SetKeybinding(vTranscripts, gocui.KeyEnter, gocui.ModNone, getLine); err != nil {
		return err
	}
	if err := g.SetKeybinding(vDestinations, gocui.KeyEnter, gocui.ModNone, getLine); err != nil {
		return err
	}
	if err := g.SetKeybinding(vSubmit, gocui.KeyEnter, gocui.ModNone, submit); err != nil {
		return err
	}
	if err := g.SetKeybinding(vMessage, gocui.KeyEnter, gocui.ModNone, delMsg); err != nil {
		return err
	}
	if err := g.SetKeybinding(vText, gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding(vText, gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding(global, gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding(global, gocui.KeyPgdn, gocui.ModNone, downText); err != nil {
		return err
	}
	if err := g.SetKeybinding(global, gocui.KeyPgup, gocui.ModNone, upText); err != nil {
		return err
	}
	return nil
}

func downText(g *gocui.Gui, v *gocui.View) error {
	if v, err := g.View(vText); err == nil {
		ox, oy := v.Origin()
		oy = oy + 10
		v.SetOrigin(ox, oy)
	}
	return nil
}

func upText(g *gocui.Gui, v *gocui.View) error {
	if v, err := g.View(vText); err == nil {
		ox, oy := v.Origin()
		oy = oy - 10
		v.SetOrigin(ox, oy)
	}
	return nil
}

func saveMain(g *gocui.Gui, v *gocui.View) error {
	f, err := ioutil.TempFile("", "gocui_demo_")
	if err != nil {
		return err
	}
	defer f.Close()

	p := make([]byte, 5)
	v.Rewind()
	for {
		n, err := v.Read(p)
		if n > 0 {
			if _, err := f.Write(p[:n]); err != nil {
				return err
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func saveVisualMain(g *gocui.Gui, v *gocui.View) error {
	f, err := ioutil.TempFile("", "gocui_demo_")
	if err != nil {
		return err
	}
	defer f.Close()

	vb := v.ViewBuffer()
	if _, err := io.Copy(f, strings.NewReader(vb)); err != nil {
		return err
	}
	return nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView(vTop, 1, 1, maxX-2, 3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		top(g, "Select target platform, source transcript and destination.")
	}

	if v, err := g.SetView(vPlatforms, 1, 4, 11, 8); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "  Hub  "
		v.Highlight = true

		fmt.Fprintln(v, " DEV ")
		fmt.Fprintln(v, " CERT ")
		fmt.Fprintln(v, " PROD ")
		if err := g.SetCurrentView(vPlatforms); err != nil {
			return err
		}
	}
	if v, err := g.SetView(vTranscripts, 12, 4, 36, 8); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "  Source Transcript  "
		v.Highlight = true
		fmt.Fprintln(v, " kpu01.xml ")
		fmt.Fprintln(v, " langara01.xml ")
		fmt.Fprintln(v, " sfu01.xml ")
		fmt.Fprintln(v, " ubc01.xml ")
		fmt.Fprintln(v, " ufv01.xml ")
	}

	if v, err := g.SetView(vDestinations, 38, 4, maxX-2, 8); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "  Destination Institution  "
		v.Highlight = true

		fmt.Fprintln(v, " Douglas College ")
		fmt.Fprintln(v, " Kwantlen Polytechnical University ")
		fmt.Fprintln(v, " Simon Fraser University ")
	}

	// submit button
	if v, err := g.SetView(vSubmit, 1, 8, 8, 10); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		v.Highlight = true
		v.BgColor = gocui.ColorYellow
		v.FgColor = gocui.ColorBlack
		v.SelBgColor = gocui.ColorWhite
		fmt.Fprintf(v, "Submit")
	}
	if v, err := g.SetView(vNote, 18, 8, maxX-1, 10); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		v.FgColor = gocui.ColorBlack
		v.SelBgColor = gocui.ColorWhite
		fmt.Fprintf(v, "Note: Scroll transcript with Page Up & Page Down")
	}
	// show contents of currently selected transcript
	if v, err := g.SetView(vText, 1, 10, maxX-2, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		b, err := ioutil.ReadFile("Mark.Twain-Tom.Sawyer.txt")
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(v, "%s", b)
		v.Wrap = true
	}

	return nil
}

type options struct {
	platform    string
	transcript  string
	destination string
}

func (p options) String() string {
	return fmt.Sprintf("%s, %s, %s", p.platform, p.transcript, p.destination)
}

// opt holds the currently selected options
var opt = new(options)

func main() {
	g := gocui.NewGui()
	if err := g.Init(); err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetLayout(layout)
	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}
	g.SelBgColor = gocui.ColorGreen
	g.SelFgColor = gocui.ColorBlack
	g.Cursor = true
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
