package ui

import (
	"fmt"
	"net/url"
	"time"

	"github.com/godot-go/godot-go/pkg/gdnative"
	"github.com/godot-go/godot-go/pkg/log"
	"github.com/jasmaa/hikawa/pkg/browsing"
	"github.com/jasmaa/hikawa/pkg/gemini"
	"github.com/jasmaa/hikawa/pkg/gemtext"
)

const VERSION = "0.1.0"

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
	versionLabel := gdnative.NewLabelWithOwner(p.GetNode(gdnative.NewNodePath("StatusPanel/VersionLabel")).GetOwnerObject())
	versionLabel.SetText(fmt.Sprintf("v%s", VERSION))
	log.Info("Browser ready")
}

func (p *Main) OnSearchBarTextEntered(newText string) {
	searchBar := gdnative.NewLineEditWithOwner(p.GetNode(gdnative.NewNodePath("NavPanel/SearchBar")).GetOwnerObject())
	p.loadPage(searchBar.GetText(), true)
}

func (p *Main) OnSearchButtonPressed() {
	searchBar := gdnative.NewLineEditWithOwner(p.GetNode(gdnative.NewNodePath("NavPanel/SearchBar")).GetOwnerObject())
	p.loadPage(searchBar.GetText(), true)
}

func (p *Main) OnContentMetaClicked(meta string) {
	currentUrl, err := p.history.GetCurrentUrl()
	if err != nil {
		return
	}
	targetUrl, err := gemini.NextUrl(currentUrl, meta)
	if err != nil {
		return
	}
	p.loadPage(targetUrl, true)
}

func (p *Main) OnBackButtonPressed() {
	err := p.history.GoBack()
	if err != nil {
		return
	}
	currentUrl, _ := p.history.GetCurrentUrl()
	p.loadPage(currentUrl, false)
}

func (p *Main) OnForwardButtonPressed() {
	err := p.history.GoForward()
	if err != nil {
		return
	}
	currentUrl, _ := p.history.GetCurrentUrl()
	p.loadPage(currentUrl, false)
}

func (p *Main) OnSubmitButtonPressed() {
	searchBar := gdnative.NewLineEditWithOwner(p.GetNode(gdnative.NewNodePath("NavPanel/SearchBar")).GetOwnerObject())
	inputPopup := gdnative.NewPopupWithOwner(p.GetNode(gdnative.NewNodePath("InputPopup")).GetOwnerObject())
	promptBar := gdnative.NewLineEditWithOwner(p.GetNode(gdnative.NewNodePath("InputPopup/PromptBar")).GetOwnerObject())

	inputPopup.Hide()
	u, _ := url.Parse(searchBar.GetText())
	u.RawQuery = promptBar.GetText()
	p.loadPage(u.String(), true)
}

func (p *Main) OnPromptBarTextEntered(newText string) {
	searchBar := gdnative.NewLineEditWithOwner(p.GetNode(gdnative.NewNodePath("NavPanel/SearchBar")).GetOwnerObject())
	inputPopup := gdnative.NewPopupWithOwner(p.GetNode(gdnative.NewNodePath("InputPopup")).GetOwnerObject())
	promptBar := gdnative.NewLineEditWithOwner(p.GetNode(gdnative.NewNodePath("InputPopup/PromptBar")).GetOwnerObject())

	inputPopup.Hide()
	u, _ := url.Parse(searchBar.GetText())
	u.RawQuery = promptBar.GetText()
	p.loadPage(u.String(), true)
}

