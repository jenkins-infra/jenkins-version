FROM --platform=${BUILDPLATFORM} curlimages/curl:7.83.1 AS build-stage0

ARG TARGETOS
ARG TARGETARCH
ARG TARGETPLATFORM
ARG JV_VERSION

RUN curl -L -o /tmp/jv-${TARGETOS}-${TARGETARCH}.tar.gz https://github.com/jenkins-infra/jenkins-version/releases/download/${JV_VERSION}/jenkins-version-${TARGETOS}-${TARGETARCH}.tar.gz && \
      tar -xvzf /tmp/jv-${TARGETOS}-${TARGETARCH}.tar.gz -C /tmp && \
      chmod a+x /tmp/jv

FROM --platform=${BUILDPLATFORM} alpine:3.15.2
LABEL maintainer="Gareth Evans <gareth@bryncynfelin.co.uk>"

COPY --from=build-stage0 /tmp/jv /usr/bin/jv
COPY github-actions-entrypoint.sh /usr/bin/github-actions-entrypoint.sh

ENTRYPOINT [ "/usr/bin/jv" ]
CMD ["--help"]
