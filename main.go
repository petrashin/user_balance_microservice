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

type Reservation struct {
  User User `json:"user"`
  Service int `json:"service"`
  Order int `json:"order"`
  Price float64 `json:"price"`
}

type Revenue struct {
  User User `json:"user"`
  Service int `json:"service"`
  Order int `json:"order"`
  Price float64 `json:"price"`
}

func update_balance(w http.ResponseWriter, r *http.Request)  {
  w.Header().Set("Content-Type", "application/json")

  id := r.URL.Query().Get("id")
  balance := r.URL.Query().Get("balance")

  if id == "" || balance == "" {
    w.WriteHeader(http.StatusBadRequest)
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
  	resp := make(map[string]string)
  	resp["message"] = "Bad Request"
  	jsonResp, err := json.Marshal(resp)
  	if err != nil {
  		panic(err)
  	}
  	w.Write(jsonResp)
  	return
  } else {
    user_id, err := strconv.Atoi(user_id)
    if err != nil {
      panic(err)
    }

    service_id, err := strconv.Atoi(service_id)
    if err != nil {
      panic(err)
    }

    order_id, err := strconv.Atoi(order_id)
    if err != nil {
      panic(err)
    }

    db, err := sql.Open(db_type, connection)
    if err != nil {
      panic(err)
    }
    defer db.Close()

    get_query := fmt.Sprintf("SELECT * FROM `users` WHERE `id`='%d'", user_id)

    get_user, err := db.Query(get_query)
    if err != nil {
      panic(err)
    }

    var result_user = User{}

    for get_user.Next() {
      var user User
      err = get_user.Scan(&user.Id, &user.Balance)
      if err != nil {
        panic(err)
      }
      result_user = user
    }

    price, err := strconv.ParseFloat(price, 64)
    if err != nil {
      panic(err)
    }

    new_balance := result_user.Balance - price

    if new_balance >= 0 {
      insert_query := fmt.Sprintf("INSERT INTO `reservations`(`user`, `service`, `order_id`, `price`)" +
                                  "VALUES ('%d', '%d', '%d', '%g')", user_id, service_id, order_id, price)

      insert, err := db.Query(insert_query)
        if err != nil {
          panic(err)
        }
      defer insert.Close()

      update_query := fmt.Sprintf("UPDATE users " +
                                  "SET `balance` = '%g' " +
                                  "WHERE `id` = '%d';", new_balance, user_id)

      result_user.Balance = new_balance
      var reservation = Reservation{User: result_user, Service: service_id, Order: order_id, Price: price}

      update, err := db.Query(update_query)
        if err != nil {
          panic(err)
        }
      defer update.Close()

      json.NewEncoder(w).Encode(reservation)

    } else {
      w.WriteHeader(http.StatusBadRequest)
    	resp := make(map[string]string)
    	resp["message"] = "Not enough money"
    	jsonResp, err := json.Marshal(resp)
    	if err != nil {
    		panic(err)
    	}
    	w.Write(jsonResp)
    	return
    }
  }
}

func revenue_recognition(w http.ResponseWriter, r *http.Request)  {
  w.Header().Set("Content-Type", "application/json")

  user_id := r.URL.Query().Get("user_id")
  service_id := r.URL.Query().Get("service_id")
  order_id := r.URL.Query().Get("order_id")
  price := r.URL.Query().Get("price")

  if user_id == "" || service_id == "" || order_id == "" || price == "" {
    w.WriteHeader(http.StatusBadRequest)
  	resp := make(map[string]string)
  	resp["message"] = "Bad Request"
  	jsonResp, err := json.Marshal(resp)
  	if err != nil {
  		panic(err)
  	}
  	w.Write(jsonResp)
  	return
  } else {
    user_id, err := strconv.Atoi(user_id)
    if err != nil {
      panic(err)
    }
    service_id, err := strconv.Atoi(service_id)
    if err != nil {
      panic(err)
    }
    order_id, err := strconv.Atoi(order_id)
    if err != nil {
      panic(err)
    }
    price, err := strconv.ParseFloat(price, 64)
    if err != nil {
      panic(err)
    }

    db, err := sql.Open(db_type, connection)
    if err != nil {
      panic(err)
    }
    defer db.Close()

    insert_query := fmt.Sprintf("INSERT INTO `revenue` " +
                                "SELECT * FROM `reservations` " +
                                "WHERE user = '%d' " +
                                "AND service = '%d' " +
                                "AND order_id = '%d' " +
                                "AND price = '%g';", user_id, service_id, order_id, price)

    delete_query := fmt.Sprintf("DELETE FROM `reservations` " +
                                "WHERE user = '%d' " +
                                "AND service = '%d' " +
                                "AND order_id = '%d' " +
                                "AND price = '%g';", user_id, service_id, order_id, price)

    insert, err := db.Query(insert_query)
    if err != nil {
      panic(err)
    }
    defer insert.Close()

    delete, err := db.Query(delete_query)
    if err != nil {
      panic(err)
    }
    defer delete.Close()

    get_query := fmt.Sprintf("SELECT * FROM `users` WHERE `id`='%d'", user_id)

    get_user, err := db.Query(get_query)
    if err != nil {
      panic(err)
    }

    var result_user = User{}

    for get_user.Next() {
      var user User
      err = get_user.Scan(&user.Id, &user.Balance)
      if err != nil {
        panic(err)
      }
      result_user = user
    }

    var revenue = Revenue{User: result_user, Service: service_id, Order: order_id, Price: price}
    json.NewEncoder(w).Encode(revenue)
  }
}

func get_balance(w http.ResponseWriter, r *http.Request)  {
  w.Header().Set("Content-Type", "application/json")

  id := r.URL.Query().Get("id")

  if id == "" {
    w.WriteHeader(http.StatusBadRequest)
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

    db, err := sql.Open(db_type, connection)
    if err != nil {
      panic(err)
    }
    defer db.Close()

    get_query := fmt.Sprintf("SELECT * FROM `users` WHERE `id`='%d'", id)

    get, err := db.Query(get_query)
    if err != nil {
      panic(err)
    }

    var result_user = User{}

    i := 0
    for get.Next() {
      var user User
      err = get.Scan(&user.Id, &user.Balance)
      i++
      if err != nil {
        panic(err)
      }
      result_user = user
    }

    if i == 0 {
      w.WriteHeader(http.StatusOK)
    	resp := make(map[string]string)
    	resp["message"] = "No such user"
    	jsonResp, err := json.Marshal(resp)
    	if err != nil {
    		panic(err)
    	}
    	w.Write(jsonResp)
    	return
    } else {
      w.WriteHeader(http.StatusOK)
    	resp := make(map[string]float64)
    	resp["balance"] = result_user.Balance
    	jsonResp, err := json.Marshal(resp)
    	if err != nil {
    		panic(err)
    	}
    	w.Write(jsonResp)
    	return
    }
  }
}

func handleFunc()  {
  rtr := mux.NewRouter()

  rtr.HandleFunc("/update_balance/", update_balance).Methods("POST")
  rtr.HandleFunc("/reserve_money/", reserve_money).Methods("POST")
  rtr.HandleFunc("/revenue_recognition/", revenue_recognition).Methods("POST")
  rtr.HandleFunc("/get_balance/", get_balance).Methods("GET")

  http.Handle("/", rtr)
  http.ListenAndServe(":8080", nil)
}

func main()  {
  handleFunc()
}
