CREATE TABLE users(
    id SERIAL primary key,
    name VARCHAR(50),
    age NUMERIC,
    email VARCHAR(100) UNIQUE,
    password VARCHAR(100)
);

CREATE TABLE movies(
    id SERIAL primary key,
    name VARCHAR(50),
    description VARCHAR(200),
    cover VARCHAR(200),
    date DATE,
    rate NUMERIC,
    user_created_id INTEGER
);

CREATE TABLE watchlist(
    id SERIAL primary key,
    user_id int,
    movie_id int,
    date DATE
);

CREATE TABLE reviews(
    id SERIAL primary key,
    user_id int,
    movie_id int,
    date DATE,
    review VARCHAR(300),
    rate NUMERIC
);
