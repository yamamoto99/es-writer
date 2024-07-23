#!/bin/sh
set -e

# マイグレーションを実行
#　マイグレーションが成功しなかった場合、再度実行する
until /app/migrate/main; do
  echo "マイグレーションに失敗しました。再試行します..."
  sleep 1
done

# APIサーバーを起動
exec /app/main