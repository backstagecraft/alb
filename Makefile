.PHONY: tags

all: test

tags:
	gotags -R -f tags .

test:
	go test ./app ./cmd/bvsd ./cmd/bvscli

install:
	cp -f ./testdata/genesis.json $(HOME)/.bvsd/config/
	go install ./cmd/bvsd ./cmd/bvscli
