package main

import (
  "fmt"
  "net/http"
  "log"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintln(w, "home")
}

func main() {
  http.HandleFunc("/", homeHandler)
  log.Fatal(http.ListenAndServe(":8181", nil))
}
