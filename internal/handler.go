package internal

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

type Handler struct {
	DB *sql.DB
}

func (h *Handler) HandleTransaction(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var transacao Transaction

	if err := json.NewDecoder(r.Body).Decode(&transacao); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Iniciar transação
	tx, err := h.DB.Begin()
	if err != nil {
		http.Error(w, "Erro ao iniciar transação no banco de dados", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Atualizar saldo e inserir transação
	var saldoAtual, limite int
	err = tx.QueryRow("SELECT saldo, limite FROM clientes WHERE id = $1 FOR UPDATE", id).Scan(&saldoAtual, &limite)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Cliente não encontrado", http.StatusNotFound)
		} else {
			http.Error(w, "Erro ao buscar cliente", http.StatusInternalServerError)
		}
		return
	}

	if transacao.Tipo == "d" && saldoAtual-transacao.Valor < -limite {
		http.Error(w, "Saldo insuficiente", http.StatusUnprocessableEntity)
		return
	}

	novoSaldo := saldoAtual
	if transacao.Tipo == "c" {
		novoSaldo += transacao.Valor
	} else {
		novoSaldo -= transacao.Valor
	}

	_, err = tx.Exec("UPDATE clientes SET saldo = $1 WHERE id = $2", novoSaldo, id)
	if err != nil {
		http.Error(w, "Erro ao atualizar saldo do cliente", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec("INSERT INTO transacoes (cliente_id, valor, tipo, descricao) VALUES ($1, $2, $3, $4)", id, transacao.Valor, transacao.Tipo, transacao.Descricao)
	if err != nil {
		http.Error(w, "Erro ao registrar transação", http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, "Erro ao finalizar transação no banco de dados", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]int{"limite": limite, "saldo": novoSaldo})
}

func (h *Handler) HandleExtract(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	extrato := Extract{}

	rows, err := h.DB.Query(`
        SELECT valor, tipo, descricao, realizada_em
        FROM transacoes
        WHERE cliente_id = $1
        ORDER BY realizada_em DESC
        LIMIT 10`, id)
	if err != nil {
		http.Error(w, "Erro ao buscar transações", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var t Transaction
		if err := rows.Scan(&t.Valor, &t.Tipo, &t.Descricao, &t.RealizadaEm); err != nil {
			http.Error(w, "Erro ao ler transação", http.StatusInternalServerError)
			return
		}
		extrato.UltimasTransacoes = append(extrato.UltimasTransacoes, t)
	}

	err = h.DB.QueryRow("SELECT saldo, limite FROM clientes WHERE id = $1", id).Scan(&extrato.Saldo, &extrato.Limite)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Cliente não encontrado", http.StatusNotFound)
		} else {
			http.Error(w, "Erro ao buscar saldo do cliente", http.StatusInternalServerError)
		}
		return
	}

	extrato.DataExtrato = time.Now().Format(time.RFC3339)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(extrato)
}
