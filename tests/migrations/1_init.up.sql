CREATE TABLE "user" (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL
);

CREATE TABLE wallet (
    id SERIAL PRIMARY KEY,
    user_id INT,
    coins INT CHECK (coins >= 0),
    FOREIGN KEY (user_id) REFERENCES "user" (id) ON DELETE CASCADE
);

CREATE TABLE item (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    price INT
);

CREATE TABLE purchase (
    id SERIAL PRIMARY KEY,
    user_id INT,
    item_id INT,
    FOREIGN KEY (user_id) REFERENCES "user" (id) ON DELETE CASCADE,
    FOREIGN KEY (item_id) REFERENCES item (id) ON DELETE CASCADE
);

CREATE TABLE transaction (
    id SERIAL PRIMARY KEY,
    from_user INT,
    to_user INT,
    amount INT CHECK (amount >= 0),
    created_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (from_user) REFERENCES "user" (id) ON DELETE CASCADE,
    FOREIGN KEY (to_user) REFERENCES "user" (id) ON DELETE CASCADE
);

INSERT INTO item (name, price) VALUES
('t-shirt', 80),
('cup', 20),
('book', 50),
('pen', 10),
('powerbank', 200),
('hoody', 300),
('umbrella', 200),
('socks', 10),
('wallet', 50),
('pink-hoody', 500);