# ビルドステージ
FROM golang:1.20-alpine AS builder

WORKDIR /app

# Go モジュールのキャッシュ
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピーしてビルド
COPY . .
RUN go build -o main .

# 実行ステージ
FROM alpine:latest

WORKDIR /app

# ビルドしたバイナリをコピー
COPY --from=builder /app/main .

# 必要な証明書をインストール（MySQL ドライバ用）
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

# ポートを開放
EXPOSE 8080

# アプリケーションを実行
CMD ["./main"]
