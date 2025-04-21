FROM golang:{{.GoVersion}}

ADD . /go/src/{{.ProjectName}}

WORKDIR /go/src/{{.ProjectName}}

RUN go build -o /go/bin/{{.ProjectName}} .
