package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type CambioAPI struct {
	USDBRL struct {
		ID         string
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

func main() {
	db, err := sql.Open("sqlite3", "./cambio.db")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS cambio (id VARCHAR(255), cotacao DECIMAL(10,2))")
	if err != nil {
		panic(err)
	}
	db.Close()

	http.HandleFunc("/cotacao", ConsultaCambioHanlder)
	http.ListenAndServe(":8080", nil)
}

func ConsultaCambioHanlder(w http.ResponseWriter, r *http.Request) {
	cambio, err := consultaCambio()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = insertCambio(cambio)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(cambio.USDBRL.Bid)
	w.Write([]byte(`{ "cotacao": "` + cambio.USDBRL.Bid + `"}`))
}

func consultaCambio() (*CambioAPI, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)

	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var c CambioAPI
	err = json.Unmarshal(body, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func insertCambio(cambio *CambioAPI) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	db, err := sql.Open("sqlite3", "./cambio.db")
	if err != nil {
		panic(err)
	}

	stmt, err := db.Prepare("insert into cambio (id, cotacao) values(?,?)")
	if err != nil {
		fmt.Printf("%v", err)
		return err
	}
	cambio.USDBRL.ID = uuid.New().String()
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, cambio.USDBRL.ID, cambio.USDBRL.Bid)
	if err != nil {
		return err
	}

	return nil
}
