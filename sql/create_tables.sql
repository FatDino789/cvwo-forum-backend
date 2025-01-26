-- Drop existing tables
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS tags;

-- Create tags table
CREATE TABLE IF NOT EXISTS tags (
   id TEXT PRIMARY KEY,
   text VARCHAR(50) NOT NULL UNIQUE,
   color VARCHAR(7) NOT NULL,
   searches INTEGER DEFAULT 0
);

-- Recreate users table with UUID and username
CREATE TABLE IF NOT EXISTS users (
  id TEXT PRIMARY KEY,
  username VARCHAR(50) UNIQUE NOT NULL,
  email VARCHAR(255) UNIQUE NOT NULL, 
  password_hash VARCHAR(255) NOT NULL,
  icon_index INTEGER NOT NULL,
  color_index INTEGER NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create posts table with UUID and tags array
CREATE TABLE IF NOT EXISTS posts (
   id TEXT PRIMARY KEY,
   user_id TEXT REFERENCES users(id),
   title VARCHAR(255) NOT NULL,
   content TEXT NOT NULL,
   created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
   likes_count INTEGER DEFAULT 0,
   views_count INTEGER DEFAULT 0,
   comments JSONB DEFAULT '[]',
   tags TEXT[] DEFAULT '{}'
);

-- Insert tags with mock UUIDs
INSERT INTO tags (id, text, color, searches) VALUES
('123e4567-e89b-12d3-a456-426614174010', 'Europe', '#DCF2E7', 10),
('123e4567-e89b-12d3-a456-426614174011', 'Summer Exchange', '#FFEDD5', 5);

-- Insert users with mock UUIDs and usernames
INSERT INTO users (id, username, email, password_hash, icon_index, color_index) VALUES
('123e4567-e89b-12d3-a456-426614174000', 'testuser123', 'testing@gmail.com', '$2a$10$mN6CaIxk7mU0QM3B2Q490euGHJS5Dx0AOTjG7v82f9dQL/Gm.gCEC', 45, 2),
('123e4567-e89b-12d3-a456-426614174001', 'johndoe', 'john@example.com', '$2a$10$mN6CaIxk7mU0QM3B2Q490euGHJS5Dx0AOTjG7v82f9dQL/Gm.gCEC', 12, 4),
('123e4567-e89b-12d3-a456-426614174002', 'janesmith', 'jane@example.com', '$2a$10$mN6CaIxk7mU0QM3B2Q490euGHJS5Dx0AOTjG7v82f9dQL/Gm.gCEC', 33, 1);

INSERT INTO posts VALUES 
('123e4567-e89b-12d3-a456-426614174003', '123e4567-e89b-12d3-a456-426614174000', 'Exchange Experience in Tokyo', 
 'Sharing my amazing semester abroad experience at Waseda University!', 
 '2024-01-15T10:00:00Z', '2024-01-15T10:00:00Z', 45, 230, 
 '[
   {
     "id": "123e4567-e89b-12d3-a456-426614174006",
     "user_id": "123e4567-e89b-12d3-a456-426614174001",
     "content": "Great post! How did you handle the language barrier?",
     "created_at": "2024-01-15T10:30:00Z",
     "username": "johndoe",
     "icon_index": 12,
     "color_index": 4
   }
 ]'::jsonb,
 ARRAY['123e4567-e89b-12d3-a456-426614174011']
),
('123e4567-e89b-12d3-a456-426614174004', '123e4567-e89b-12d3-a456-426614174001', 'NUS Exchange Application Guide', 
 'Step-by-step guide on applying for exchange at NUS', 
 '2024-01-16T14:30:00Z', '2024-01-16T14:30:00Z', 78, 456,
 '[
   {
     "id": "123e4567-e89b-12d3-a456-426614174008",
     "user_id": "123e4567-e89b-12d3-a456-426614174000",
     "content": "This guide is super helpful! Could you add more details about visa application?",
     "created_at": "2024-01-16T15:00:00Z",
     "username": "testuser123",
     "icon_index": 45,
     "color_index": 2
   }
 ]'::jsonb,
 ARRAY['123e4567-e89b-12d3-a456-426614174010', '123e4567-e89b-12d3-a456-426614174011']
);