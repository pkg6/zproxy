build:
	go build -ldflags="-s -w" -o zproxy main.go
	$(if $(shell command -v upx), upx zproxy)

mac:
	GOOS=darwin go build -ldflags="-s -w" -o zproxy-darwin .
	$(if $(shell command -v upx), upx zproxy-darwin)

win:
	GOOS=windows go build -ldflags="-s -w" -o zproxy.exe .
	$(if $(shell command -v upx), upx zproxy.exe)

linux:
	GOOS=linux go build -ldflags="-s -w" -o zproxy-linux .
	$(if $(shell command -v upx), upx zproxy-linux)
