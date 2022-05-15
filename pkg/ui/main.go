package ui

import (
	"fmt"
	"net/url"
	"time"

	g "github.com/AllenDang/giu"
	"github.com/jasmaa/hikawa/pkg/browsing"
	"github.com/jasmaa/hikawa/pkg/gemini"
	"github.com/jasmaa/hikawa/pkg/gemtext"
)

var (
	searchText              string
	inputText               string
	content                 string
	isBackButtonDisabled    bool
	isForwardButtonDisabled bool
	isSearchButtonDisabled  bool
	isInputMode             bool
	client                  gemini.Client
	history                 browsing.History
)

func init() {
	searchText = "gemini://gemini.circumlunar.space/"
	isBackButtonDisabled = true
	isForwardButtonDisabled = true
	isSearchButtonDisabled = false
	isInputMode = false
	client = gemini.Client{
		Timeout: 3 * time.Second,
	}
}

func onSubmitSearch() {
	setLoading()
	go func() {
		newUrl := navigatePage(searchText, true)
		searchText = newUrl
		setNavigationButtons()
	}()
}

func onContentMetaClicked(meta string) {
	currentUrl, err := history.GetCurrentUrl()
	if err != nil {
		return
	}
	targetUrl, err := gemini.NextUrl(currentUrl, meta)
	if err != nil {
		return
	}

	parsedTargetUrl, err := url.Parse(targetUrl)
	if err != nil {
		return
	}

	targetScheme := parsedTargetUrl.Scheme
	if targetScheme == "http" || targetScheme == "https" {
		g.OpenURL(targetUrl)
	} else {
		setLoading()
		go func() {
			newUrl := navigatePage(targetUrl, true)
			searchText = newUrl
			setNavigationButtons()
			g.Update()
		}()
	}
}

func onBackButtonPressed() {
	err := history.GoBack()
	if err != nil {
		return
	}

	setLoading()
	go func() {
		currentUrl, _ := history.GetCurrentUrl()
		newUrl := navigatePage(currentUrl, false)
		searchText = newUrl
		setNavigationButtons()
		g.Update()
	}()
}

func onForwardButtonPressed() {
	err := history.GoForward()
	if err != nil {
		return
	}

	setLoading()
	go func() {
		currentUrl, _ := history.GetCurrentUrl()
		newUrl := navigatePage(currentUrl, false)
		searchText = newUrl
		setNavigationButtons()
		g.Update()
	}()
}

func onSubmitInput() {
	u, _ := url.Parse(searchText)
	u.RawQuery = inputText

	setLoading()
	go func() {
		newUrl := navigatePage(u.String(), true)
		searchText = newUrl
		inputText = ""
		setNavigationButtons()
	}()
}

func navigatePage(rawurl string, shouldPushHistory bool) string {
	isInputMode = false
	clientResp, err := client.NavigatePage(rawurl)

	if err != nil {
		content = err.Error()
		if shouldPushHistory {
			history.Push(rawurl)
		}
		return rawurl
	}

	if clientResp.Response.Header.Status == gemini.STATUS_SUCCESS {
		if _, ok := clientResp.MimeTypes["text/gemini"]; ok {
			content = gemtext.ConvertToMarkdown(clientResp.Response.Body)
		} else {
			content = fmt.Sprintf("cannot display MIME type: %s", clientResp.Response.Header.Meta)
		}
	} else if clientResp.Response.Header.Status == gemini.STATUS_INPUT {
		isInputMode = true
	} else {
		content = fmt.Sprintf("[%d] %s", clientResp.Response.Header.Status, clientResp.Response.Header.Meta)
	}

	if shouldPushHistory {
		history.Push(clientResp.Url)
	}

	return clientResp.Url
}

func setNavigationButtons() {
	isBackButtonDisabled = !history.CanGoBack()
	isForwardButtonDisabled = !history.CanGoForward()
	isSearchButtonDisabled = false
}

func setLoading() {
	content = "Loading..."
	isBackButtonDisabled = true
	isForwardButtonDisabled = true
	isSearchButtonDisabled = true
}

func Loop() {
	var contentWidget g.Widget
	if isInputMode {
		contentWidget = g.Row(
			g.InputText(&inputText),
			g.Event().OnKeyPressed(g.KeyEnter, onSubmitInput),
			g.Button("Submit").OnClick(onSubmitInput),
		)
	} else {
		contentWidget = g.Markdown(&content).OnLink(func(url string) {
			go onContentMetaClicked(url)
		})
	}

	g.SingleWindow().Layout(
		g.Table().Rows(
			g.TableRow(
				g.Row(
					g.Button("<").OnClick(onBackButtonPressed).Disabled(isBackButtonDisabled),
					g.Button(">").OnClick(onForwardButtonPressed).Disabled(isForwardButtonDisabled),
					g.InputText(&searchText),
					g.Event().OnKeyPressed(g.KeyEnter, onSubmitSearch),
					g.Button("Go").OnClick(onSubmitSearch).Disabled(isSearchButtonDisabled),
				),
			),
			g.TableRow(
				contentWidget,
			),
		).Freeze(1, 1),
	)
}
