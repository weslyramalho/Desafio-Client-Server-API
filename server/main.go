package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Cotacao struct {
	Usdbrl struct {
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
	} `json:"USDBRL"`
}

type Cotacaodb struct {
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

func newCota(code string, codein string, name string, high string, low string, varbid string, pctchange string, bid string, ask string, timestamp string, create_date string) *Cotacaodb {
	return &Cotacaodb{
		Code:       code,
		Codein:     code,
		Name:       name,
		High:       high,
		Low:        low,
		VarBid:     varbid,
		PctChange:  pctchange,
		Bid:        bid,
		Ask:        ask,
		Timestamp:  timestamp,
		CreateDate: create_date,
	}
}

func main() {
	http.HandleFunc("/cotacao", BuscaCotacaoHandler)
	http.ListenAndServe(":8080", nil)

}
func BuscaCotacaoHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "cotacao.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	w.Header().Set("Content-Type", "application/json")
	ctx, cancel := context.WithTimeout(r.Context(), 2000*time.Millisecond)
	defer cancel()

	cota, error := BuscaCotacao(ctx)
	if error != nil {
		panic(fmt.Sprintf("Falha ao tentar pegar cotação: %v", error))
	}

	cot, _ := ioutil.ReadAll(cota)
	ctx = nil
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Nanosecond)
	cota, err := newCota(cot.Usdbrl.Code, cot.Usdbrl.Codein, cot.Usdbrl.Name, cot.Usdbrl.High, cot.Usdbrl.Low, cot.Usdbrl.VarBid, cot.Usdbrl.PctChange, cot.Usdbrl.Bid, cot.Usdbrl.Ask, cot.Usdbrl.Timestamp, cot.Usdbrl.CreateDate)
	err := salvarCotacao(ctx, db, cota)

	if err != nil {
		panic(fmt.Sprintf("Erro do buscaCotacao %v", err))
	}

	json.NewEncoder(w).Encode(cot)

}
func BuscaCotacao(c context.Context) (*Cotacao, error) {
	req, err := http.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error ao fazer requisição: %v\n", err)
	}
	defer req.Body.Close()
	res, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao ler a resposta: %v\n", err)
	}

	var data Cotacao
	err = json.Unmarshal(res, &data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao fazer parse da resposta: %v\n", err)
	}
	return &data, nil
}

func salvarCotacao(c context.Context, db *sql.DB, cota *Cotacaodb) (int64, error) {
	stmt, err := db.Prepare("INSERT INTO cotacoes(code,codein,name,high,low,varbid,pctchange,bid,ask,timestamp,create_date) VALUES (?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		panic(err)
	}

	defer stmt.Close()
	rt, err := stmt.Exec(cota.Code, cota.Codein, cota.Ask, cota.Bid, cota.CreateDate, cota.High, cota.Low, cota.Name, cota.PctChange, cota.Timestamp, cota.VarBid)
	if err != nil {
		return int64(0), err
	}
	lastID, err := rt.LastInsertId()

	if err != nil {
		return int64(0), err
	}
	return lastID, nil
}
