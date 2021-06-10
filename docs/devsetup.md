
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

[Back to contents](../Readme.md) | Next: [Building](build.md)
---------------------------------|---------------------------
