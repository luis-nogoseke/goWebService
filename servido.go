package main
import (
    "fmt"
    "log"
    "net/http"
    "gorilla/mux"
)

type Torrents struct {
  Title string
  Description string
  MagnetLink string
  Size string
  Downloads int
  Seeders int
  Leechers int
}

/*var movies = map[string]*Movie{
  "tt0246578": &Movie{Title: "Donnie Darko",
  Rating: 8.1, Year: 2001},
  "tt0080120": &Movie{Title: "The Warriors",
  Rating: 7.7, Year: 1979},
  "tt0074486": &Movie{Title: "Eraserhead",
  Rating: 7.4, Year: 1977},
}*/

var torrents = map[string]*Torrent{

}


//Não mexi em nada nessas funções handle ainda
func handleTorrents(res http.ResponseWriter, req *http.Request) {
  res.Header().Set("Content-Type", "application/json")
  outgoingJSON, error := json.Marshal(movies)
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
  imdbKey := vars["imdbKey"]
  log.Println("Request for:", imdbKey)
  movie, ok := movies[imdbKey]
  if !ok {
    res.WriteHeader(http.StatusNotFound)
    fmt.Fprint(res, string("Movie not found"))
    return
  }
  outgoingJSON, error := json.Marshal(movie)
  if error != nil {
    log.Println(error.Error())
    http.Error(res, error.Error(), http.StatusInternalServerError)
    return
  }
  fmt.Fprint(res, string(outgoingJSON))
}

func main () {
  router := mux.NewRouter().StrictSlash(true)
  router.HandleFunc("/torrents", handleTorrents).Methods("GET")
  router.HandleFunc("/torrent/{Key}", handleTorrent).Methods("GET") //Sei la que Key seria (magnetlink talvez?), qualquer coisa tirar essa linha
  handleMovie).Methods("GET")  log.Fatal(http.ListenAndServe("localhost:8080", router))
}
