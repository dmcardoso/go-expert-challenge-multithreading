package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type CEP struct {
	Cep      string `json:"cep"`
	State    string `json:"state"`
	City     string `json:"city"`
	District string `json:"district"`
	Address  string `json:"address"`
	Origin   string `json:"origin"`
}

type ViaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	UF          string `json:"uf"`
	IBGE        string `json:"ibge"`
	GIA         string `json:"gia"`
	DDD         string `json:"ddd"`
	SIAFI       string `json:"siafi"`
}

type ApiCEP struct {
	Code       string `json:"code"`
	State      string `json:"state"`
	City       string `json:"city"`
	District   string `json:"district"`
	Address    string `json:"address"`
	Status     int    `json:"status"`
	Ok         bool   `json:"ok"`
	StatusText string `json:"statusText"`
}

func NewCEP(cep, state, city, district, address, origin string) *CEP {
	return &CEP{Cep: cep, State: state, City: city, District: district, Address: address, Origin: origin}
}

func main() {
	cep := os.Args[1:][0]
	channel := make(chan *CEP)

	go SearchCEPApiCEP(cep, channel)
	go SearchCEPViaCEP(cep, channel)

	select {
	case cepResponse := <-channel:
		fmt.Printf("CEP response -> CEP: %v, State: %v, City: %v, District: %v, Address: %v, Origin: %v\n", cepResponse.Cep, cepResponse.State, cepResponse.City, cepResponse.District, cepResponse.Address, cepResponse.Origin)
	case <-time.After(time.Second * 1):
		println("timeout")
	}
}

func SearchCEPViaCEP(cep string, channel chan<- *CEP) {
	origin := "http://viacep.com.br/ws/" + cep + "/json/"
	req, err := http.NewRequest(http.MethodGet, origin, nil)

	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		panic(err)
	}

	var viaCEPResponse ViaCEP
	err = json.Unmarshal(body, &viaCEPResponse)

	if err != nil {
		panic(err)
	}

	cepReponse := NewCEP(viaCEPResponse.Cep, viaCEPResponse.UF, viaCEPResponse.Localidade, viaCEPResponse.Bairro, viaCEPResponse.Logradouro+", "+viaCEPResponse.Complemento, origin)

	channel <- cepReponse
}

func SearchCEPApiCEP(cep string, channel chan<- *CEP) {
	origin := "https://cdn.apicep.com/file/apicep/" + cep + ".json"
	req, err := http.NewRequest(http.MethodGet, origin, nil)

	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		panic(err)
	}

	var apiCEPResponse ApiCEP
	err = json.Unmarshal(body, &apiCEPResponse)

	if err != nil {
		panic(err)
	}

	cepReponse := NewCEP(apiCEPResponse.Code, apiCEPResponse.State, apiCEPResponse.City, apiCEPResponse.District, apiCEPResponse.Address, origin)

	channel <- cepReponse
}
