CREATE TABLE auth_tokens (
     id BIGINT PRIMARY KEY AUTO_INCREMENT,
     user_id INT NOT NULL,
     token VARCHAR(512) NOT NULL,
     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
     expired_at TIMESTAMP NOT NULL,
     status TINYINT DEFAULT 1,
     INDEX idx_user_id (user_id),
     INDEX idx_token (token)
);

