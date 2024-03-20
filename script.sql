-- Coloque scripts iniciais aqui
CREATE TABLE customers (
  id SERIAL NOT NULL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  "limit" INTEGER NOT NULL,
  balance INTEGER NOT NULL DEFAULT 0
);


CREATE TABLE transactions (
  id SERIAL NOT NULL PRIMARY KEY,
  amount INTEGER NOT NULL,
  "type" VARCHAR(1) NOT NULL,
  customer_id INTEGER NOT NULL,
  FOREIGN KEY (customer_id) REFERENCES customers(id)
);


DO $$
BEGIN
  INSERT INTO customers (name, "limit")
  VALUES
    ('o barato sai caro', 1000 * 100),
    ('zan corp ltda', 800 * 100),
    ('les cruders', 10000 * 100),
    ('padaria joia de cocaia', 100000 * 100),
    ('kid mais', 5000 * 100);
END; $$
