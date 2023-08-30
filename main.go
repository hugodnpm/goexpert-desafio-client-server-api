package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Dolar struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

type ReturnContacao struct {
	Bid string `json:"bid"`
}

func main() {
	http.HandleFunc("/cotacao", handler)
	http.ListenAndServe(":8080", nil)

}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/cotacao" {
		w.WriteHeader(http.StatusNotFound)
	}

	dolar, error := QueryExchange()
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dolar)

}
func QueryExchange() (*ReturnContacao, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var data map[string]Dolar
	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		return nil, err
	}
	usdbrl := data["USDBRL"]

	cotacao := ReturnContacao{
		Bid: usdbrl.Bid,
	}
	err = conexaoDB(&usdbrl)
	if err != nil {
		return nil, err
	}
	return &cotacao, nil
}
func conexaoDB(usdbrl *Dolar) error {
	db, err := sql.Open("sqlite3", "cotacao.db")
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS cotacao (id INTEGER PRIMARY KEY, code TEXT,codein TEXT, name TEXT, high TEXT, low TEXT, varbid TEXT, pctchange TEXT, bid TEXT, ask TEXT, timestamp TEXT, createdate TEXT)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	err = insertCotacao(db, usdbrl)
	if err != nil {
		return err
	}
	return nil
}
func insertCotacao(db *sql.DB, dolar *Dolar) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	stmt, err := db.Prepare("insert into cotacao(code, codein,name,high,low,varBid,pctChange,bid,ask,timestamp,createDate) values(?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, dolar.Code, dolar.Codein, dolar.Name, dolar.High, dolar.Low, dolar.VarBid, dolar.PctChange, dolar.Bid, dolar.Ask, dolar.Timestamp, dolar.CreateDate)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
