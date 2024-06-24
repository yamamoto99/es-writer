# Goの公式イメージをベースにする
FROM golang:latest

# 作業ディレクトリを設定
WORKDIR /app

# Goのモジュールを有効にする
ENV GO111MODULE=on

# 必要なGoのパッケージをダウンロード
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# ソースコードをコピー
COPY backend/ ./

# APIサーバーをビルド
RUN go build -o main .

# ポート8080を公開
EXPOSE 8080

# APIサーバーを実行
CMD ["./main"]