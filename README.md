# user_balance_microservice

## Локальный запуск

1) Убедитесь, что у вас установлены [Golang](https://go.dev/doc/install) и [Docker](https://docs.docker.com/get-docker/).

2) Выполните следующие команды в консоли:

```sh
$ git clone https://github.com/petrashin/user_balance_microservice.git
$ cd user_balance_microservice
```

3) Создайте файл .env по примеру, приведенному в env.example, где:
- DB_TYPE тип СУБД, который вы используете (mysql или postgres)
- DB_USERNAME ваш логин для подключения к базе данных (в mysql по умолчанию root)
- DB_PASSWORD ваш пароль для подключения к базе данных (в mysql по умолчанию root)
- IP ваш текущий IP
- DB_POPT порт для подключения к базе данных (в mysql по умолчанию 3306)
- DB_NAME название базы данных, которую вы создадите в пункте 3 с (например, avito)

4) Создайте базу данных, импортировав предоставленный скрипт db.sql

5) Запустите тесты:

```sh
$ go test -v
```

6) Запустите docker

7) Создайте и запустите веб-приложение:

```sh
$ docker compose up --build
```

## Эндпоинты API:
- (POST) http://127.0.0.1:8080/update_balance/?id={{id}}&balance={{balance}} - метод начисления средств на баланс <br>
- (POST) http://127.0.0.1:8080/reserve_money/?user_id={{id}}&service_id={{service_id}}&order_id={{order_id}}&price={{price}} - метод резервирования средств с основного баланса на отдельном счете <br>
- (POST) http://127.0.0.1:8080/revenue_recognition/?user_id={{id}}&service_id={{service_id}}&order_id={{order_id}}&price={{price}} - метод признания выручки <br>
- (GET) http://127.0.0.1:8080/get_balance/?id={{id}} - метод получения баланса пользователя <br>
- (POST) http://127.0.0.1:8080/unreserve_money/?user_id={{id}}&service_id={{service_id}}&order_id={{order_id}}&price={{price}} - метод разрезервирования денег, если услугу применить не удалось <br>
