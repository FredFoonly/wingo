package main

import (
	"fmt"

	"github.com/BurntSushi/xgb/xproto"

	"github.com/BurntSushi/xgbutil/xwindow"

	"github.com/BurntSushi/wingo/focus"
	"github.com/BurntSushi/wingo/heads"
	"github.com/BurntSushi/wingo/workspace"
)

type clients []*client

func (cls clients) Get(i int) heads.Client {
	return cls[i]
}

func (cls clients) Len() int {
	return len(cls)
}

type wingoState struct {
	root    *xwindow.Window
	clients clients
	heads   *heads.Heads
	prompts prompts
	conf    *config
	theme   *theme
}

func newWingoState() *wingoState {
	wingo := &wingoState{
		clients: make(clients, 0, 50),
		heads:   nil,
	}
	return wingo
}

func (wingo *wingoState) initializeHeads() {
	wingo.heads = heads.NewHeads(X)
	for _, wrkName := range wingo.conf.workspaces {
		wingo.addWorkspace(wrkName)
	}
	wingo.heads.Initialize(wingo.clients)
}

func (wingo *wingoState) addClient(c *client) {
	if cliIndex(c, wingo.clients) != -1 {
		panic("BUG: Cannot add client that is already managed.")
	}
	wingo.clients = append(wingo.clients, c)
}

func (wingo *wingoState) removeClient(c *client) {
	if i := cliIndex(c, wingo.clients); i > -1 {
		wingo.clients = append(wingo.clients[:i], wingo.clients[i+1:]...)
	}
}

func (wingo *wingoState) findManagedClient(id xproto.Window) *client {
	for _, client := range wingo.clients {
		if client.Id() == id {
			return client
		}
	}
	return nil
}

func (wingo *wingoState) focusFallback() {
	wrk := wingo.workspace()
	for i := len(focus.Clients) - 1; i >= 0; i-- {
		switch client := focus.Clients[i].(type) {
		case *client:
			if client.frame.IsMapped() && client.workspace == wrk {
				focus.Focus(client)
				return
			}
		default:
			fmt.Printf("Unsupported client type: %T", client)
			panic("Not implemented.")
		}
	}
	focus.Root()
}

func (wingo *wingoState) workspace() *workspace.Workspace {
	return wingo.heads.ActiveWorkspace()
}

func (wingo *wingoState) addWorkspace(name string) {
	wrk := wingo.heads.NewWorkspace(name)
	wrk.PromptSlctGroup = wingo.prompts.slct.AddGroup(wrk)
	wrk.PromptSlctItem = wingo.prompts.slct.AddChoice(wrk)

	wingo.heads.AddWorkspace(wrk)
}
