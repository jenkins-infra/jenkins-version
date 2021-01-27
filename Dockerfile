FROM --platform=${BUILDPLATFORM} golang:1.15-alpine AS builder

ARG TARGETOS
ARG TARGETARCH
ARG TARGETPLATFORM

WORKDIR /go/src/app
COPY . .

ARG version
ARG build_date
ARG sha

ENV CGO_ENABLED=0

RUN go version

RUN GOOS=$TARGETOS GOARCH=$TARGETARCH make build

RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -v -a \
  -ldflags "-w -s \
    -X github.com/garethjevans/jenkins-version/pkg/version.BuildDate=$build_date \
    -X github.com/garethjevans/jenkins-version/pkg/version.Version=$version \
    -X github.com/garethjevans/jenkins-version/pkg/version.GoVersion=\'"$(go version)"\' \
    -X github.com/garethjevans/jenkins-version/pkg/version.Sha1=$sha" \
  -o bin/jv cmd/jv/jv.go

FROM alpine:3.12.3

LABEL maintainer="Gareth Evans <gareth@bryncynfelin.co.uk>"
COPY --from=builder /go/src/app/bin/jv /usr/bin/jv

ENTRYPOINT [ "/usr/bin/jv" ]

CMD ["--help"]
