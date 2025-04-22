.PHONY : dev prod install compile

dev:
	echo "Compiling..."
	go build -o bin/jrx . 

prod:
	echo "Compiling..."
	go build -ldflags="-s -w" -o bin/jrx .

move:
	mv bin/jrx ~/bin/jrx

install: prod move

compile:
	echo "Compiling for every OS and Platform"
	GOOS=freebsd GOARCH=amd64 go build -ldflags="-s -w" -o bin/jrx-freebsd-amd64 .
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/jrx-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o bin/jrx-linux-arm64 .
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o bin/jrx-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o bin/jrx-darwin-m1 .