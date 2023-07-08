-- +migrate Up
CREATE TABLE users(
                        id INT PRIMARY KEY AUTO_INCREMENT,
                        name VARCHAR(255) NOT NULL,
                        phone_number VARCHAR(255) NOT NULL UNIQUE,
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        password VARCHAR(255) NOT NULL
);

-- +migrate Down
TABLE users;