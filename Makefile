build:
	rsrc -manifest main.manifest -o rsrc.syso
	go build -ldflags="-H windowsgui"

clean:
	rm rsrc.syso *.exe