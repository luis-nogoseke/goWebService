package main
import (
  "encoding/json"
  "fmt"
  "log"
 // "net/url"
  "os"
  "strings"
  "net/http"
  "bufio"
  "strconv"
  "bytes"
)

const serverURL = "http://localhost:8080/"

var in *bufio.Reader


type Torrent struct {
  Title string
  Description string
  MagnetLink string
  Size string
  Downloads int
  Seeders int
  Leechers int
}

func printTorrent (t *Torrent) {
    fmt.Println("Title: ", t.Title)
    fmt.Println("Description: ", t.Description)
    fmt.Println("Size: ", t.Size)
    fmt.Println("Total number of downloads: ", t.Downloads)
    fmt.Println("Seeders: ", t.Seeders)
    fmt.Println("Leechers: ", t.Leechers)
    fmt.Println("Magnet Link: ", t.MagnetLink)
    fmt.Printf("\n\n")
}


func get() {
    resp, err := http.Get(serverURL+"torrents")
    if err != nil {
        log.Fatalf("could not fetch: %v", err)
    }  
    fmt.Println(resp)

    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        fmt.Println("Erro 1")
        return
    }
    var result map[string]*Torrent
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        fmt.Println("Erro 2")
        return
    }
    for _, tor := range result {
        printTorrent(tor)
    }
}

func getWithKey(key string) {
    resp, err := http.Get(serverURL+"torrent/"+strings.TrimSuffix(key, "\n"))
    if err != nil {
        log.Fatalf("Falha: %v", err)
    }  

    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        fmt.Println("Torrent não encontrado")
        return
    }
    var result Torrent
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        fmt.Println("Falha de decodificação")
        return
    }
    printTorrent(&result)  
}

func post() {
    var postT Torrent
    fmt.Println("Nome:")
    aux, _ := in.ReadString('\n')
    postT.Title = strings.TrimSuffix(aux, "\n")
    fmt.Println("Descrição:")
    aux, _ = in.ReadString('\n')
    postT.Description = strings.TrimSuffix(aux, "\n")
    fmt.Println("Magnet Link:")
    aux, _ = in.ReadString('\n')
    postT.MagnetLink = strings.TrimSuffix(aux, "\n")
    fmt.Println("Tamanho:")
    aux, _ = in.ReadString('\n')
    postT.Size = strings.TrimSuffix(aux, "\n")
    fmt.Println("Total de Downloads:")
    aux, _ = in.ReadString('\n')
    postT.Downloads, _ = strconv.Atoi(strings.TrimSuffix(aux, "\n"))
    fmt.Println("Seeders:")
    aux, _ = in.ReadString('\n')
    postT.Seeders, _ = strconv.Atoi(strings.TrimSuffix(aux, "\n"))
    fmt.Println("Leechers:")
    aux, _ = in.ReadString('\n')
    postT.Leechers, _ = strconv.Atoi(strings.TrimSuffix(aux, "\n"))

    postJSON, _ := json.Marshal(postT)
    resp, err := http.Post(serverURL+"torrent/"+postT.Title, "application/json", bytes.NewBuffer(postJSON))    

    if err != nil {
        log.Fatalf("Erro na operação de POST: %v", err)
    }

    defer resp.Body.Close()

    if resp.StatusCode != http.StatusCreated {
        fmt.Println("Erro ao criar torrent: ", resp.StatusCode)
        return
    }
}

func put() {
    var putT Torrent

    fmt.Println("Torrent para atualizar:")
    aux, _ := in.ReadString('\n')
    key := strings.TrimSuffix(aux, "\n")
    fmt.Println("Nova descrição:")
    aux, _ = in.ReadString('\n')
    putT.Description = strings.TrimSuffix(aux, "\n")
    putJSON, _ := json.Marshal(putT)

    req, err := http.NewRequest("PUT", serverURL+"torrent/"+key, bytes.NewBuffer(putJSON))
    if err != nil {
        return
    }
    resp, err := http.DefaultClient.Do(req)

    if err != nil {
        log.Fatalf("could not fetch: %v", err)
    }

    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
        fmt.Println("Torrent não encontrado.")
        return
    }
}

func delete(key string) {
    req, err := http.NewRequest("DELETE", serverURL+"torrent/"+key, nil)
    if err != nil {
        return
    }
    resp, err := http.DefaultClient.Do(req)

    if err != nil {
        log.Fatalf("could not fetch: %v", err)
    }

    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        fmt.Println("Torrent não encontrado.")
        return
    }
}

func main () {
    in = bufio.NewReader(os.Stdin)
    var op byte
    for {
        fmt.Println("\n MENU")
        fmt.Println("1 - Todos os torrents")
        fmt.Println("2 - Procurar torrent por nome")
        fmt.Println("3 - Post torrent")
        fmt.Println("4 - Delete torrent")
        fmt.Println("5 - Atualizar descrição do torrent")
        fmt.Println("6 - Sair")
        fmt.Scanf("%c\n", &op)
        switch op {
        case '1':
            get()
            break
        case '2':
            fmt.Println("Nome do torrent:")
            key, _ := in.ReadString('\n')
            getWithKey(key)
            break
        case '3':
            fmt.Println("Criar novo torrent")
            post()
            break
        case '4':
            fmt.Println("Nome do torrent para deletar:")
            key, _ := in.ReadString('\n')
            delete(key)
            break
        case '5':
            put()
            break
        case '6':
            os.Exit(0)
        }
    }
}
