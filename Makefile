.PHONY: build install test vet clean skill

BINARY=asana-cli
SKILL_DIR=$(HOME)/.claude/skills/asana

VERSION ?= dev
LDFLAGS = -X github.com/danilodrobac/asana-cli/cmd.Version=$(VERSION)

build:
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) .

install: skill
	go install -ldflags "$(LDFLAGS)" .

test:
	go test ./...

vet:
	go vet ./...

clean:
	rm -f $(BINARY)

skill:
	mkdir -p $(SKILL_DIR)
	cp skill/SKILL.md $(SKILL_DIR)/SKILL.md
	@echo "Skill installed to $(SKILL_DIR)"

uninstall:
	rm -f $(shell go env GOPATH)/bin/$(BINARY)
	rm -rf $(SKILL_DIR)
