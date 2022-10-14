SHELL := /bin/bash

GITCOMMIT := $(shell git rev-parse HEAD)
GITDATE := $(shell git show -s --format='%ct')
GITVERSION := $(shell cat package.json | jq .version)

LDFLAGSSTRING +=-X main.GitCommit=$(GITCOMMIT)
LDFLAGSSTRING +=-X main.GitDate=$(GITDATE)
LDFLAGSSTRING +=-X main.GitVersion=$(GITVERSION)
LDFLAGS :=-ldflags "$(LDFLAGSSTRING)"

CONTRACTS_PATH := "./packages/key-contract/artifacts/contracts"

key-locker:
	env GO111MODULE=on go build $(LDFLAGS)
.PHONY: key-locker

clean:
	rm key-locker

test:
	go test -v ./...

lint:
	golangci-lint run ./...

abi:
	cat $(CONTRACTS_PATH)/KeyLocker.sol/KeyLocker.json \
		| jq '{abi,bytecode}' \
		> packages/key-contract/abis/KeyLocker.json

binding: abi
	$(eval temp := $(shell mktemp))

	cat packages/key-contract/abis/KeyLocker.json \
		| jq -r .bytecode > $(temp)

	cat packages/key-contract/abis/KeyLocker.json \
		| jq .abi \
		| abigen --pkg bindings \
		--abi - \
		--out blockchain/ethereum/bindings/keylocker.go \
		--type ethereum \
		--bin $(temp)

	rm $(temp)