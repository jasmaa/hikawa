package main

import (
	"fmt"
	"image"
	"os"

	g "github.com/AllenDang/giu"
	"github.com/jasmaa/hikawa/pkg/ui"
)

const VERSION = "1.0.0"

func getImageFromFilePath(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	image, _, err := image.Decode(f)
	return image, err
}

func main() {
	wnd := g.NewMasterWindow(fmt.Sprintf("Hikawa - v%s", VERSION), 800, 600, 0)
	img, err := getImageFromFilePath("assets/icon.png")
	if err == nil {
		wnd.SetIcon([]image.Image{img})
	}
	wnd.Run(ui.Loop)
}
