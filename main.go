package main

//go:generate go run cmd/main.go --gdnative --types --classes

import (
	_ "github.com/jasmaa/hikawa/pkg/export"
	_ "github.com/jasmaa/hikawa/pkg/ui"
)

func main() {
}
