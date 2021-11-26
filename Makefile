BINARY_NAME=apc_compare
VERSION=1.0.1


build:
#	SET CGO_ENABLED=0 set GOARCH=amd64 set GOOS=darwin go build -o ${BINARY_NAME}-darwin main.go
#	SET CGO_ENABLED=0 set GOARCH=amd64 set GOOS=linux go build -o ${BINARY_NAME}-linux main.go

	SET CGO_ENABLED=0&&SET GOARCH=386&& SET GOOS=windows&& go build -ldflags "-s -w" -o bin/${BINARY_NAME}_x86.exe main.go
	SET CGO_ENABLED=0&&SET GOARCH=amd64&& SET GOOS=windows&& go build -ldflags "-s -w" -o bin/${BINARY_NAME}_x64.exe main.go

	tar -a -f  bin/${BINARY_NAME}_x64.zip -c conf\*.*  -C ".\bin" ${BINARY_NAME}_x64.exe
	tar -a -f  bin/${BINARY_NAME}_x86.zip -c conf\*.*  -C ".\bin" ${BINARY_NAME}_x86.exe
run:
	./bin/${BINARY_NAME}_x64.exe

release:
	tar -a -f  ./release/${BINARY_NAME}_${VERSION}.zip -c .\bin\${BINARY_NAME}_x64.zip  -C ".\bin" ${BINARY_NAME}_x86.zip
.PHONY : release
clean:
	go clean
	rd /s /q "./bin"
