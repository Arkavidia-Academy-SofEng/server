CREATE TABLE users (
                       id VARCHAR(26) PRIMARY KEY,
                       email VARCHAR(255) UNIQUE NOT NULL,
                       password VARCHAR(255),
                       name VARCHAR(255) NOT NULL,
                       role VARCHAR(50) NOT NULL,
                        phone_number VARCHAR(20),
                       profile_picture VARCHAR(255),
                        banner_picture VARCHAR(255),
                        about_me TEXT,
                       is_premium BOOLEAN DEFAULT FALSE,
                       premium_until TIMESTAMP,
                       headline VARCHAR(255),
                       location TEXT,
                       created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP,
                       deleted_at TIMESTAMP
);