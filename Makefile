# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get

# Main package
MAIN_PACKAGE = ./cmd/main.go

# Output binary
OUTPUT_BINARY = gertygo

# Build the project
build:
	$(GOBUILD) -o $(OUTPUT_BINARY) $(MAIN_PACKAGE)
	chmod +x $(OUTPUT_BINARY)

# Clean the project
clean:
	$(GOCLEAN)
	rm -f $(OUTPUT_BINARY)

# Run the project
run: build
	./$(OUTPUT_BINARY)

# Install project dependencies
deps:
	$(GOGET) ./...

# Run all tests
test:
	$(GOTEST) ./...

.PHONY: build clean run deps test