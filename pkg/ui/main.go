package ui

import (
	"fmt"
	"time"

	"github.com/godot-go/godot-go/pkg/gdnative"
	"github.com/godot-go/godot-go/pkg/log"
	"github.com/jasmaa/hikawa/pkg/gemini"
	"github.com/jasmaa/hikawa/pkg/gemtext"
)

type Main struct {
	gdnative.NodeImpl
	gdnative.UserDataIdentifiableImpl

	client gemini.Client
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
	p.client = gemini.Client{
		Timeout: 3 * time.Second,
	}
	log.Info("Browser ready")
}

func (p *Main) OnSearchBarTextEntered(newText string) {
	p.navigatePage()
}

func (p *Main) OnSearchButtonPressed() {
	p.navigatePage()
}

func (p *Main) OnClassRegistered(e gdnative.ClassRegisteredEvent) {
	// methods
	e.RegisterMethod("_ready", "Ready")
	e.RegisterMethod("_on_SearchButton_pressed", "OnSearchButtonPressed")
	e.RegisterMethod("_on_SearchBar_text_entered", "OnSearchBarTextEntered")
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

func (p *Main) navigatePage() {
	searchBar := gdnative.NewLineEditWithOwner(p.GetNode(gdnative.NewNodePath("SearchBar")).GetOwnerObject())
	content := gdnative.NewRichTextLabelWithOwner(p.GetNode(gdnative.NewNodePath("Content")).GetOwnerObject())

	clientResp, err := p.client.NavigatePage(searchBar.GetText())
	if err != nil {
		content.SetBbcode(err.Error())
	} else {
		if clientResp.Response.Header.Status == gemini.STATUS_SUCCESS {
			if clientResp.Response.Header.Meta == "text/gemini" {
				contentBbcode := gemtext.ConvertToBbcode(clientResp.Response.Body)
				content.SetBbcode(contentBbcode)
			} else {
				content.SetBbcode(fmt.Sprintf("cannot display MIME type: %s", clientResp.Response.Header.Meta))
			}
		} else {
			content.SetBbcode(fmt.Sprintf("[%d] %s", clientResp.Response.Header.Status, clientResp.Response.Header.Meta))
		}
		searchBar.SetText(clientResp.Url)
	}
}
