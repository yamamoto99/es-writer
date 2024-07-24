# Backend
## docker 利用方法
Docker Composeを使ってポート8080でサービスを起動
```shell
docker compose up -d --build
```
Docker Composeを使ってサービスを停止し、ボリュームも削除
```shell
docker compose down -v
```
noneイメージを全削除できる有能コマンド
```shell
docker image prune
```
## TerminalからDBを操作する
動作しているコンテナの一覧を表示
```shell
docker ps
```
指定したコンテナ内で、インタラクティブなbashシェルを起動する
```shell
docker exec -it <-コンテナID-> bash 
```
PostgreSQLに接続する。ユーザー名はpostgresで、データベースはtestdb
```shell
psql -U postgres -d testdb
```
## テーブル構造
| フィールド名 | データ型 | JSON タグ | GORM タグ |
|-------------|---------|-----------|-----------|
| UserID      | string  | id        | gorm:unique not null |
| Username    | string  | username  | unique not null |
| Email       | string  | email     | - |
| Bio         | string  | bio       | - |
| Experience  | string  | experience| - |
| Projects    | string  | projects  | - |
| CreatedAt   | time.Time| created_at| - |
| UpdatedAt   | time.Time| updated_at| - |

## EC2接続方法
>[!CAUTION]
>現在はEC2にデプロイしていないです

EC2に設定されたキーペアがある状態で
```
ssh -i <秘密鍵名> <ユーザー名>@<パブリックIP>
```
## EC2からRDSインスタンス内のDBにアクセスする
```
psql -h=<DBエンドポイント>　--U=<username> -d<dbname>
```
