-- +goose Up
CREATE TABLE 
    items (
        id bigserial PRIMARY KEY,
        title VARCHAR(255) NOT NULL,
        amount INT NOT NULL,
        quantity INT NOT NULL,
        status VARCHAR(20) NOT NULL,
        owner_id INT NOT NULL
        );

-- insert seed data
INSERT INTO items (title, amount, quantity, status, owner_id) VALUES ('iphone16', 40000, 10, 'REJECTED', 1);
INSERT INTO items (title, amount, quantity, status, owner_id) VALUES ('airpod', 2000, 20, 'APPROVED', 1);
INSERT INTO items (title, amount, quantity, status, owner_id) VALUES ('ipad', 30000, 30, 'PENDING', 1);

-- +goose Down
DROP TABLE items;