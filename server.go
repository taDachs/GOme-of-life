package gameoflife

import (
  "context"
  "encoding/json"
  "fmt"

  "log"
  "net/http"
  "os"
  "os/signal"
  "syscall"
  "time"

  "github.com/gorilla/mux"
)

func RunServer(port string, game *Game) {
  var wait time.Duration

  r := mux.NewRouter()
  api := r.PathPrefix("/api").Subrouter()

  if game.IsHost {
    api.HandleFunc("/init", func(w http.ResponseWriter, r *http.Request) {
      initHandle(game, w, r)
    })
  }

  api.HandleFunc("/sync", syncHandle)

  api.HandleFunc("/update", updateHandle).Methods(http.MethodPost)

  srv := &http.Server{
    Addr:         ":" + port,
    WriteTimeout: time.Second * 15,
    ReadTimeout:  time.Second * 15,
    IdleTimeout:  time.Second * 60,
    Handler:      r,
  }

  go func() {
    if err := srv.ListenAndServe(); err != nil {
      fmt.Println(err)
    }
  }()

  end := make(chan os.Signal, 1)
  signal.Notify(end, os.Interrupt, syscall.SIGTERM)
  <-end

  ctx, cancel := context.WithTimeout(context.Background(), wait)
  defer cancel()

  srv.Shutdown(ctx)
  os.Exit(0)
}


func syncHandle(w http.ResponseWriter, r *http.Request) {
  var resp Sync

  err := json.NewDecoder(r.Body).Decode(&resp)
  if err != nil {
    HttpErrWrite("Error while reading body", err, http.StatusUnprocessableEntity, w)
    return
  }

  SyncChannel <- resp

  w.Header().Set("content-type", "application/json")
  w.WriteHeader(http.StatusOK)
}

func HttpErrWrite(msg string, err error, status int, w http.ResponseWriter) {
  log.Println(msg, err)

  w.Header().Set("content-type", "application/json")
  w.WriteHeader(status)
  if err := json.NewEncoder(w).Encode(&HttpErrorBody{Err: msg + err.Error()}); err != nil {
    fmt.Println(err)
  }
  return
}

func updateHandle(w http.ResponseWriter, r *http.Request) {
  var resp []Change

  err := json.NewDecoder(r.Body).Decode(&resp)
  if err != nil {
    HttpErrWrite("Error while reading body", err, http.StatusUnprocessableEntity, w)
    return
  }

  for _, chg := range resp {
    ChangeChannel <- chg
  }

  w.Header().Set("content-type", "application/json")
  w.WriteHeader(http.StatusOK)
}

func initHandle(game *Game, w http.ResponseWriter, r *http.Request) {
  w.Header().Set("content-type", "application/json")
  if game.Started  {
    json, err := json.Marshal("Game already started")

    if err != nil {
      fmt.Println("Error while converting init to json: ", err)
    }

    w.WriteHeader(http.StatusInternalServerError)
    w.Write(json)
    return
  }
  var init Init
  init.Board = *game.Board

  json, err := json.Marshal(init)

  if err != nil {
    fmt.Println("Error while converting init to json: ", err)
    w.WriteHeader(http.StatusInternalServerError)
    w.Write(json)
    return
  }

  w.WriteHeader(http.StatusOK)
  w.Write(json)
  game.Started = true
}

type HttpErrorBody struct {
  Err string `json:"error"`
}
