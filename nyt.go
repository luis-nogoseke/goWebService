package main
import (
  "encoding/json"
  "fmt"
  "log"
  "net/http"
  "net/url"
  "os"
  "strings"
  "time"
)

const IssuesURL = "https://api.nytimes.com/svc/books/v3/lists.json"


type IssuesSearchResult struct {
  Status string `json:"status"`
  Copyrigth string `json:"copyright"`
  N_results int `json:"num_results"`
  Last string `json:"last_modified"`
  Results []*Issue
}

type Issue struct {
  List string `json:"list_name"`
  Name string `json:"display_name"`
  Date string `json:"bestsellers_date"`
  Pdate string `json:"published_date"`
  Rank int
  Lrank int `json:"rank_last_week"`
  Weeks int `json:"weeks_on_list"`
  Asterisk int
  Dagger int
  Amazon_url string `json:"amazon_product_url"`
  Isbns []*Isbn
  Details []*Detail `json:"book_details"`
  Reviews []*Review
}

type Isbn struct {
  Isbn10 string
  Isbn13 string
}

type Detail struct {
  Title string
  Description string
  Contributor string
  Author string
  Note string `json:"contributor_note"`
  Price int
  Age string `json:"age_group"`
  Publisher string
  Isbn13 string `json:"primary_isbn13"`
  Isbn10 string `json:"primary_isbn10"`
}


type Review struct {
  Link string `json:"book_review_link"`
  First_ch string `json:"first_chapter_link"`
  Sunday string `json:"sunday_review_link"`
  Article string `json:"article_chapter_link"`
}

type User struct {
  Login string
  HTMLURL string `json:"html_url"`
}

func SearchIssues (terms []string) (*IssuesSearchResult, error) {
  q := url.QueryEscape(strings.Join(terms, " "))
  req, err := http.NewRequest("GET", IssuesURL+"?api-key=054670dd36264476b6456e1b8e24c7d1&list="+q, nil)

  if err != nil {
    return nil, err
  }

  resp, err := http.DefaultClient.Do(req)
  if err != nil {
    return nil, err
  }

  defer resp.Body.Close()

  if resp.StatusCode != http.StatusOK {
    return nil, fmt.Errorf("search query failed: %s", resp.Status)
  }

  var result IssuesSearchResult
  if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
    return nil, err
  }

  return &result, nil
}

func daysAgo (t time.Time) int {
  return int(time.Since(t).Hours() / 24)
}

func main () {
  fmt.Println("Rodando")
  result, err := SearchIssues(os.Args[1:])
  if err != nil {
    log.Fatal(err)
  }
  fmt.Printf("%d issues:\n", result.N_results)
  fmt.Println(result.Status)
  fmt.Println(result.Copyrigth)
  fmt.Println(result.Last)
  for _, item := range result.Results {
    fmt.Println(item.Name)
    fmt.Println(item.Details[0].Title)
  }
}
