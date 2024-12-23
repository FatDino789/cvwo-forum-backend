-- In sql/create_tables.sql
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert mock data
INSERT INTO users (username, email, password_hash) VALUES
('john_doe', 'john@example.com', 'hash123'),
('jane_smith', 'jane@example.com', 'hash456');

INSERT INTO posts (user_id, title, content) VALUES
(1, 'First Post', 'This is my first post content'),
(1, 'Second Post', 'Another post by John'),
(2, 'Jane''s Post', 'Hello from Jane!');