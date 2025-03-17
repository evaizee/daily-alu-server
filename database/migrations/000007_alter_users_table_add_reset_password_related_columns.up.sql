ALTER TABLE users
ADD reset_password_token VARCHAR(255),
ADD reset_password_requested_at TIMESTAMP;