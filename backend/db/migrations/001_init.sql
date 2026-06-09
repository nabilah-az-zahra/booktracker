CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Users
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    yearly_goal INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Books
CREATE TABLE IF NOT EXISTS books (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255),
    cover_url TEXT,
    total_pages INT DEFAULT 0,
    status VARCHAR(20) DEFAULT 'want_to_read' CHECK (status IN ('want_to_read', 'reading', 'finished')),
    rating INT DEFAULT 0 CHECK (rating >= 0 AND rating <= 5),
    finished_at TIMESTAMP,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Reading Sessions
CREATE TABLE IF NOT EXISTS reading_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    book_id UUID REFERENCES books(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    started_at TIMESTAMP NOT NULL,
    ended_at TIMESTAMP,
    duration_seconds INT DEFAULT 0,
    pages_read INT DEFAULT 0,
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'paused', 'completed')),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Reading Progress
CREATE TABLE IF NOT EXISTS reading_progress (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    book_id UUID REFERENCES books(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    current_page INT DEFAULT 0,
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT unique_progress_per_book_user UNIQUE (book_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_books_user_id ON books(user_id);
CREATE INDEX IF NOT EXISTS idx_reading_sessions_book_id ON reading_sessions(book_id);
CREATE INDEX IF NOT EXISTS idx_reading_sessions_user_id ON reading_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_reading_progress_book_id ON reading_progress(book_id);
CREATE INDEX IF NOT EXISTS idx_reading_progress_user_id ON reading_progress(user_id);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);