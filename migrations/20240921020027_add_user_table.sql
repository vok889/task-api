-- +goose Up
CREATE TABLE users (
    id          bigserial PRIMARY KEY,
    username    VARCHAR(255) UNIQUE NOT NULL,
    password    VARCHAR(255) NOT NULL,
    email    	VARCHAR(100),
    first_name  VARCHAR(100),
    last_name   VARCHAR(100)
);

-- +goose Down
DROP TABLE users;