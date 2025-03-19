CREATE TABLE experiences (
                             id VARCHAR(26) PRIMARY KEY,
                             user_id VARCHAR(26) NOT NULL,
                             image_url VARCHAR(255),
                             job_title VARCHAR(255) NOT NULL,
                            job_location VARCHAR(255) NOT NULL,
                             skill_used VARCHAR(255),
                             start_date VARCHAR(50) NOT NULL,
                             end_date VARCHAR(50),
                             description TEXT,
                             created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                             updated_at TIMESTAMP,
                             FOREIGN KEY (user_id) REFERENCES users(id)
);