
## Development Environment Setup ##

Starling backend is built in Golang and the frontend is built in ReactJs.
You can build Starling on Windows, MacOS and Linux.

If you are using Windows, it is recommended to use [WSL 2 (Windows Sub System for Linux)](https://docs.microsoft.com/en-us/windows/wsl/install-win10).

### Golang ###
Follow the instructions to [install Go](https://golang.org/doc/install). Pick the appropriate package to install the
latest 1.16.x release of Go. This will give you access to the Go toolchain and compiler.

- If you are on Windows, use the MSI to install. It will set the necessary environment variables.
- If you installed via the tarball, you will need to add a GOROOT environment variable pointing to the
  folder where you installed Go (typically /usr/local/go on linux-based systems)
- You should also check to make sure that you can access the Go compiler and tools. They are available at $GOROOT/bin
  (or $GOROOT\bin) and should be added to your path if they are not already. You can verify this by running the following:
    - Max/Linux: `which go`
    - Windows (CMD): `where go`

### GCC - GNU Tools ###
One of the packages used in the server depends on gcc (GNU C compiler).
You can install GCC using appropriate means e.g. *sudo apt-get install gcc*.

Also, you need **make** to build based on makefile. On windows you can use *nmake.exe*.

### Node JS ###
Starling UX is built using ReactJS. The pre-built optimized packages are checked into [static](pkg/serving/static) folder.
However if you want to modify the UX, you need to install the following:
- [NodeJS](https://nodejs.org/en/download/)
- [Yarn](https://yarnpkg.com/getting-started/install)


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

[Back to contents](../README.md)| Next: [Running server](running.md)
--------------------------------|------------------------------------
