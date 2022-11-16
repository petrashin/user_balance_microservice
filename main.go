package main

import (
  "fmt"
  "log"
  "strconv"
  "net/http"
  "encoding/json"
  "github.com/gorilla/mux"

  "database/sql"
  _ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type User struct {
  Id int `json:"id"`
  Balance float64 `json:"balance"`
}

func get_user_balance(id int) float64 {
  query := "SELECT * FROM `users` WHERE `id`=?"
  rows, err := db.Query(query, id)
  if err != nil {
    panic(err)
  }

  var user = User{}
  for rows.Next() {
    var tmp_user User
    err = rows.Scan(&tmp_user.Id, &tmp_user.Balance)
    if err != nil {
      panic(err)
    }
    user = tmp_user
  }
  return user.Balance
}

func check_user_exists(id int) bool {
  query := "SELECT * FROM `users` WHERE `id`=?"
  rows, err := db.Query(query, id)
  if err != nil {
    panic(err)
  }
  i := 0
  for rows.Next() {
    i++
  }
  if i > 0 {
    return true
  }
  return false
}

func check_reservation_exists(user_id int, service_id int, order_id int, price float64) bool {
  query := "SELECT * FROM `reservations` WHERE `user`=? AND `service`=? AND `order_id`=? AND `price`=?"
  rows, err := db.Query(query, user_id, service_id, order_id, price)
  if err != nil {
    panic(err)
  }
  i := 0
  for rows.Next() {
    i++
  }
  if i > 0 {
    return true
  }
  return false
}

func update_balance(w http.ResponseWriter, r *http.Request)  {
  w.Header().Set("Content-Type", "application/json")

  id, err1 := strconv.Atoi(r.URL.Query().Get("id"))
  balance, err2 := strconv.ParseFloat(r.URL.Query().Get("balance"), 64)

  if err1 != nil || err2 != nil {
    w.WriteHeader(http.StatusBadRequest)
  	json.NewEncoder(w).Encode(map[string]string{"result": "Bad request"})
  	return
  }

  query := "INSERT INTO `users`(`id`, `balance`) VALUES (?, ?) ON DUPLICATE KEY UPDATE `balance` = `balance` + ?"
  insert, err := db.Query(query, id, balance, balance)
  if err != nil {
      panic(err)
  }
  defer insert.Close()

  w.WriteHeader(http.StatusOK)
  json.NewEncoder(w).Encode(map[string]string{"result": "Balance updated"})
}

func reserve_money(w http.ResponseWriter, r *http.Request)  {
  w.Header().Set("Content-Type", "application/json")

  user_id, err1 := strconv.Atoi(r.URL.Query().Get("user_id"))
  service_id, err2 := strconv.Atoi(r.URL.Query().Get("service_id"))
  order_id, err3 := strconv.Atoi(r.URL.Query().Get("order_id"))
  price, err4 := strconv.ParseFloat(r.URL.Query().Get("price"), 64)

  if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
    w.WriteHeader(http.StatusBadRequest)
  	json.NewEncoder(w).Encode(map[string]string{"result": "Bad request"})
  	return
  }

  if check_user_exists(user_id) {
    new_balance := get_user_balance(user_id) - price

    if new_balance >= 0 {
      query := "INSERT INTO `reservations`(`user`, `service`, `order_id`, `price`) VALUES (?, ?, ?, ?)"
      insert, err := db.Query(query, user_id, service_id, order_id, price)
      if err != nil {
        panic(err)
      }
      defer insert.Close()

      query = "UPDATE users SET `balance` = ? WHERE `id` = ?"
      update, err := db.Query(query, new_balance, user_id)
      if err != nil {
        panic(err)
      }
      defer update.Close()

      w.WriteHeader(http.StatusOK)
      json.NewEncoder(w).Encode(map[string]string{"result": "Money reserved"})
      return
    } else {
      w.WriteHeader(http.StatusBadRequest)
      json.NewEncoder(w).Encode(map[string]string{"result": "Not enough money"})
      return
    }
  } else {
    w.WriteHeader(http.StatusNotFound)
    json.NewEncoder(w).Encode(map[string]string{"result": "User not found"})
    return
  }
}

