package ui

import (
	"fmt"
	"net/url"
	"path"
	"time"

	"github.com/godot-go/godot-go/pkg/gdnative"
	"github.com/godot-go/godot-go/pkg/log"
	"github.com/jasmaa/hikawa/pkg/gemini"
	"github.com/jasmaa/hikawa/pkg/gemtext"
)

type Main struct {
	gdnative.NodeImpl
	gdnative.UserDataIdentifiableImpl

	client     gemini.Client
	currentUrl string
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
	searchBar := gdnative.NewLineEditWithOwner(p.GetNode(gdnative.NewNodePath("SearchBar")).GetOwnerObject())
	newUrl := p.navigatePage(searchBar.GetText())
	p.currentUrl = newUrl
	searchBar.SetText(newUrl)
	searchBar.ReleaseFocus()
}

func (p *Main) OnSearchButtonPressed() {
	searchBar := gdnative.NewLineEditWithOwner(p.GetNode(gdnative.NewNodePath("SearchBar")).GetOwnerObject())
	newUrl := p.navigatePage(searchBar.GetText())
	p.currentUrl = newUrl
	searchBar.SetText(newUrl)
}

func (p *Main) OnContentMetaClicked(meta string) {
	searchBar := gdnative.NewLineEditWithOwner(p.GetNode(gdnative.NewNodePath("SearchBar")).GetOwnerObject())

	target, err := url.Parse(meta)
	if err != nil {
		return
	}
	var targetUrl string
	if len(target.Scheme) == 0 {
		u, _ := url.Parse(p.currentUrl)
		u.Path = path.Join(u.Path, meta)
		targetUrl = u.String()
	} else {
		switch target.Scheme {
		case "gemini":
			targetUrl = meta
		default:
			return
		}
	}

	newUrl := p.navigatePage(targetUrl)
	p.currentUrl = newUrl
	searchBar.SetText(newUrl)
}

func (p *Main) OnClassRegistered(e gdnative.ClassRegisteredEvent) {
	// methods
	e.RegisterMethod("_ready", "Ready")
	e.RegisterMethod("_on_SearchButton_pressed", "OnSearchButtonPressed")
	e.RegisterMethod("_on_SearchBar_text_entered", "OnSearchBarTextEntered")
	e.RegisterMethod("_on_Content_meta_clicked", "OnContentMetaClicked")
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

func (p *Main) navigatePage(url string) string {
	content := gdnative.NewRichTextLabelWithOwner(p.GetNode(gdnative.NewNodePath("Content")).GetOwnerObject())
	defer content.ScrollToLine(0)

	clientResp, err := p.client.NavigatePage(url)
	if err != nil {
		content.SetBbcode(err.Error())
		return url
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
		return clientResp.Url
	}
}
