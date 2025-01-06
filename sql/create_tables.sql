-- Drop existing tables
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS users;

-- Recreate users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create new posts table with comments as JSONB
CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    picture_url VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    likes_count INTEGER DEFAULT 0,
    views_count INTEGER DEFAULT 0,
    discussion_thread TEXT,
    comments JSONB DEFAULT '[]',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert users
INSERT INTO users (email, password_hash) VALUES
('testing@gmail.com', '$2a$10$mN6CaIxk7mU0QM3B2Q490euGHJS5Dx0AOTjG7v82f9dQL/Gm.gCEC'),
('john@example.com', '$2a$10$mN6CaIxk7mU0QM3B2Q490euGHJS5Dx0AOTjG7v82f9dQL/Gm.gCEC'),
('jane@example.com', '$2a$10$mN6CaIxk7mU0QM3B2Q490euGHJS5Dx0AOTjG7v82f9dQL/Gm.gCEC');

-- Insert posts with comments as JSONB
INSERT INTO posts (user_id, title, content, picture_url, created_at, likes_count, views_count, discussion_thread, comments) VALUES 
(1, 'Exchange Experience in Tokyo', 
   'Sharing my amazing semester abroad experience at Waseda University!', 
   'https://example.com/tokyo-campus.jpg', 
   '2024-01-15 10:00:00', 
   45, 
   230, 
   'Looking for advice on accommodation and cultural adjustments',
   '[
     {
       "id": 1,
       "user_id": 2,
       "content": "Great post! How did you handle the language barrier?",
       "created_at": "2024-01-15 10:30:00"
     },
     {
       "id": 2,
       "user_id": 3,
       "content": "The campus looks amazing! Did you stay in university accommodation?",
       "created_at": "2024-01-15 11:00:00"
     }
   ]'::jsonb
),
(2, 'NUS Exchange Application Guide', 
   'Step-by-step guide on applying for exchange at NUS', 
   'https://example.com/nus-guide.jpg', 
   '2024-01-16 14:30:00', 
   78, 
   456, 
   'Tips on module mapping and application timeline',
   '[
     {
       "id": 3,
       "user_id": 1,
       "content": "This guide is super helpful! Could you add more details about visa application?",
       "created_at": "2024-01-16 15:00:00"
     }
   ]'::jsonb
);