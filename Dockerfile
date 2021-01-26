FROM golang:1.15 as builder

# Install TARGETPLATFORM parser to translate its value to GOOS, GOARCH, and GOARM
COPY --from=tonistiigi/xx:golang / /
# Bring TARGETPLATFORM to the build scope
ARG TARGETPLATFORM

WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...

ARG version
ARG build_date
ARG sha

RUN \
  go build -v -a \
  -ldflags "-w -s \
    -X github.com/garethjevans/jenkins-version/pkg/version.BuildDate=$build_date \
    -X github.com/garethjevans/jenkins-version/pkg/version.Version=$version \
    -X github.com/garethjevans/jenkins-version/pkg/version.Sha1=$sha" \
  -o bin/jv

FROM alpine:3.12.3

LABEL maintainer="Gareth Evans <gareth@bryncynfelin.co.uk>"
COPY --from=builder /go/src/app/bin/jv /usr/bin/jv

ENTRYPOINT [ "/usr/bin/jv" ]

CMD ["--help"]
