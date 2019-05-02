# Basic go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Binary name
BINARY_NAME=artemis
    
		
all: test build
build: 
	$(GOBUILD) -o $(BINARY_NAME) -v
test: 
	$(GOTEST) -v ./...
clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME) solution.log
run:
	$(GOBUILD) -o $(BINARY_NAME) -v
	./$(BINARY_NAME) schedule -d example/datacenter.json -w example/workload.json
deps:
	$(GOGET) github.com/spf13/cobra
    
