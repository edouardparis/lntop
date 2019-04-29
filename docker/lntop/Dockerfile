FROM golang:1.12-alpine as builder

# install build dependencies
RUN apk add --no-cache --update git gcc musl-dev

ARG LNTOP_SRC_PATH

WORKDIR /root/_build

# we want to populate the module cache based on the go.{mod,sum} files.
COPY "$LNTOP_SRC_PATH/go.mod" .
COPY "$LNTOP_SRC_PATH/go.sum" .

WORKDIR $GOPATH/src/github.com/edouardparis/lntop
COPY "$LNTOP_SRC_PATH" .

ENV GO111MODULE=on
RUN go install -mod=vendor ./...

# ---------------------------------------------------------------------------------------------------------------------------

FROM golang:1.12-alpine as final

RUN apk add --no-cache \
    bash fish \
    ca-certificates \
    tini

ENTRYPOINT ["/sbin/tini", "--"]

ENV PATH $PATH:/root

ARG LNTOP_CONF_PATH

# copy the binaries and entrypoint from the builder image.
COPY --from=builder /go/bin/lntop /bin/

WORKDIR /root

COPY "home" .