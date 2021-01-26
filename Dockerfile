FROM golang:1.15 as builder

WORKDIR /go/src/app

COPY . .

RUN go get -d -v ./...

RUN \
  go build -v -a \
  -ldflags "-w -s \
    -X \"github.com/garethjevans/jenkins-version/pkg/version.BuildDate=`date -R`\" \
    -X \"github.com/garethjevans/jenkins-version/pkg/version.GoVersion=`go version`\" \
    -X \"github.com/garethjevans/jenkins-version/pkg/version.Version=`git describe --tags`\""\
  -o bin/jv

FROM alpine:3.12.3

LABEL maintainer="Gareth Evans <gareth@bryncynfelin.co.uk>"
COPY --from=builder /go/src/app/bin/jv /usr/bin/jv

ENTRYPOINT [ "/usr/bin/jv" ]

CMD ["--help"]