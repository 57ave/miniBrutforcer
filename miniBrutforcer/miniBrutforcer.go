
package main

import (
  "fmt"
  "io/ioutil"
  "log"
  "net/http"
  "time"
  "strings"
)

func queryExecute(path string, key string, value string, method string) {

  client := &http.Client{}
  req, err := http.NewRequest(method, "http://localhost:3333" + path, nil)
  if err != nil {
    log.Fatal(err)
  }

  q := req.URL.Query()
  q.Add(key, value)

  req.URL.RawQuery = q.Encode()

  resp, err := client.Do(req)
  if err != nil {
    fmt.Println(err)
  }

  defer resp.Body.Close()

  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    log.Fatal(err)
  }

  fmt.Println(string(body))
}

func executeQueryFromFile(filePathQuery []string, filePathRoot chan string) {
  for eachtask := range filePathRoot {
    for i := 0 ; i < len(filePathQuery) ; i++ {
      splitData := strings.Split(filePathQuery[i], "=")
      queryExecute(eachtask, splitData[0], splitData[1], "POST")
    }
  }
}

func main() {
  filePathTab := "./wordList/query1.csv"
  channel := make(chan string)
  data_tab := GetFileContent("./wordList/rootList")

  for i := 0 ;i < len(data_tab); i++ {
		go executeQueryFromFile(GetFileContent(filePathTab), channel)
	}

  for i := 0; i < len(data_tab); i++ {
    channel <- data_tab[i]
  }
  time.Sleep(1 * time.Second)
  close(channel)
}
