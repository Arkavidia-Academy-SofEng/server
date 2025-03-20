CREATE TABLE portfolios (
                           id VARCHAR(26) PRIMARY KEY,
                           user_id VARCHAR(26) NOT NULL,
                           image VARCHAR(255),
                           project_name VARCHAR(255) NOT NULL,
                           project_location VARCHAR(255) NOT NULL,
                           description_image VARCHAR(255),
                           project_link VARCHAR(255),
                           start_date VARCHAR(50) NOT NULL,
                           end_date VARCHAR(50),
                           description TEXT,
                           created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                           updated_at TIMESTAMP,
                           FOREIGN KEY (user_id) REFERENCES users(id)
);