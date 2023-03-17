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

type Cotacao struct {
	Valor string `json:"cotacao"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		fmt.Printf("Erro ao fazer a requisição: %v\n", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Erro ao processar a requisição: %v\n", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Erro ao fazer o parse da requisição: %v\n", err)
	}

	var c Cotacao
	err = json.Unmarshal(body, &c)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao fazer o parse dos dados: %v\n", err)
	}

	file, err := os.Create("cotacao.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao criar o arquivo: %v", err)
	}
	_, err = file.WriteString(fmt.Sprintf(`Dólar: {%s}`, c.Valor))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao escrever no arquivo: %v", err)
	}
	defer file.Close()

	fmt.Printf("Arquivo criado com sucesso!!\n")
}
