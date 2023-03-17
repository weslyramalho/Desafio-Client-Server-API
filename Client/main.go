package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type Dolar struct {
	Bid string `json:"bid"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3000*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(fmt.Sprintf("Erro ao fazer requisição: %v", err))
	}

	f, err := os.Create("cotacao.txt")
	if err != nil {
		panic(fmt.Sprintf("Erro ao salvar cotaçao: %v", err))
	}

	var dolar Dolar
	defer f.Close()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(fmt.Sprintf("Erro ao fazer requisição %v", err))
	}

	defer res.Body.Close()

	cota, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(fmt.Sprintf("Erro ao pegar dados de requisição: %v", err))
	}

	err = json.Unmarshal(cota, &dolar)
	if err != nil {
		panic(fmt.Sprintf("Erro ao tranformar json: %v", err))
	}

	var d string = "Dólar: "
	_, err = f.Write([]byte(d))
	if err != nil {
		panic(fmt.Sprintf("Erro ao salvar arquivo txt: %v", err))
	}
	_, err = f.Write([]byte(dolar.Bid))
	if err != nil {
		panic(fmt.Sprintf("Erro ao salvar arquivo txt: %v", err))
	}
}
