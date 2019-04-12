linux: statik
	CGO_ENABLED=1 CC=gcc GOOS=linux GOARCH=amd64 go build -tags static -ldflags "-s -w" -o gpp-linux gpp.go

freebsd: statik
	CGO_ENABLED=1 CC=gcc GOOS=freebsd GOARCH=amd64 go build -tags static -ldflags "-s -w" -o gpp-freebsd gpp.go

macos: statik
	CGO_ENABLED=1 CC=gcc GOOS=darwin GOARCH=amd64 go build -tags static -ldflags "-s -w" -o gpp-macos gpp.go

windows: statik
	CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go build -tags static -ldflags "-s -w -H=windowsgui" -o gpp-windows.exe gpp.go


rpi: statik
	CGO_ENABLED=1 CC=gcc GOOS=linux GOARCH=arm go build -o gpp-rpi gpp.go

statik:
	statik -src fonts/

clean:
	rm -rf statik gpp-linux gpp-freebsd gpp-macos gpp-windows.exe
