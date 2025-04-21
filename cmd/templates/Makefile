.PHONY : dev prod compile

dev:
	echo "Compiling..."
	go build -o isaac . 

prod:
	echo "Compiling..."
	go build -ldflags="-s -w" -o isaac .

compile:
	echo "Compiling for every OS and Platform"
	GOOS=freebsd GOARCH=amd64 go build -ldflags="-s -w" -o bin/{{.ProjectName}}-bsd-amd64 .
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/{{.ProjectName}}-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o bin/{{.ProjectName}}-linux-arm64 .