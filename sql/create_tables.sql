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

-- Recreate users table with UUID
CREATE TABLE IF NOT EXISTS users (
   id TEXT PRIMARY KEY,
   email VARCHAR(255) UNIQUE NOT NULL,
   password_hash VARCHAR(255) NOT NULL,
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

-- Insert users with mock UUIDs
INSERT INTO users (id, email, password_hash) VALUES
('123e4567-e89b-12d3-a456-426614174000', 'testing@gmail.com', '$2a$10$mN6CaIxk7mU0QM3B2Q490euGHJS5Dx0AOTjG7v82f9dQL/Gm.gCEC'),
('123e4567-e89b-12d3-a456-426614174001', 'john@example.com', '$2a$10$mN6CaIxk7mU0QM3B2Q490euGHJS5Dx0AOTjG7v82f9dQL/Gm.gCEC'),
('123e4567-e89b-12d3-a456-426614174002', 'jane@example.com', '$2a$10$mN6CaIxk7mU0QM3B2Q490euGHJS5Dx0AOTjG7v82f9dQL/Gm.gCEC');

-- Insert posts with mock UUIDs, tags, and updated_at
INSERT INTO posts (id, user_id, title, content, created_at, updated_at, likes_count, views_count, comments, tags) VALUES 
('123e4567-e89b-12d3-a456-426614174003', '123e4567-e89b-12d3-a456-426614174000', 'Exchange Experience in Tokyo', 
  'Sharing my amazing semester abroad experience at Waseda University!', 
  '2024-01-15T10:00:00Z',
  '2024-01-15T10:00:00Z',  -- Added updated_at
  45, 
  230, 
  '[
    {
      "id": "123e4567-e89b-12d3-a456-426614174006",
      "user_id": "123e4567-e89b-12d3-a456-426614174001",
      "content": "Great post! How did you handle the language barrier?",
      "created_at": "2024-01-15T10:30:00Z"
    },
    {
      "id": "123e4567-e89b-12d3-a456-426614174007",
      "user_id": "123e4567-e89b-12d3-a456-426614174002",
      "content": "The campus looks amazing! Did you stay in university accommodation?",
      "created_at": "2024-01-15T11:00:00Z"
    }
  ]'::jsonb,
  ARRAY['123e4567-e89b-12d3-a456-426614174011']  -- Summer Exchange tag
),
('123e4567-e89b-12d3-a456-426614174004', '123e4567-e89b-12d3-a456-426614174001', 'NUS Exchange Application Guide', 
  'Step-by-step guide on applying for exchange at NUS', 
  '2024-01-16T14:30:00Z',
  '2024-01-16T14:30:00Z',  -- Added updated_at
  78, 
  456, 
  '[
    {
      "id": "123e4567-e89b-12d3-a456-426614174008",
      "user_id": "123e4567-e89b-12d3-a456-426614174000",
      "content": "This guide is super helpful! Could you add more details about visa application?",
      "created_at": "2024-01-16T15:00:00Z"
    }
  ]'::jsonb,
  ARRAY['123e4567-e89b-12d3-a456-426614174010', '123e4567-e89b-12d3-a456-426614174011']  -- Europe and Summer Exchange tags
);