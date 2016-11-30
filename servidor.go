package main
import (
  "encoding/json"
  "fmt"
  "log"
  "net/http"
  "gorilla/mux"
  "strings"
)

type Torrent struct {
  Title string
  Description string
  MagnetLink string
  Size string
  Downloads int
  Seeders int
  Leechers int
}

var torrents = map[string]*Torrent{
  "Filme1": &Torrent{Title: "Filme1",Description:"Teste",MagnetLink:"1",Size:"1.2GB",Downloads:20,Seeders:1,Leechers:0},
  "Filme2": &Torrent{Title: "Filme2",Description:"Filme muito legal.",MagnetLink:"1245464",Size:"500MB",Downloads:300,Seeders:13,Leechers:2},
  "Jogo1": &Torrent{Title: "Jogo1",Description:"Um jogo de ação.",MagnetLink:"767834",Size:"5GB",Downloads:20,Seeders:78,Leechers:91},
}

func handleTorrents(res http.ResponseWriter, req *http.Request) {
  res.Header().Set("Content-Type", "application/json")
  outgoingJSON, error := json.Marshal(torrents)
  if error != nil {
    log.Println(error.Error())
    http.Error(res, error.Error(),
    http.StatusInternalServerError)
    return
  }
  fmt.Fprint(res, string(outgoingJSON))
}

func handleTorrent(res http.ResponseWriter, req *http.Request) {
  res.Header().Set("Content-Type", "application/json")
  vars := mux.Vars(req)
  Key := strings.TrimSuffix(vars["Key"], "\n")
  log.Println("Request for:", Key)

  switch req.Method {
    case "GET":
      torrent, ok := torrents[Key]
      if !ok {
        res.WriteHeader(http.StatusNotFound)
        fmt.Fprint(res, string("Torrent not found"))
        return
      }
      outgoingJSON, error := json.Marshal(torrent)
      if error != nil {
        log.Println(error.Error())
        http.Error(res, error.Error(), http.StatusInternalServerError)
        return
      }
      fmt.Fprint(res, string(outgoingJSON))
      break
    case "DELETE":
      _, ok := torrents[Key]
      if ok {
        delete(torrents, Key)
        res.WriteHeader(http.StatusOK)
        break
      }
      res.WriteHeader(http.StatusNoContent)
      break
    case "POST":
      torrent := new(Torrent)
      decoder := json.NewDecoder(req.Body)
      error := decoder.Decode(&torrent)
      if error != nil {
        fmt.Println(error.Error())
        log.Println(error.Error())
        http.Error(res, error.Error(), http.StatusInternalServerError)
        return
      }
      torrents[Key] = torrent
      outgoingJSON, err := json.Marshal(torrent)
      if err != nil {
        log.Println(error.Error())
        http.Error(res, err.Error(), http.StatusInternalServerError)
        return
      }
      res.WriteHeader(http.StatusCreated)
      fmt.Fprint(res, string(outgoingJSON))

    case "PUT":
      torrent, ok := torrents[Key]
      if !ok {
        res.WriteHeader(http.StatusNotFound)
        fmt.Fprint(res, string("Torrent not found"))
        return
      }
      update := new(Torrent)
      decoder := json.NewDecoder(req.Body)
      error := decoder.Decode(&update)
      if error != nil {
        log.Println(error.Error())
        http.Error(res, error.Error(), http.StatusInternalServerError)
        return
      }
      if update.Title != "" {
        torrent.Title = update.Title
      }
      if update.Description != "" {
        torrent.Description = update.Description
      }
      if update.MagnetLink != "" {
        torrent.MagnetLink = update.MagnetLink
      }
      if update.Size != "" {
        torrent.Size = update.Size
      }
      if update.Downloads != 0 {
        torrent.Downloads = update.Downloads
      }
      if update.Seeders != 0 {
        torrent.Seeders = update.Seeders
      }
      if update.Leechers != 0 {
        torrent.Leechers = update.Leechers
      }
      res.WriteHeader(http.StatusNoContent)
    }
}

func main () {
  router := mux.NewRouter().StrictSlash(true)
  router.HandleFunc("/torrents", handleTorrents).Methods("GET")
  router.HandleFunc("/torrent/{Key}", handleTorrent).Methods("GET", "DELETE","POST","PUT")
  log.Fatal(http.ListenAndServe("localhost:8080", router))
}
