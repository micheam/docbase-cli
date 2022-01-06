TARGET = docbase
REV_PARSE = git rev-parse --short HEAD
BIN_DIR = 

$(TARGET) : clean cmd/$(TARGET)/main.go
	@go build -ldflags "-X main.Githash=`$(REV_PARSE)`" -o ./bin/$(TARGET) ./cmd/$(TARGET)

.PHONY: clean
clean :
	@rm -f ./bin/$(TARGET)

.PHONY: test
test : 
	@go test ./...

.PHONY: install
install : test $(TARGET)
	@cp ./bin/$(TARGET) $(GOBIN)
