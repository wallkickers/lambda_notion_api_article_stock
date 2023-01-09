# ベースとなるDockerイメージ指定
FROM alpine
RUN apk --no-cache add bash go git ca-certificates curl vim

ENV CGO_ENABLED=0

# コンテナログイン時のディレクトリ指定
WORKDIR /go/src/app
