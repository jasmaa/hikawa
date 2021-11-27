build:
	go run mage.go -v build

test:
	go test ./...

clean:
	go run mage.go -v clean
	rm *.pck *.exe