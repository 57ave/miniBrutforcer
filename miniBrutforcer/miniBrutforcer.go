
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

func executeQueryFromFile(filePathPromise chan string) {
  for filePath := range filePathPromise {
    dataTab := GetFileContent(filePath)
    for i := 0 ; i < len(dataTab) ; i++ {
      splitData := strings.Split(dataTab[i], "=")
      queryExecute("/admin", splitData[0], splitData[1], "POST")
    }
  }
}

func main() {
  filePathTab := []string{
    "./wordList/query1.csv",
    "./wordList/query2.csv",
    "./wordList/query3.csv",
    "./wordList/query4.csv",
  }
  channel := make(chan string)

  for i := 0; i < 2; i++ {
		go executeQueryFromFile(channel)
	}

  for i := 0; i < len(filePathTab); i++ {
    channel <- filePathTab[i]
  }
  time.Sleep(1 * time.Second)
  close(channel)
}
