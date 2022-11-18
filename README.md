# user_balance_microservice

## Running Locally

1) Make sure you have [Golang](https://go.dev/doc/install) and [Docker](https://docs.docker.com/get-docker/) installed.

```sh
$ git clone https://github.com/petrashin/user_balance_microservice.git # or clone your own fork
$ cd user_balance_microservice
```

2) Create .env file as env.example where:
- DB_TYPE is your database DBMS (mysql or postgres)
- DB_USERNAME is your database username (root is default in mysql)
- DB_PASSWORD is your database password (root is default in mysql)
- IP is your current ip
- DB_POPT is your database port (for mysql 3306 is default)
- DB_NAME is your database name which you will create in 3-rd paragraph (e.g. avito)

3) Create database with provided db.sql script

4) Run provided tests:

```sh
$ go test -v
```

5) Start docker daemon

6) Build and deploy web application service:

```sh
$ docker compose up --build
```

## API endpoints:
- (POST) http://127.0.0.1:8080/update_balance/?id={{id}}&balance={{balance}} - update user balance <br>
- (POST) http://127.0.0.1:8080/reserve_money/?user_id={{id}}&service_id={{service_id}}&order_id={{order_id}}&price={{price}} - reserve money <br>
- (POST) http://127.0.0.1:8080/revenue_recognition/?user_id={{id}}&service_id={{service_id}}&order_id={{order_id}}&price={{price}} - write of to revenue <br>
- (GET) http://127.0.0.1:8080/get_balance/?id={{id}} - get user balance <br>
- (POST) http://127.0.0.1:8080/unreserve_money/?user_id={{id}}&service_id={{service_id}}&order_id={{order_id}}&price={{price}} - unreserve money <br>
