-- Coloque scripts iniciais aqui
CREATE TABLE clientes (
  id SERIAL NOT NULL PRIMARY KEY,
  nome VARCHAR(255) NOT NULL,
  limite INTEGER NOT NULL,
  saldo_inicial INTEGER NOT NULL DEFAULT 0
);


CREATE TABLE transacoes (
  id SERIAL NOT NULL PRIMARY KEY,
  amount INTEGER NOT NULL,
  type VARCHAR(1) NOT NULL,
  customer_id INTEGER NOT NULL,
  FOREIGN KEY (customer_id) REFERENCES clientes(id)
);


DO $$
BEGIN
  INSERT INTO clientes (nome, limite)
  VALUES
    ('o barato sai caro', 1000 * 100),
    ('zan corp ltda', 800 * 100),
    ('les cruders', 10000 * 100),
    ('padaria joia de cocaia', 100000 * 100),
    ('kid mais', 5000 * 100);
END; $$
