package main
import (
  "encoding/json"
  "fmt"
  "log"
  "net/http"
  "gorilla/mux"
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
  "1": &Torrent{Title: "Filme1",Description:"Teste",MagnetLink:"1",Size:"1.2GB",Downloads:20,Seeders:1,Leechers:0},
}

//Não mexi em nada nessas funções handle ainda
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
  Key := vars["Key"]
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
    case "DELETE":
      delete(torrents, Key)
      res.WriteHeader(http.StatusNoContent)
  }
}

func main () {
  router := mux.NewRouter().StrictSlash(true)
  router.HandleFunc("/torrents", handleTorrents).Methods("GET")
  router.HandleFunc("/torrent/{Key}", handleTorrent).Methods("GET", "DELETE") //Sei la que Key seria (magnetlink talvez?), qualquer coisa tirar essa linha
  log.Fatal(http.ListenAndServe("localhost:8080", router))
}