func (p *Main) OnClassRegistered(e gdnative.ClassRegisteredEvent) {
	// methods
	e.RegisterMethod("_ready", "Ready")
	e.RegisterMethod("_on_SearchButton_pressed", "OnSearchButtonPressed")
	e.RegisterMethod("_on_SearchBar_text_entered", "OnSearchBarTextEntered")
	e.RegisterMethod("_on_Content_meta_clicked", "OnContentMetaClicked")
	e.RegisterMethod("_on_BackButton_pressed", "OnBackButtonPressed")
	e.RegisterMethod("_on_ForwardButton_pressed", "OnForwardButtonPressed")
	e.RegisterMethod("_on_SubmitButton_pressed", "OnSubmitButtonPressed")
	e.RegisterMethod("_on_PromptBar_text_entered", "OnPromptBarTextEntered")
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

func (p *Main) navigatePage(url string, shouldPushHistory bool) string {
	content := gdnative.NewRichTextLabelWithOwner(p.GetNode(gdnative.NewNodePath("Content")).GetOwnerObject())
	defer content.ScrollToLine(0)

	clientResp, err := p.client.NavigatePage(url)

	if err != nil {
		content.SetBbcode(err.Error())
		if shouldPushHistory {
			p.history.Push(url)
		}
		return url
	}

	if clientResp.Response.Header.Status == gemini.STATUS_SUCCESS {
		if _, ok := clientResp.MimeTypes["text/gemini"]; ok {
			contentBbcode := gemtext.ConvertToBbcode(clientResp.Response.Body)
			content.SetBbcode(contentBbcode)
		} else {
			content.SetBbcode(fmt.Sprintf("cannot display MIME type: %s", clientResp.Response.Header.Meta))
		}
	} else if clientResp.Response.Header.Status == gemini.STATUS_INPUT {
		inputPopup := gdnative.NewPopupWithOwner(p.GetNode(gdnative.NewNodePath("InputPopup")).GetOwnerObject())
		question := gdnative.NewLabelWithOwner(p.GetNode(gdnative.NewNodePath("InputPopup/Question")).GetOwnerObject())
		promptBar := gdnative.NewLineEditWithOwner(p.GetNode(gdnative.NewNodePath("InputPopup/PromptBar")).GetOwnerObject())
		question.SetText(clientResp.Response.Header.Meta)
		promptBar.Clear()
		inputPopup.PopupCentered(gdnative.NewVector2(400, 200))
	} else {
		content.SetBbcode(fmt.Sprintf("[%d] %s", clientResp.Response.Header.Status, clientResp.Response.Header.Meta))
	}

	if clientResp.Response.Header.Status != gemini.STATUS_INPUT &&
		clientResp.Response.Header.Status != gemini.STATUS_SENSITIVE_INPUT &&
		shouldPushHistory {
		p.history.Push(clientResp.Url)
	}

	return clientResp.Url
}

func (p *Main) loadPage(targetUrl string, shouldPushHistory bool) {
	statusLabel := gdnative.NewLabelWithOwner(p.GetNode(gdnative.NewNodePath("StatusPanel/StatusLabel")).GetOwnerObject())
	searchBar := gdnative.NewLineEditWithOwner(p.GetNode(gdnative.NewNodePath("NavPanel/SearchBar")).GetOwnerObject())

	statusLabel.SetText("Loading...")
	p.setEnableNavPanel(false)

	go func() {
		newUrl := p.navigatePage(targetUrl, shouldPushHistory)
		searchBar.SetText(newUrl)
		searchBar.ReleaseFocus()
		p.setEnableNavPanel(true)
		p.setNavigationButtons()
		statusLabel.SetText("Done.")
	}()
}

func (p *Main) setNavigationButtons() {
	backButton := gdnative.NewButtonWithOwner(p.GetNode(gdnative.NewNodePath("NavPanel/BackButton")).GetOwnerObject())
	backButton.SetDisabled(!p.history.CanGoBack())
	forwardButton := gdnative.NewButtonWithOwner(p.GetNode(gdnative.NewNodePath("NavPanel/ForwardButton")).GetOwnerObject())
	forwardButton.SetDisabled(!p.history.CanGoForward())
}

func (p *Main) setEnableNavPanel(enabled bool) {
	searchBar := gdnative.NewLineEditWithOwner(p.GetNode(gdnative.NewNodePath("NavPanel/SearchBar")).GetOwnerObject())
	searchButton := gdnative.NewButtonWithOwner(p.GetNode(gdnative.NewNodePath("NavPanel/SearchButton")).GetOwnerObject())
	backButton := gdnative.NewButtonWithOwner(p.GetNode(gdnative.NewNodePath("NavPanel/BackButton")).GetOwnerObject())
	forwardButton := gdnative.NewButtonWithOwner(p.GetNode(gdnative.NewNodePath("NavPanel/ForwardButton")).GetOwnerObject())
	searchBar.SetEditable(enabled)
	searchButton.SetDisabled(!enabled)
	backButton.SetDisabled(!enabled)
	forwardButton.SetDisabled(!enabled)
}
