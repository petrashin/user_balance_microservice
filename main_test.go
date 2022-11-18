package main

import (
  "os"
  "testing"
  "net/http"
  "net/http/httptest"
  "github.com/gorilla/mux"
  "github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
  DB_TYPE, _ := os.LookupEnv("DB_TYPE")
  DB_USERNAME, _ := os.LookupEnv("DB_USERNAME")
  DB_PASSWORD, _ := os.LookupEnv("DB_PASSWORD")
  IP, _ := os.LookupEnv("IP")
  DB_PORT, _ := os.LookupEnv("DB_PORT")
  DB_NAME, _ := os.LookupEnv("DB_NAME")

  if DB_TYPE == "" || DB_USERNAME == "" || DB_PASSWORD == "" || IP == "" || DB_PORT == "" || DB_NAME == "" {
    t.Error("Problem with getting params from .env file")
  }

  err := init_db(DB_TYPE, DB_USERNAME, DB_PASSWORD, IP, DB_PORT, DB_NAME)

  if err != nil {
    t.Error("Problem with connection to database")
  }
}

func Router() *mux.Router {
  rtr := mux.NewRouter()
  rtr.HandleFunc("/update_balance/", update_balance).Methods("POST")
  rtr.HandleFunc("/reserve_money/", reserve_money).Methods("POST")
  rtr.HandleFunc("/revenue_recognition/", revenue_recognition).Methods("POST")
  rtr.HandleFunc("/get_balance/", get_balance).Methods("GET")
  rtr.HandleFunc("/unreserve_money/", unreserve_money).Methods("POST")
  return rtr
}

func create_test_user (t *testing.T) {
  insert, err := db.Query("INSERT INTO `users`(`id`, `balance`) VALUES (9999, 1000)")
  if err != nil {
    t.Error("Problem with inserting into database")
  }
  defer insert.Close()
}

func create_test_reservation (t *testing.T) {
  insert, err := db.Query("INSERT INTO `reservations`(`user`, `service`, `order_id`, `price`) VALUES (9999, 1, 1, 1000)")
  if err != nil {
    t.Error("Problem with inserting into database")
  }
  defer insert.Close()
}

func delete_test_revenue_recognition (t *testing.T) {
  delete, err := db.Query("DELETE FROM `revenue` WHERE user = 9999 AND service = 1 AND order_id = 1 AND price = 1000")
  if err != nil {
    t.Error("Problem with deleting from database")
  }
  defer delete.Close()
}

func delete_test_user (t *testing.T) {
  delete, err := db.Query("DELETE FROM `users` WHERE id = 9999")
  if err != nil {
    t.Error("Problem with deleting from database")
  }
  defer delete.Close()
}

func delete_test_reservation (t *testing.T) {
  delete, err := db.Query("DELETE FROM `reservations` WHERE user = 9999 AND service = 1 AND order_id = 1 AND price = 1000")
  if err != nil {
    t.Error("Problem with deleting from database")
  }
  defer delete.Close()
}

func TestUpdatingBalance(t *testing.T) {
  // bad request
  request1, _ := http.NewRequest("POST", "/update_balance/", nil)
  response1 := httptest.NewRecorder()
  Router().ServeHTTP(response1, request1)
  assert.Equal(t, 400, response1.Code, "Bad request response is expected")

  // not enough params
  request2, _ := http.NewRequest("POST", "/update_balance/?id=9999", nil)
  response2 := httptest.NewRecorder()
  Router().ServeHTTP(response2, request2)
  assert.Equal(t, 400, response2.Code, "Bad request response is expected")

  // status 200
  request3, _ := http.NewRequest("POST", "/update_balance/?id=9999&balance=1000", nil)
  response3 := httptest.NewRecorder()
  Router().ServeHTTP(response3, request3)
  delete_test_user(t)
  assert.Equal(t, 200, response3.Code, "OK response is expected")
}

func TestGettingBalance(t *testing.T) {
  // bad request
  request1, _ := http.NewRequest("GET", "/get_balance/", nil)
  response1 := httptest.NewRecorder()
  Router().ServeHTTP(response1, request1)
  assert.Equal(t, 400, response1.Code, "Bad request response is expected")

  // user doesnt exist
  request2, _ := http.NewRequest("GET", "/get_balance/?id=0", nil)
  response2 := httptest.NewRecorder()
  Router().ServeHTTP(response2, request2)
  assert.Equal(t, 404, response2.Code, "Not Found response is expected")

  // status 200
  create_test_user(t)
  request3, _ := http.NewRequest("GET", "/get_balance/?id=9999", nil)
  response3 := httptest.NewRecorder()
  Router().ServeHTTP(response3, request3)
  delete_test_user(t)
  assert.Equal(t, 200, response3.Code, "OK response is expected")
}

