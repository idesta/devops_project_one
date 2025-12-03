// api/main.go
package main

import (
  "database/sql"
  "encoding/json"
  "log"
  "net/http"
  "os"

  _ "github.com/lib/pq"
)

var db *sql.DB

func main() {
  dsn := os.Getenv("DATABASE_DSN") // e.g. "postgres://user:pass@postgres:5432/appdb?sslmode=disable"
  var err error
  db, err = sql.Open("postgres", dsn)
  if err != nil { log.Fatal(err) }
  defer db.Close()

  http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(200); w.Write([]byte("ok"))
  })

  http.HandleFunc("/register", registerHandler)
  http.HandleFunc("/login", loginHandler)

  port := os.Getenv("PORT")
  if port == "" { port = "8080" }
  log.Println("api listening on", port)
  log.Fatal(http.ListenAndServe(":"+port, nil))
}

type User struct {
  Username string `json:"username"`
  Password string `json:"password"`
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodPost { http.Error(w, "method", 405); return}
  var u User; json.NewDecoder(r.Body).Decode(&u)
  _, err := db.Exec("INSERT INTO users(username,password) VALUES($1, $2)", u.Username, u.Password)
  if err != nil { http.Error(w, err.Error(), 500); return }
  w.WriteHeader(201)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodPost { http.Error(w, "method", 405); return}
  var u User; json.NewDecoder(r.Body).Decode(&u)
  var pw string
  err := db.QueryRow("SELECT password FROM users WHERE username=$1", u.Username).Scan(&pw)
  if err != nil { http.Error(w, "invalid", 401); return}
  if pw != u.Password { http.Error(w, "invalid", 401); return}
  w.Write([]byte(`{"ok":true}`))
}
