EXECUTABLE=./bin/starling
WINDOWS=$(EXECUTABLE)_windows_amd64.exe
LINUX=$(EXECUTABLE)_linux_amd64
DARWIN=$(EXECUTABLE)_darwin_amd64
WEBUX=./webux
STATIC_CONTENT=./pkg/serving/static

.PHONY: all clean

## Build for all platforms
all: build						## Build for all platforms

# build binaries
build: windows linux darwin		## Build binaries for all platforms

ux:								## Build React UX
	cd $(WEBUX) && yarn install && yarn build
	rm -rf $(STATIC_CONTENT)
	mkdir $(STATIC_CONTENT)
	mv $(WEBUX)/build/* $(STATIC_CONTENT)

windows:						## Build for Windows
	env GOOS=windows GOARCH=amd64 go build -v -o $(WINDOWS)

linux:							## Build for Linux
	env GOOS=linux GOARCH=amd64 go build -v -o $(LINUX)

darwin:							## Build for Darwin (macOS)
	env GOOS=darwin GOARCH=amd64 go build -v -o $(DARWIN)

clean:							## Remove previous build
	go clean
	rm -f $(WINDOWS_EXECUTABLE) $(LINUX_EXECUTABLE) $(DARWIN_EXECUTABLE)
	rm -rf $(WEBUX)/build

help: 							## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
