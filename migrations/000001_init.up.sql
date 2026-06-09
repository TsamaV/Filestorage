CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(200) UNIQUE NOT NULL,
    password_hash VARCHAR(200) NOT NULL
);

CREATE TABLE files (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    url VARCHAR(1000) NOT NULL,
    file_name VARCHAR(200) NOT NULL,
    size INT NOT NULL,
    created_at TIMESTAMP NOT NULL
);