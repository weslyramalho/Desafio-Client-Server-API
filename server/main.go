package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

	cot, error := BuscaCotacao(ctx)
	if error != nil {
		panic(fmt.Sprintf("Falha ao tentar pegar cotação: %v", error))
	}

	ctx = nil
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Nanosecond)

	lastID, err := salvarCotacao(ctx, db, cot)

	if err != nil {
		panic(fmt.Sprintf("Erro do buscaCotacao %v", err))
	}
	fmt.Printf("main.CotacaoHandler - Último ID inserido: %d", lastID)

	json.NewEncoder(w).Encode(cot)

}
func BuscaCotacao(c context.Context) (map[string]Cotacaodb, error) {
	f, err := os.Create("arquivos.txt")
	if err != nil {
		panic(err)
	}
	req, err := http.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error ao fazer requisição: %v\n", err)
	}
	defer req.Body.Close()
	res, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao ler a resposta: %v\n", err)
	}
	f.Write(res)
	var data map[string]Cotacaodb
	err = json.Unmarshal(res, &data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao fazer parse da resposta: %v\n", err)
	}
	return data, nil
}

func salvarCotacao(c context.Context, db *sql.DB, cota map[string]Cotacaodb) (int64, error) {
	obj := cota["USDBRL"]

	ct := `
		INSERT INTO
			cotacoes(
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

	select {
	case <-c.Done():
		return int64(0), errors.New("error tempo exedido")
	}

	re, err := db.Exec(ct,
		obj.Code,
		obj.Codein,
		obj.Name,
		obj.High,
		obj.Low,
		obj.VarBid,
		obj.PctChange,
		obj.Bid,
		obj.Ask,
		obj.Timestamp,
		obj.CreateDate,
	)

	if err != nil {
		return int64(0), err
	}

	// Inserir o registro na tabela
	lastID, err := re.LastInsertId()

	if err != nil {
		return int64(0), err
	}

	return lastID, nil
	/*

		stmt, err := db.Prepare("INSERT INTO cotacoes(code,codein,name,high,low,varbid,pctchange,bid,ask,timestamp,create_date) VALUES (?,?,?,?,?,?,?,?,?,?,?)")
		if err != nil {
			panic(err)
		}

		obj := cota["USDBRL"]

		defer stmt.Close()
		obj, err := stmt.Exec(cota.Code, cota.Codein, cota.Ask, cota.Bid, cota.CreateDate, cota.High, cota.Low, cota.Name, cota.PctChange, cota.Timestamp, cota.VarBid)
		if err != nil {
			return int64(0), err
		}
		lastID, err := rt.LastInsertId()

		if err != nil {
			return int64(0), err
		}
		return lastID, nil
	*/
}
