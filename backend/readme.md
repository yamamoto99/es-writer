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
## phpmyadminからDBを操作する
localhost5050にアクセス
```shell
http://localhost:5050
```
下記でログイン
```
user: user
Password: password
```
>[!CAUTION]
>現在はEC2にデプロイしていないです
## EC2接続方法
EC2に設定されたキーペアがある状態で
```
ssh -i <秘密鍵名> <ユーザー名>@<パブリックIP>
```
## EC2からRDSインスタンス内のDBにアクセスする
```
psql -h=<DBエンドポイント>　--U=<username> -d<dbname>
```
