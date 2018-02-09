mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
current_dir := $(patsubst %/,%,$(dir $(mkfile_path)))

export GOBIN=$(current_dir)/bin

GOCC := gorunpkg github.com/goccmack/gocc
GOFLAGS :=

.PHONY: all clean test

all: vendor bin/lexer bin/codegen

debug: GOFLAGS += -tags debug
debug: all

test-debug: GOFLAGS += -tags debug
test-debug: test

vendor:
	dep ensure -v

bin/lexer: src/main/lexer/lexer.go src/lexer/lexer.go src/lexer/*.go
	@echo -e "\e[1;32mCompiling Lexer \e[0m"
	go install $(GOFLAGS) $(current_dir)/src/main/lexer/lexer.go

bin/codegen: src/main/codegen/codegen.go src/codegen/*.go
	@echo -e "\e[1;32mCompiling Codegen \e[0m"
	go install $(GOFLAGS) $(current_dir)/src/main/codegen/codegen.go

src/lexer/lexer.go: src/tango.ebnf
	@echo -e "\e[1;33mGenerating Lexer \e[0m"
	cd $(current_dir)/src && $(GOCC) tango.ebnf

test:
	go test $(GOFLAGS) src/lexer

clean:
	@echo -e "\e[1;31mCleaning Files \e[0m"
	@echo -e "\e[1;31m  Clearing pkg and bin \e[0m"
	@rm -rf $(current_dir)/pkg $(current_dir)/bin/**
	@echo -e "\e[1;31m  Clearing generated files \e[0m"
	@rm -rf $(current_dir)/src/util
	@rm -rf $(current_dir)/src/token
	@rm -rf $(current_dir)/src/lexer/lexer.go
	@rm -rf $(current_dir)/src/lexer/acttab.go
	@rm -rf $(current_dir)/src/lexer/transitiontable.go

# Prepare for submitting. Warning don't run this lightly.
nuke: clean
	@echo -e "\e[1;31m  Clearing downloaded libraries \e[0m"
	@rm -rf $(current_dir)/vendor
	@echo -e "\e[1;31m  Clearing git stuff \e[0m"
	@rm -rf $(current_dir)/.git
	@rm -rf $(current_dir)/.gitignore
	@rm -rf $(current_dir)/README.md
	@rm -rf $(current_dir)/tango.ebnf
