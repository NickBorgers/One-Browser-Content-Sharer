package main

import (
  "log"
  "io/ioutil"
  "fmt"
  "bytes"
  "net/http"
  "time"
  "github.com/google/uuid"
)

func main() {
  http.HandleFunc("/", handler)

  log.Println("Listening...")
  http.ListenAndServe(":3000", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
  dat, err := ioutil.ReadFile("lock.cookie")
  if err == nil {
    fmt.Println("Found existing lock on content")
    cookie, err := r.Cookie("lock")
    if err != nil {
      fmt.Println("Did not find lock cookie on request, definitely isn't authorized")
      w.WriteHeader(http.StatusUnauthorized)
      return
    }

    comparison := bytes.Compare([]byte(cookie.Value), dat)

    if comparison != 0 {
      fmt.Println("Presented lock cookie does not match known lock")
      w.WriteHeader(http.StatusUnauthorized)
      return
    }
  } else {
    fmt.Println("Did not find existing lock cookie value, generating and setting for requester")
    expire := time.Now().AddDate(10, 0, 0)
    secret := uuid.New()
    cookie := http.Cookie{
        Name:    "lock",
        Value:   secret.String(),
        Expires: expire,
        Secure: true,
        HttpOnly: true,
    }
    http.SetCookie(w, &cookie)
    ioutil.WriteFile("lock.cookie", []byte(secret.String()), 0644)
  }

  fmt.Println("Requester is authorized to consume")
  data, err := ioutil.ReadFile("static/content.html")
  if err != nil { fmt.Fprint(w, err) }
  http.ServeContent(w, r, "", time.Now(), bytes.NewReader(data))
}
