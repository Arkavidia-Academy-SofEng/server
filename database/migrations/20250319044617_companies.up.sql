CREATE TABLE companies (
                           id VARCHAR(26) PRIMARY KEY,
                           profile_picture VARCHAR(255),
                           banner_picture VARCHAR(255),
                           email VARCHAR(255) UNIQUE NOT NULL,
                           password VARCHAR(255),
                           name VARCHAR(255) NOT NULL,
                           about_us TEXT,
                           industry_types VARCHAR(255),
                           number_employees INTEGER,
                           established_date DATE,
                           company_url VARCHAR(255),
                           required_skill TEXT,
                           location TEXT,
                           phone_number VARCHAR(20),
                           created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                           updated_at TIMESTAMP,
                           deleted_at TIMESTAMP
);