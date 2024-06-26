ARG BUILDER_IMAGE=golang
ARG BUILDER_VERSION=1.18-alpine

FROM $BUILDER_IMAGE:$BUILDER_VERSION AS builder

WORKDIR /go/src/app

ENV GOPRIVATE=github.com/anoideaopen
#TODO  actaul regisrte URL and Logopass
ARG APP_VERSION=unknown

RUN echo "$REGISTRY_NETRC" > ~/.netrc

COPY go.mod go.sum ./
RUN apk add git=~2 binutils=~2 upx=~3 && CGO_ENABLED=0 go mod download

COPY . .
RUN CGO_ENABLED=0 go build -v -ldflags="-X 'main.AppInfoVer=$APP_VERSION'" -o /go/bin/app && strip /go/bin/app && upx -5 -q /go/bin/app

FROM alpine:3.15
COPY --chown=65534:65534 --from=builder /go/bin/app /
USER 65534

ENTRYPOINT [ "/app" ]
