package main

import (
  "fmt"
  "strconv"
  "net/http"
  "encoding/json"
  "github.com/gorilla/mux"

  "database/sql"
  _ "github.com/go-sql-driver/mysql"
)

const db_type = "mysql"
const connection = "root:root@tcp(127.0.0.1:3306)/avito"

type User struct {
  Id int `json:"id"`
  Balance float64 `json:"balance"`
}

type Service struct {
  Id int `json:"id"`
}

type Order struct {
  Id int `json:"id"`
}

type Reservation struct {
  Id int `json:"id"`
  User *User `json:"user"`
  Service *Service `json:"service"`
  Order *Order `json:"order"`
  price float64 `json:"price"`
}

// TODO: get_response(message) -> jsonResp

func update_balance(w http.ResponseWriter, r *http.Request)  {
  w.Header().Set("Content-Type", "application/json")

  id := r.URL.Query().Get("id")
  balance := r.URL.Query().Get("balance")

  if id == "" || balance == "" {

    w.WriteHeader(http.StatusBadRequest)
    // TODO change for get_response
  	resp := make(map[string]string)
  	resp["message"] = "Bad Request"
  	jsonResp, err := json.Marshal(resp)
  	if err != nil {
  		panic(err)
  	}
  	w.Write(jsonResp)
  	return

  } else {

    id, err := strconv.Atoi(id)
    if err != nil {
      panic(err)
    }

    balance, err := strconv.ParseFloat(balance, 64)
    if err != nil {
      panic(err)
    }

    db, err := sql.Open(db_type, connection)
    if err != nil {
      panic(err)
    }
    defer db.Close()

    upsert_query := fmt.Sprintf("INSERT INTO `users`(`id`, `balance`) VALUES ('%d', '%g')" +
                                "ON DUPLICATE KEY UPDATE `balance` = `balance` + '%g';", id, balance, balance)

    insert, err := db.Query(upsert_query)
      if err != nil {
        panic(err)
      }
    defer insert.Close()

    get_query := fmt.Sprintf("SELECT * FROM `users` WHERE `id`='%d'", id)

    get, err := db.Query(get_query)
    if err != nil {
      panic(err)
    }

    var result_user = User{}

    for get.Next() {
      var user User
      err = get.Scan(&user.Id, &user.Balance)
      if err != nil {
        panic(err)
      }
      result_user = user
    }

    json.NewEncoder(w).Encode(result_user)
  }
}

func reserve_money(w http.ResponseWriter, r *http.Request)  {
  w.Header().Set("Content-Type", "application/json")

  user_id := r.URL.Query().Get("user_id")
  service_id := r.URL.Query().Get("service_id")
  order_id := r.URL.Query().Get("order_id")
  price := r.URL.Query().Get("price")

  if user_id == "" || service_id == "" || order_id == "" || price == "" {

    w.WriteHeader(http.StatusBadRequest)
    // TODO change for get_response
  	resp := make(map[string]string)
  	resp["message"] = "Bad Request"
  	jsonResp, err := json.Marshal(resp)
  	if err != nil {
  		panic(err)
  	}
  	w.Write(jsonResp)
  	return
  } else {

    resp := make(map[string]string)
  	resp["message"] = "Money reserved"
  	jsonResp, err := json.Marshal(resp)
  	if err != nil {
  		panic(err)
  	}
  	w.Write(jsonResp)
  	return

  }
}

func handleFunc()  {
  rtr := mux.NewRouter()

  rtr.HandleFunc("/update_balance/", update_balance).Methods("POST")
  rtr.HandleFunc("/reserve_money/", reserve_money).Methods("POST")

  http.Handle("/", rtr)
  http.ListenAndServe(":8080", nil)
}

func main()  {
  handleFunc()
}
