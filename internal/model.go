package internal

import "time"

type Account struct {
	Transactions []Transaction
	Balanco      int `json:"balanco"`
	Limite       int `json:"limite"`
}

const (
	Debito  = "d"
	Credito = "c"
)

type Transaction struct {
	RealizadaEm time.Time `json:"realizada_em"`
	Tipo        string    `json:"tipo"`
	Descricao   string    `json:"descricao"`
	Valor       int       `json:"valor"`
}

type Extract struct {
	DataExtrato       string `json:"data_extrato"`
	UltimasTransacoes []Transaction
	Saldo             int `json:"saldo"`
	Limite            int `json:"limite"`
}
