.PHONY: build install test vet clean skill

BINARY=asana-cli
SKILL_DIR=$(HOME)/.claude/skills/asana

build:
	go build -o $(BINARY) .

install: build skill
	go install .

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
