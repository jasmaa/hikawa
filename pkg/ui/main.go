package ui

import (
	"time"

	"github.com/godot-go/godot-go/pkg/gdnative"
	"github.com/godot-go/godot-go/pkg/log"
	"github.com/jasmaa/hikawa/pkg/gemini"
)

type Main struct {
	gdnative.NodeImpl
	gdnative.UserDataIdentifiableImpl
}

func (p *Main) ClassName() string {
	return "Main"
}

func (p *Main) BaseClass() string {
	return "Control"
}

func (p *Main) Init() {
}

func (p *Main) Ready() {
	log.Info("Browser ready")
}

func (p *Main) OnSearchButtonPressed() {
	searchBar := gdnative.NewLineEditWithOwner(p.GetNode(gdnative.NewNodePath("SearchBar")).GetOwnerObject())
	content := gdnative.NewRichTextLabelWithOwner(p.GetNode(gdnative.NewNodePath("Content")).GetOwnerObject())

	r, err := gemini.ParseRequest(searchBar.GetText())
	if err != nil {
		content.SetText(err.Error())
		return
	}

	textChan := make(chan string, 1)
	go func() {
		resp, err := r.Send()
		if err != nil {
			textChan <- err.Error()
			return
		}
		if resp.Header.Status/10 == 2 {
			textChan <- resp.Body
		}
		textChan <- ""
	}()

	select {
	case res := <-textChan:
		content.SetText(res)
	case <-time.After(3 * time.Second):
		content.SetText("request timed out")
	}
}

func (p *Main) OnClassRegistered(e gdnative.ClassRegisteredEvent) {
	// methods
	e.RegisterMethod("_ready", "Ready")
	e.RegisterMethod("_on_SearchButton_pressed", "OnSearchButtonPressed")
}

func NewMainWithOwner(owner *gdnative.GodotObject) Main {
	inst := gdnative.GetCustomClassInstanceWithOwner(owner).(*Main)
	return *inst
}

func init() {
	gdnative.RegisterInitCallback(func() {
		gdnative.RegisterClass(&Main{})
	})
}
