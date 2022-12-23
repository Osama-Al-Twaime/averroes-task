CREATE TABLE users(
    name VARCHAR(50),
    age NUMERIC,
    email VARCHAR(100) UNIQUE,
    password VARCHAR(100)
);