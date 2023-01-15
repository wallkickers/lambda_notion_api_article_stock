FROM golang:1.19.1-alpine
RUN apk update && apk --no-cache add bash git ca-certificates curl vim

ENV CGO_ENABLED=0

# コンテナログイン時のディレクトリ指定
WORKDIR /go/src/app
