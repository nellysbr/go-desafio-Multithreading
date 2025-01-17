package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Address struct {
	CEP         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	UF          string `json:"uf"`
	Source      string
}

func main() {
	cep := "01153000"

	resultChan := make(chan Address)

	go fetchBrasilAPI(cep, resultChan)
	go fetchViaCEP(cep, resultChan)

	timeout := time.After(1 * time.Second)

	select {
	case result := <-resultChan:
		printAddress(result)
	case <-timeout:
		fmt.Println("Erro: Timeout - Nenhuma API respondeu em tempo hábil.")
	}
}

func fetchBrasilAPI(cep string, resultChan chan<- Address) {
	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)
	fetchAPI(url, "BrasilAPI", resultChan)
}

func fetchViaCEP(cep string, resultChan chan<- Address) {
	url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)
	fetchAPI(url, "ViaCEP", resultChan)
}

func fetchAPI(url, source string, resultChan chan<- Address) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var address Address
	err = json.Unmarshal(body, &address)
	if err != nil {
		return
	}

	address.Source = source
	resultChan <- address
}

func printAddress(address Address) {
	fmt.Printf("Resultado da API mais rápida (%s):\n", address.Source)
	fmt.Printf("CEP: %s\n", address.CEP)
	fmt.Printf("Logradouro: %s\n", address.Logradouro)
	fmt.Printf("Complemento: %s\n", address.Complemento)
	fmt.Printf("Bairro: %s\n", address.Bairro)
	fmt.Printf("Localidade: %s\n", address.Localidade)
	fmt.Printf("UF: %s\n", address.UF)
}