func revenue_recognition(w http.ResponseWriter, r *http.Request)  {
  w.Header().Set("Content-Type", "application/json")

  user_id, err1 := strconv.Atoi(r.URL.Query().Get("user_id"))
  service_id, err2 := strconv.Atoi(r.URL.Query().Get("service_id"))
  order_id, err3 := strconv.Atoi(r.URL.Query().Get("order_id"))
  price, err4 := strconv.ParseFloat(r.URL.Query().Get("price"), 64)

  if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
    w.WriteHeader(http.StatusBadRequest)
  	json.NewEncoder(w).Encode(map[string]string{"result": "Bad request"})
  	return
  }

  if check_reservation_exists(user_id, service_id, order_id, price) {
    query := "INSERT INTO `revenue` SELECT * FROM `reservations` WHERE user = ? AND service = ? AND order_id = ? AND price = ?"
    insert, err := db.Query(query, user_id, service_id, order_id, price)
    if err != nil {
      panic(err)
    }
    defer insert.Close()

    query = "DELETE FROM `reservations` WHERE user = ? AND service = ? AND order_id = ? AND price = ?"
    delete, err := db.Query(query, user_id, service_id, order_id, price)
    if err != nil {
      panic(err)
    }
    defer delete.Close()

    json.NewEncoder(w).Encode(map[string]string{"result": "Successfully"})
    return
  } else {
    w.WriteHeader(http.StatusNotFound)
  	json.NewEncoder(w).Encode(map[string]string{"result": "Not found"})
  	return
  }
}

func get_balance(w http.ResponseWriter, r *http.Request)  {
  w.Header().Set("Content-Type", "application/json")
  user_id, err := strconv.Atoi(r.URL.Query().Get("id"))

  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
  	json.NewEncoder(w).Encode(map[string]string{"result": "Bad request"})
  	return
  }
  user_exists := check_user_exists(user_id)
  if user_exists {
    balance := get_user_balance(user_id)
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]float64{"balance": balance})
    return
  } else {
    w.WriteHeader(http.StatusNotFound)
    json.NewEncoder(w).Encode(map[string]string{"result": "Not found"})
    return
  }
}

func unreserve_money(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")

  user_id, err1 := strconv.Atoi(r.URL.Query().Get("user_id"))
  service_id, err2 := strconv.Atoi(r.URL.Query().Get("service_id"))
  order_id, err3 := strconv.Atoi(r.URL.Query().Get("order_id"))
  price, err4 := strconv.ParseFloat(r.URL.Query().Get("price"), 64)

  if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
    w.WriteHeader(http.StatusBadRequest)
  	json.NewEncoder(w).Encode(map[string]string{"result": "Bad request"})
  	return
  }

  if check_reservation_exists(user_id, service_id, order_id, price) {
    new_balance := get_user_balance(user_id) + price

    query := "DELETE FROM `reservations` WHERE user = ? AND service = ? AND order_id = ? AND price = ?"
    delete, err := db.Query(query, user_id, service_id, order_id, price)
    if err != nil {
      panic(err)
    }
    defer delete.Close()

    query = "UPDATE users SET `balance` = ? WHERE `id` = ?"
    update, err := db.Query(query, new_balance, user_id)
    if err != nil {
      panic(err)
    }
    defer update.Close()

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"result": "Money unreserved"})
    return
  } else {
    w.WriteHeader(http.StatusNotFound)
    json.NewEncoder(w).Encode(map[string]string{"result": "Reservation not found"})
    return
  }
}

func handleFunc()  {
  rtr := mux.NewRouter()

  rtr.HandleFunc("/update_balance/", update_balance).Methods("POST")
  rtr.HandleFunc("/reserve_money/", reserve_money).Methods("POST")
  rtr.HandleFunc("/revenue_recognition/", revenue_recognition).Methods("POST")
  rtr.HandleFunc("/get_balance/", get_balance).Methods("GET")
  rtr.HandleFunc("/unreserve_money/", unreserve_money).Methods("POST")

  http.Handle("/", rtr)
  http.ListenAndServe(":8080", nil)
}

func init_db(db_type string, username string, password string, port string, database_name string) error {
    var err error
    connectionString := fmt.Sprintf("%s:%s@tcp(127.0.0.1:%s)/%s", username, password, port, database_name)

    db, err = sql.Open(db_type, connectionString)
    if err != nil {
        return err
    }

    return db.Ping()
}

func main()  {
  err := init_db("mysql", "root", "root", "3306", "avito")

  if err != nil {
    log.Fatal(err)
  }

  handleFunc()
}
