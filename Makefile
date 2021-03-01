BIN_NAME=./bin/starling
WIN_BIN_NAME=${BIN_NAME}.exe

buildwindows:
	env GOOS=windows GOARCH=amd64 go build -o ${WIN_BIN_NAME} -v

buildmac:
	env GOOS=darwin GOARCH=amd64 go build -o ${BIN_NAME} -v
buildlinux:
	env GOOS=linux GOARCH=amd64 go build -o ${BIN_NAME} -v

run:
	go run *.go

clean:
	go clean
	rm -f ${BIN_NAME}
	rm -f ${WIN_BIN_NAME}