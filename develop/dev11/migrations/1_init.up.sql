CREATE TABLE IF NOT EXISTS events
(
    id serial PRIMARY KEY,
    user_id        integer NOT NULL,
    date date NOT NULL,
    event     VARCHAR(500) NOT NULL
);