CREATE TABLE users (
   id BIGINT AUTO_INCREMENT PRIMARY KEY,
   uuid CHAR(36) NOT NULL UNIQUE,
   email VARCHAR(255) NOT NULL UNIQUE,
   password TEXT NOT NULL,
   created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
   updated_at DATETIME NULL DEFAULT NULL,
   deleted_at DATETIME NULL DEFAULT NULL
);