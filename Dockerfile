FROM --platform=${BUILDPLATFORM} alpine:3.13.1

ARG TARGETOS
ARG TARGETARCH
ARG TARGETPLATFORM

LABEL maintainer="Gareth Evans <gareth@bryncynfelin.co.uk>"
COPY dist/jv-${TARGETOS}_${TARGETOS}_${TARGETARCH}/jv /usr/bin/jv

ENTRYPOINT [ "/usr/bin/jv" ]

CMD ["--help"]
