package main

import (
	"fmt"

	g "github.com/AllenDang/giu"
	"github.com/jasmaa/hikawa/pkg/ui"
)

const VERSION = "2.0.0"

func main() {
	wnd := g.NewMasterWindow(fmt.Sprintf("Hikawa - v%s", VERSION), 800, 600, 0)
	wnd.Run(ui.Loop)
}
