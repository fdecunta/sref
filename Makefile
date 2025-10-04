install: sref.go
	go build sref.go
	sudo mv sref /usr/local/bin

uninstall:
	sudo rm /usr/local/bin/sref

build: sref.go
	go build sref.go

.PHONY: install uninstall build
