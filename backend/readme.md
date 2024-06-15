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
## Terminalからpostgresを操作する
指定したコンテナ内で、インタラクティブなbashシェルを起動する
```shell
docker exec -it <-コンテナID-> bash 
```
PostgreSQLに接続する。ユーザー名はpostgresで、データベースはtestdb
```shell
psql -U postgres -d testdb
```