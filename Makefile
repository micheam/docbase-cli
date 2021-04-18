TARGET = docbase
REV_PARSE = git rev-parse --short HEAD
BIN_DIR = 

$(TARGET) : cmd/$(TARGET)/main.go
	@go build -ldflags "-X main.Githash=`$(REV_PARSE)`" -o ./bin/$(TARGET) ./cmd/$(TARGET)

.PHONY: claen
clean :
	@rm -f ./bin/$(TARGET)

.PHONY: test
test : 
	@go test ./...

.PHONY: install
install : $(TARGET)
	@cp ./bin/$(TARGET) $(GOBIN)/$(TARGET)
