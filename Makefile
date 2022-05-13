build:
	go build cmd/main.go

test:
	go test ./...

clean:
	rm *.exe