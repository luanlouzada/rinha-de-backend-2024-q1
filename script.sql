CREATE TABLE clientes (
    id SERIAL PRIMARY KEY,
    limite INT NOT NULL,
    saldo INT NOT NULL DEFAULT 0
);

CREATE TABLE transacoes (
    id SERIAL PRIMARY KEY,
    cliente_id INT REFERENCES clientes(id) ON DELETE CASCADE,
    valor INT NOT NULL,
    tipo CHAR(1) NOT NULL CHECK (tipo IN ('c', 'd')),
    descricao VARCHAR(10) NOT NULL,
    realizada_em TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO clientes (id, limite, saldo) VALUES
(1, 100000, 0),
(2, 80000, 0),
(3, 1000000, 0),
(4, 10000000, 0),
(5, 500000, 0);


