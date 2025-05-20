-- db/migrations/<timestamp>_create_purchases_table.up.sql
CREATE TABLE purchases (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    lottery_id INT NOT NULL,
    numbers INT[] NOT NULL
);
