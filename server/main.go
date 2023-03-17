package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

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

type Dolar struct {
	Bid string `json:"bid"`
}

func main() {
	http.HandleFunc("/cotacao", cotaHandler)
	http.ListenAndServe(":8080", nil)

}
func cotaHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, cancel := context.WithTimeout(r.Context(), 2000*time.Millisecond)
	defer cancel()

	var cot map[string]Cotacaodb
	cot, err := BuscaCotacao(ctx)
	if err != nil {
		panic(fmt.Sprintf("Erro ao buscar cotação %v", err))
	}

	ctx = nil
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Nanosecond)

	err = salvCota(ctx, cot)
	if err != nil {
		panic(fmt.Sprintf("Erro ao salvar cotação %v", err))
	}
	dol := Dolar{Bid: cot["USDBRL"].Bid}

	json.NewEncoder(w).Encode(dol)

}

func BuscaCotacao(c context.Context) (map[string]Cotacaodb, error) {
	req, err := http.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error ao fazer requisição: %v\n", err)
	}
	defer req.Body.Close()

	res, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao ler a resposta: %v\n", err)
	}

	var data map[string]Cotacaodb
	err = json.Unmarshal(res, &data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao fazer parse da resposta: %v\n", err)
	}
	return data, nil
}
func salvCota(c context.Context, cota map[string]Cotacaodb) error {
	db, err := sql.Open("sqlite3", "cotacao.db")
	if err != nil {
		return err
	}

	defer db.Close()
	cot := cota["USDBRL"]
	sts := `
	CREATE TABLE cotacao(id INTEGER PRIMARY KEY,code TEXT, codein TEXT, name TEXT, high TEXT, low TEXT, varbid TEXT, pctchange TEXT, bid TEXT, ask TEXT, timestamp TEXT, create_date TEXT);
			INSERT INTO
				cotacao(
					code,
					codein,
					name,
					high,
					low,
					varbid,
					pctchange,
					bid,
					ask,
					timestamp,
					create_date
				)
				VALUES (
					?,
					?,
					?,
					?,
					?,
					?,
					?,
					?,
					?,
					?,
					?
				);
	`

	result, err := db.Exec(sts, cot.Code, cot.Codein, cot.Name, cot.High, cot.Low, cot.VarBid, cot.PctChange, cot.Bid, cot.Ask, cot.Timestamp, cot.CreateDate)

	if err != nil {
		log.Fatal(err)
	}
	lastID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	fmt.Println(lastID)
	fmt.Println("table data created")
	return nil
}
