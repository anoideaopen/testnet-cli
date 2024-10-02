###############################################################################
# Build image
###############################################################################

ARG UBUNTU_VER
FROM ubuntu:${UBUNTU_VER:-22.04} AS builder

ARG TARGETARCH
ARG TARGETOS
ARG GO_VER
ARG APP_VER=unknown
ARG COMMIT=unknown
ARG BUILD_TIME=unknown

RUN apt update && apt install -y \
    git \
    gcc \
    curl \
    make

RUN curl -sL https://go.dev/dl/go${GO_VER}.${TARGETOS}-${TARGETARCH}.tar.gz | tar zxf - -C /usr/local
ENV PATH="/usr/local/go/bin:$PATH"

ADD . .

RUN CGO_ENABLED=0 go build -v -ldflags="-X 'main.version=$APP_VER' -X 'main.commit=$COMMIT' -X 'main.date=$BUILD_TIME'" -o /go/bin/cli

###############################################################################
# Runtime image
###############################################################################

ARG UBUNTU_VER
FROM ubuntu:${UBUNTU_VER:-22.04}

ARG APP_VER

# set up nsswitch.conf for Go's "netgo" implementation
# - https://github.com/golang/go/blob/go1.9.1/src/net/conf.go#L194-L275
# - docker run --rm debian:stretch grep '^hosts:' /etc/nsswitch.conf
RUN echo 'hosts: files dns' > /etc/nsswitch.conf

ENV APP_VER=${APP_VER}

COPY    --chown=65534:65534 --from=builder /go/bin/cli /
USER 65534

ENTRYPOINT [ "/cli" ]
