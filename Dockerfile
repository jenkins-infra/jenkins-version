FROM --platform=${BUILDPLATFORM} alpine:3.13.0

ARG TARGETOS
ARG TARGETARCH
ARG TARGETPLATFORM

LABEL maintainer="Gareth Evans <gareth@bryncynfelin.co.uk>"
COPY dist/jv-$TARGETARCH_$TARGETARCH_$TARGETOS/jv /usr/bin/jv

ENTRYPOINT [ "/usr/bin/jv" ]

CMD ["--help"]
