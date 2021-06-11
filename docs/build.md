

## Build ##
Golang produces OS specific executables. You can create a binary targeting any OS from any OS.

To build starling, you can use the makefile. Help can be displayed by running `make help`.
```
$ make help
all                            Build for all platforms
build                          Build binaries for all platforms
ux                             Build React UX
windows                        Build for Windows (AMD 64bit)
linux                          Build for Linux (AMD 64bit)
pi                             Build for Raspberry Pi (ARM 64 bit)
darwin                         Build for Darwin (macOS)
clean                          Remove previous build
help                           Display available commands
```

#### Build binaries ####
Use the `make` command to build binaries for all platforms.
use the following commands for building only targeting specific platform:
1. Windows: `make windows` or `GOOS=windows GOARCH=amd64 go build -v -o bin/starling_windows_amd64.exe`
2. Linux: `make linux`
3. Raspberry Pi: `make pi`
4. macOS: `make darwin` 

```
$ make
env GOOS=windows GOARCH=amd64 go build -v -o ./bin/starling_windows_amd64.exe
env GOOS=linux GOARCH=amd64 go build -v -o ./bin/starling_linux_amd64
env GOOS=linux GOARCH=arm64 go build -v -o ./bin/starling_linux_arm64
env GOOS=darwin GOARCH=amd64 go build -v -o ./bin/starling_darwin_amd64
```

#### Build UX ####
UX is bundled into the Starling source. Unless you are changing the UX, there is no need to run this step.
Production ready UX can be built using `make ux` command.
This pulls in all the `npm` packages that are needed for the UX.
After building the code, it copies the built site into the [static](../pkg/serving/static) folder. 

#### Cleanup ####
To clean up all the binaries, use `make clean`

[Back to contents](../README.md)| Previous: [Development environment setup](devsetup.md) | Next: [Running server](running.md)
---------------------------------|-------------------------------------------------------|------------------------------------