func TestReserveMoney(t *testing.T) {
  // bad request
  request1, _ := http.NewRequest("POST", "/reserve_money/", nil)
  response1 := httptest.NewRecorder()
  Router().ServeHTTP(response1, request1)
  assert.Equal(t, 400, response1.Code, "Bad request response is expected")

  // not enough params
  request2, _ := http.NewRequest("POST", "/reserve_money/?user_id=9999&service_id=1&order_id=1", nil)
  response2 := httptest.NewRecorder()
  Router().ServeHTTP(response2, request2)
  assert.Equal(t, 400, response2.Code, "Bad request response is expected")

  // status 200
  create_test_user(t)
  request3, _ := http.NewRequest("POST", "/reserve_money/?user_id=9999&service_id=1&order_id=1&price=1000", nil)
  response3 := httptest.NewRecorder()
  Router().ServeHTTP(response3, request3)
  delete_test_reservation(t)
  delete_test_user(t)
  assert.Equal(t, 200, response3.Code, "OK response is expected")

  // not enough money
  create_test_user(t)
  request4, _ := http.NewRequest("POST", "/reserve_money/?user_id=9999&service_id=1&order_id=1&price=10000", nil)
  response4 := httptest.NewRecorder()
  Router().ServeHTTP(response4, request4)
  delete_test_user(t)
  assert.Equal(t, 400, response4.Code, "Bad request response is expected")

  // user not found
  create_test_user(t)
  request5, _ := http.NewRequest("POST", "/reserve_money/?user_id=0&service_id=1&order_id=1&price=1000", nil)
  response5 := httptest.NewRecorder()
  Router().ServeHTTP(response5, request5)
  delete_test_user(t)
  assert.Equal(t, 404, response5.Code, "Not found response is expected")
}

func TestRevenueRecognition(t *testing.T) {
  // bad request
  request1, _ := http.NewRequest("POST", "/revenue_recognition/", nil)
  response1 := httptest.NewRecorder()
  Router().ServeHTTP(response1, request1)
  assert.Equal(t, 400, response1.Code, "Bad request response is expected")

  // not enough params
  request2, _ := http.NewRequest("POST", "/revenue_recognition/?user_id=9999&service_id=1&order_id=1", nil)
  response2 := httptest.NewRecorder()
  Router().ServeHTTP(response2, request2)
  assert.Equal(t, 400, response2.Code, "Bad request response is expected")

  // status 200
  create_test_reservation(t)
  request3, _ := http.NewRequest("POST", "/revenue_recognition/?user_id=9999&service_id=1&order_id=1&price=1000", nil)
  response3 := httptest.NewRecorder()
  Router().ServeHTTP(response3, request3)
  delete_test_revenue_recognition(t)
  assert.Equal(t, 200, response3.Code, "OK response is expected")

  // reservation not found
  create_test_user(t)
  create_test_reservation(t)
  request5, _ := http.NewRequest("POST", "/revenue_recognition/?user_id=0&service_id=1&order_id=1&price=1000", nil)
  response5 := httptest.NewRecorder()
  Router().ServeHTTP(response5, request5)
  delete_test_reservation(t)
  delete_test_user(t)
  assert.Equal(t, 404, response5.Code, "Not found response is expected")
}

func TestUnreservingMoney(t *testing.T) {
  // bad request
  request1, _ := http.NewRequest("POST", "/unreserve_money/", nil)
  response1 := httptest.NewRecorder()
  Router().ServeHTTP(response1, request1)
  assert.Equal(t, 400, response1.Code, "Bad request response is expected")

  // not enough params
  request2, _ := http.NewRequest("POST", "/unreserve_money/?user_id=9999&service_id=1&order_id=1", nil)
  response2 := httptest.NewRecorder()
  Router().ServeHTTP(response2, request2)
  assert.Equal(t, 400, response2.Code, "Bad request response is expected")

  // status 200
  create_test_user(t)
  create_test_reservation(t)
  request3, _ := http.NewRequest("POST", "/unreserve_money/?user_id=9999&service_id=1&order_id=1&price=1000", nil)
  response3 := httptest.NewRecorder()
  Router().ServeHTTP(response3, request3)
  delete_test_user(t)
  assert.Equal(t, 200, response3.Code, "OK response is expected")

  // reservation not found
  create_test_user(t)
  create_test_reservation(t)
  request5, _ := http.NewRequest("POST", "/unreserve_money/?user_id=0&service_id=1&order_id=1&price=1000", nil)
  response5 := httptest.NewRecorder()
  Router().ServeHTTP(response5, request5)
  delete_test_reservation(t)
  delete_test_user(t)
  assert.Equal(t, 404, response5.Code, "Not found response is expected")
}
