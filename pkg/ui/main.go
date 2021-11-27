package ui

import (
	"fmt"
	"net/url"
	"path"
	"time"

	"github.com/godot-go/godot-go/pkg/gdnative"
	"github.com/godot-go/godot-go/pkg/log"
	"github.com/jasmaa/hikawa/pkg/browsing"
	"github.com/jasmaa/hikawa/pkg/gemini"
	"github.com/jasmaa/hikawa/pkg/gemtext"
)

type Main struct {
	gdnative.NodeImpl
	gdnative.UserDataIdentifiableImpl

	client  gemini.Client
	history browsing.History
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
	p.history = browsing.NewHistory()
	p.setNavigationButtons()
	log.Info("Browser ready")
}

func (p *Main) OnSearchBarTextEntered(newText string) {
	searchBar := gdnative.NewLineEditWithOwner(p.GetNode(gdnative.NewNodePath("SearchBar")).GetOwnerObject())
	newUrl := p.navigatePage(searchBar.GetText())
	p.history.Push(newUrl)
	searchBar.SetText(newUrl)
	searchBar.ReleaseFocus()
	p.setNavigationButtons()
}

func (p *Main) OnSearchButtonPressed() {
	searchBar := gdnative.NewLineEditWithOwner(p.GetNode(gdnative.NewNodePath("SearchBar")).GetOwnerObject())
	newUrl := p.navigatePage(searchBar.GetText())
	p.history.Push(newUrl)
	searchBar.SetText(newUrl)
	p.setNavigationButtons()
}

func (p *Main) OnContentMetaClicked(meta string) {
	searchBar := gdnative.NewLineEditWithOwner(p.GetNode(gdnative.NewNodePath("SearchBar")).GetOwnerObject())

	target, err := url.Parse(meta)
	if err != nil {
		return
	}
	var targetUrl string
	if len(target.Scheme) == 0 {
		currentUrl, err := p.history.GetCurrentUrl()
		if err != nil {
			return
		}
		u, _ := url.Parse(currentUrl)
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
	p.history.Push(newUrl)
	searchBar.SetText(newUrl)
	p.setNavigationButtons()
}

func (p *Main) OnBackButtonPressed() {
	searchBar := gdnative.NewLineEditWithOwner(p.GetNode(gdnative.NewNodePath("SearchBar")).GetOwnerObject())

	err := p.history.GoBack()
	if err != nil {
		return
	}
	currentUrl, _ := p.history.GetCurrentUrl()
	newUrl := p.navigatePage(currentUrl)
	searchBar.SetText(newUrl)
	p.setNavigationButtons()

}

func (p *Main) OnForwardButtonPressed() {
	searchBar := gdnative.NewLineEditWithOwner(p.GetNode(gdnative.NewNodePath("SearchBar")).GetOwnerObject())

	err := p.history.GoForward()
	if err != nil {
		return
	}
	currentUrl, _ := p.history.GetCurrentUrl()
	newUrl := p.navigatePage(currentUrl)
	searchBar.SetText(newUrl)
	p.setNavigationButtons()
}

func (p *Main) OnClassRegistered(e gdnative.ClassRegisteredEvent) {
	// methods
	e.RegisterMethod("_ready", "Ready")
	e.RegisterMethod("_on_SearchButton_pressed", "OnSearchButtonPressed")
	e.RegisterMethod("_on_SearchBar_text_entered", "OnSearchBarTextEntered")
	e.RegisterMethod("_on_Content_meta_clicked", "OnContentMetaClicked")
	e.RegisterMethod("_on_BackButton_pressed", "OnBackButtonPressed")
	e.RegisterMethod("_on_ForwardButton_pressed", "OnForwardButtonPressed")
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
			if _, ok := clientResp.MimeTypes["text/gemini"]; ok {
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

func (p *Main) setNavigationButtons() {
	backButton := gdnative.NewButtonWithOwner(p.GetNode(gdnative.NewNodePath("BackButton")).GetOwnerObject())
	backButton.SetDisabled(!p.history.CanGoBack())
	forwardButton := gdnative.NewButtonWithOwner(p.GetNode(gdnative.NewNodePath("ForwardButton")).GetOwnerObject())
	forwardButton.SetDisabled(!p.history.CanGoForward())
}
