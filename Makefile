.PHONY: download
download: 
	go mod download

.PHONY: tidy
tidy: 
	go mod tidy

.PHONY: test
test: 
	go test ./... -v

