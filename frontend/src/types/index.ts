export interface User {
    id: string
    name: string
    email: string
    yearly_goal: number
    created_at: string
}

export interface Book {
    id: string
    user_id: string
    title: string
    author: string
    cover_url: string
    total_pages: number
    status: 'want_to_read' | 'reading' | 'finished'
    rating: number
    finished_at: string | null
    created_at: string
}

export interface ReadingSession {
    id: string
    book_id: string
    user_id: string
    started_at: string
    ended_at: string | null
    duration_seconds: number | null
    pages_read: number
    status: 'active' | 'paused' | 'completed'
    created_at: string
}

export interface ReadingProgress {
    id: string
    book_id: string
    user_id: string
    current_page: number
    updated_at: string
}

export interface SessionNote {
    id: string
    session_id: string
    user_id: string
    chapter: string
    pages: string
    note: string
    created_at: string
}

export interface StatsData {
    total_books: number
    finished_books: number
    total_reading_time_seconds: number
    total_pages_read: number
    current_streak: number
    yearly_goal: number
    yearly_finished: number
}

export interface DailyReading {
    date: string
    pages: number
    seconds: number
}

export interface AuthResponse {
    token: string
    user: User
}

export interface CreateBookRequest {
    title: string
    author: string
    cover_url: string
    total_pages: number
    status: 'want_to_read' | 'reading' | 'finished'
}

export interface UpdateBookRequest {
    title: string
    author: string
    cover_url: string
    total_pages: number
    status: 'want_to_read' | 'reading' | 'finished'
    rating: number
}

export interface GoogleBook {
    id: string
    volumeInfo: {
        title: string
        authors?: string[]
        imageLinks?: {
            thumbnail: string
        }
        pageCount?: number
    }
}
