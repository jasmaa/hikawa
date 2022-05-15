build:
	go build -ldflags "-s -w -H=windowsgui -extldflags=-static" cmd/main.go

test:
	go test ./...

clean:
	rm *.exe