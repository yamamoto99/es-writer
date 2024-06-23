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
## pgAdminからDBを操作する
localhost5050にアクセス
```shell
http://localhost:5050
```
下記でログイン
```
Email Address/Username:pgadmin4@pgadmin.org
Password:admin
```
DBに接続する
```
Servers->登録->サーバー
名前:自由
ホスト名/アドレス:db
管理者用データベース:postgres
ユーザー名：postgres
パスワード:postgres
```
## EC2接続方法
EC2に設定されたキーペアがある状態で
```
ssh -i <秘密鍵名> <ユーザー名>@<パブリックIP>
```
## EC2からRDSインスタンス内のDBにアクセスする
```
psql --host=<DBエンドポイント> --port=<portnum> --username=<username> --dbname<dbname>
```
