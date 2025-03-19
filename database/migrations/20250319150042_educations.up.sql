CREATE TABLE educations (
                           id VARCHAR(26) PRIMARY KEY,
    user_id VARCHAR(26) NOT NULL,
                           image VARCHAR(255),
                           title_degree VARCHAR(255) NOT NULL,
                           institutional_name VARCHAR(255) NOT NULL,
                           start_date VARCHAR(50) NOT NULL,
                           end_date VARCHAR(50),
                           description TEXT,
                           created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                           updated_at TIMESTAMP
);