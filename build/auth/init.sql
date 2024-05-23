CREATE TABLE IF NOT EXISTS Users(
    id SERIAL PRIMARY KEY,
    username text NOT NULL,
    password text NOT NULL
)