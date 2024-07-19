# Goの公式イメージをベースにする
FROM golang:latest

# 作業ディレクトリを設定
WORKDIR /app

# Goのモジュールを有効にする
ENV GO111MODULE=on

# 必要なGoのパッケージをダウンロード
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピー
COPY ./ ./

# マイグレーションツールとAPIサーバーをビルド
RUN go build -o /app/migrate/main ./migrate
RUN go build -o /app/main .

# ポート8080を公開
EXPOSE 8080

# マイグレーションを実行してからAPIサーバーを起動するスクリプトを作成
COPY start.sh /app/start.sh
RUN chmod +x /app/start.sh

# スクリプトを実行
CMD ["/app/start.sh"]