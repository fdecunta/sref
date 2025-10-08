CONFIG_DIR := $(HOME)/.config/sref
LOCALBIN := /usr/local/bin

install: sref.go
	go build sref.go
	mkdir -p $(CONFIG_DIR)
	touch $(CONFIG_DIR)/email.conf
	touch $(CONFIG_DIR)/references.json
	sudo mv sref $(LOCALBIN)/

uninstall:
	rm -rf $(CONFIG_DIR)
	sudo rm -f $(LOCALBIN)/sref

build: sref.go
	go build sref.go

.PHONY: install uninstall build
