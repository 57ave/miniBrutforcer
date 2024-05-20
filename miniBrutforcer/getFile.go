package main

import (
    "os"
    "fmt"
    "log"
	"strings"
)

func GetFileContent(filePath string) []string{
    body, err := os.ReadFile(filePath)
    if err != nil {
        log.Fatalf("unable to read file: %v", err)
    }
	dataTab := strings.Split(string(body), ";\n")
    fmt.Println(dataTab[0])
	return dataTab
}
