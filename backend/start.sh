#!/bin/sh
set -e

# マイグレーションを実行
/app/migrate/main

# APIサーバーを起動
exec /app/main