# user_balance_microservice

POST http://127.0.0.1:8080/update_balance/?id=32&balance=10000 - update user balance
POST http://127.0.0.1:8080/reserve_money/?user_id=31&service_id=1&order_id=1&price=10000 - reserve money
POST http://127.0.0.1:8080/revenue_recognition/?user_id=6&service_id=1&order_id=1&price=15000 - write of to revenue
GET http://127.0.0.1:8080/get_balance/?id=31 - get user balance
POST http://127.0.0.1:8080/unreserve_money/?user_id=31&service_id=1&order_id=1&price=10000 - unreserve money
