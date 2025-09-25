build: sref.go
	go build sref.go

run:
	./sref -a "10.1111/j.1461-0248.2008.01192.x"

.PHONY: build run
