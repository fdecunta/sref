CONFIG_DIR := $(HOME)/.config/sref
LOCALBIN := /usr/local/bin

install: sref-db.go
	go build sref-db.go
	mkdir -p $(CONFIG_DIR)
	touch $(CONFIG_DIR)/email.conf
	touch $(CONFIG_DIR)/references.json
	sudo mv sref-db $(LOCALBIN)/

uninstall:
	rm -rf $(CONFIG_DIR)
	sudo rm -f $(LOCALBIN)/sref-db

test:
	go test -v

.PHONY: install uninstall test
