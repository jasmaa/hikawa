package main

import (
	"time"

	"github.com/jasmaa/hikawa/pkg/gemini"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func main() {

	var inTE, outTE *walk.TextEdit

	MainWindow{
		Title:   "Hikawa",
		MinSize: Size{600, 400},
		Layout:  VBox{},
		Children: []Widget{
			HSplitter{
				Children: []Widget{
					PushButton{
						Text: "<",
					},
					PushButton{
						Text: ">",
					},
					TextEdit{
						AssignTo:      &inTE,
						Text:          "gemini://gemini.circumlunar.space",
						StretchFactor: 20,
					},
					PushButton{
						OnClicked: func() {
							r, err := gemini.ParseRequest(inTE.Text())
							if err != nil {
								outTE.SetText(err.Error())
								return
							}

							textChan := make(chan string, 1)
							go func() {
								resp, err := r.Send()
								if err != nil {
									outTE.SetText(err.Error())
									return
								}
								if resp.Header.Status/10 == 2 {
									textChan <- resp.Body
								}
								textChan <- ""
							}()

							select {
							case res := <-textChan:
								outTE.SetText(res)
							case <-time.After(3 * time.Second):
								outTE.SetText("request timed out")
							}
						},
						Text: "Go!",
					},
				},
			},
			TextEdit{
				AssignTo:      &outTE,
				ReadOnly:      true,
				StretchFactor: 30,
			},
		},
	}.Run()
}
