package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type ReturnContacao struct {
	Bid string `json:"bid"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
		//panic(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var data ReturnContacao

	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err)
	}
	var file *os.File
	defer file.Close()
	_, err = os.Stat("cotacao.txt")
	if os.IsNotExist(err) {
		file, err = os.Create("cotacao.txt")
		if err != nil {
			panic(err)
		}

	} else {
		file, err = os.OpenFile("cotacao.txt", os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
	}

	currentTime := time.Now()

	currentTimeString := currentTime.Format("2006-01-02 15:04:05")
	_, err = file.Write([]byte("Dolar: " + data.Bid + " - " + currentTimeString + "\n"))
	if err != nil {
		panic(err)
	}
}
