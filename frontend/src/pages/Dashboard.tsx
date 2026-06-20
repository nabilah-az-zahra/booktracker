import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import Layout from '../components/Layout'
import PageState from '../components/PageState'
import { useAuth } from '../context/useAuth'
import api from '../api/axios'
import type { Book, StatsData, ReadingSession } from '../types'
import {
    BookOpen,
    Clock,
    Target,
    Flame,
    ArrowRight,
    Plus,
    Timer,
    BookMarked,
    Search,
} from 'lucide-react'
import { statusBorderClassName, statusLabel, statusTextClassName } from '../utils/bookUtils'
import { formatTime } from '../utils/formatUtils'

const Dashboard = () => {
    const { user } = useAuth()
    const [books, setBooks] = useState<Book[]>([])
    const [stats, setStats] = useState<StatsData | null>(null)
    const [activeSession, setActiveSession] = useState<ReadingSession | null>(null)
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState(false)

    const fetchData = async () => {
        try {
            setError(false)
            const [booksRes, statsRes] = await Promise.all([
                api.get('/api/books'),
                api.get('/api/stats'),
            ])
            setBooks(booksRes.data.data || [])
            setStats(statsRes.data.data)
        } catch (err) {
            console.error(err)
            setError(true)
        } finally {
            setLoading(false)
        }

        try {
            const sessionRes = await api.get('/api/sessions/active')
            setActiveSession(sessionRes.data.data)
        } catch {
            //
        }
    }

    useEffect(() => {
        // eslint-disable-next-line react-hooks/set-state-in-effect
        fetchData()
    }, [])

    const currentlyReading = books.filter((b) => b.status === 'reading')

    const now = new Date()
    const hours = now.getHours()
    const greeting = hours < 12 ? 'morning' : hours < 18 ? 'afternoon' : 'evening'
    const dateString = now.toLocaleDateString('en-US', {
        weekday: 'long',
        month: 'long',
        day: 'numeric',
    })

    if (loading || error) {
        return (
            <PageState
                loading={loading}
                error={error}
                loadingText="Loading your library..."
                onRetry={fetchData}
            />
        )
    }

    return (
        <Layout>
            <div className="mb-8 flex items-start justify-between">
                <div>
                    <h1 className="text-bt-dark mb-1 font-serif text-2xl font-semibold">
                        Good {greeting}, {user?.name?.split(' ')[0]}
                    </h1>
                    <p className="text-bt-muted text-sm">{dateString}</p>
                </div>
                <Link
                    to="/books/search"
                    className="bg-bt-dark text-bt-cream hover:bg-bt-gold flex cursor-pointer items-center gap-2 rounded-md px-4 py-2 text-sm font-medium transition-colors duration-200"
                >
                    <Plus size={14} />
                    Add Book
                </Link>
            </div>

            <div className="space-y-8">
                {activeSession && (
                    <div className="border-bt-banner-border border-l-bt-gold bg-bt-banner-bg flex items-center justify-between rounded-lg border border-l-[3px] px-5 py-4">
                        <div className="flex items-center gap-3">
                            <Timer size={16} className="text-bt-gold" />
                            <div>
                                <p className="text-bt-dark text-sm font-medium">
                                    Reading session in progress
                                </p>
                                <p className="text-bt-muted mt-0.5 text-xs">
                                    Tap to continue where you left off
                                </p>
                            </div>
                        </div>
                        <Link
                            to={`/reading/${activeSession.book_id}`}
                            className="text-bt-gold hover:text-bt-muted-dark flex items-center gap-1 text-xs font-medium transition-colors"
                        >
                            Resume <ArrowRight size={12} />
                        </Link>
                    </div>
                )}

                {stats && (
                    <div className="grid grid-cols-2 gap-3 md:grid-cols-4">
                        <div className="border-bt-border bg-bt-surface hover:border-bt-gold rounded-xl border px-5 py-4 transition-colors duration-200">
                            <BookOpen size={16} strokeWidth={1.5} className="text-bt-gold mb-3" />
                            <p className="text-bt-dark mb-0.5 font-serif text-2xl font-semibold">
                                {stats.finished_books}
                            </p>
                            <p className="text-bt-warm text-xs font-medium">Books Read</p>
                            <p className="text-bt-muted-light mt-0.5 text-xs">
                                of {stats.total_books} total
                            </p>
                        </div>
                        <div className="border-bt-border bg-bt-surface hover:border-bt-gold rounded-xl border px-5 py-4 transition-colors duration-200">
                            <Clock size={16} strokeWidth={1.5} className="text-bt-gold mb-3" />
                            <p className="text-bt-dark mb-0.5 font-serif text-2xl font-semibold">
                                {formatTime(stats.total_reading_time_seconds)}
                            </p>
                            <p className="text-bt-warm text-xs font-medium">Reading Time</p>
                            <p className="text-bt-muted-light mt-0.5 text-xs">
                                {stats.total_pages_read} pages
                            </p>
                        </div>
                        <div className="border-bt-border bg-bt-surface hover:border-bt-gold rounded-xl border px-5 py-4 transition-colors duration-200">
                            <Flame size={16} strokeWidth={1.5} className="text-bt-gold mb-3" />
                            <p className="text-bt-dark mb-0.5 font-serif text-2xl font-semibold">
                                {stats.current_streak}
                            </p>
                            <p className="text-bt-warm text-xs font-medium">Day Streak</p>
                            <p className="text-bt-muted-light mt-0.5 text-xs">days in a row</p>
                        </div>
                        <div className="border-bt-border bg-bt-surface hover:border-bt-gold rounded-xl border px-5 py-4 transition-colors duration-200">
                            <Target size={16} strokeWidth={1.5} className="text-bt-gold mb-3" />
                            <p className="text-bt-dark mb-0.5 font-serif text-2xl font-semibold">
                                {stats.yearly_finished}/{stats.yearly_goal || '—'}
                            </p>
                            <p className="text-bt-warm text-xs font-medium">Yearly Goal</p>
                            <p className="text-bt-muted-light mt-0.5 text-xs">
                                {stats.yearly_goal
                                    ? `${Math.round((stats.yearly_finished / stats.yearly_goal) * 100)}% done`
                                    : 'No goal set'}
                            </p>
                        </div>
                    </div>
                )}

                <div>
                    <div className="mb-4 flex items-center justify-between">
                        <h2 className="text-bt-dark font-serif text-lg font-semibold">
                            Currently Reading
                        </h2>
                        <Link
                            to="/library"
                            className="text-bt-muted hover:text-bt-gold flex items-center gap-1 text-xs transition-colors"
                        >
                            All books <ArrowRight size={12} />
                        </Link>
                    </div>

                    {currentlyReading.length === 0 ? (
                        <div className="border-bt-dashed bg-bt-surface rounded-xl border border-dashed py-12 text-center">
                            <BookMarked
                                size={28}
                                strokeWidth={1.5}
                                className="text-bt-placeholder mx-auto mb-3"
                            />
                            <p className="text-bt-muted mb-3 text-sm">Nothing in progress yet</p>
                            <Link
                                to="/books/search"
                                className="text-bt-gold hover:text-bt-muted-dark text-xs font-medium transition-colors"
                            >
                                Find something to read{' '}
                                <Search
                                    size={13}
                                    strokeWidth={2}
                                    className="ml-1 inline-block align-text-bottom"
                                />
                            </Link>
                        </div>
                    ) : (
                        <div className="grid grid-cols-1 gap-3 md:grid-cols-2 lg:grid-cols-3">
                            {currentlyReading.map((book) => (
                                <div
                                    key={book.id}
                                    className="border-bt-border bg-bt-surface hover:border-bt-gold flex gap-4 rounded-xl border p-4 transition-colors duration-200"
                                >
                                    {book.cover_url ? (
                                        <img
                                            src={book.cover_url}
                                            alt={book.title}
                                            className="h-16 w-12 shrink-0 rounded object-cover"
                                        />
                                    ) : (
                                        <div className="bg-bt-accent-bg flex h-16 w-12 shrink-0 items-center justify-center rounded">
                                            <BookOpen
                                                size={16}
                                                className="text-bt-gold"
                                                strokeWidth={1.5}
                                            />
                                        </div>
                                    )}
                                    <div className="min-w-0 flex-1">
                                        <p className="text-bt-dark mb-1 line-clamp-2 text-sm leading-snug font-medium">
                                            {book.title}
                                        </p>
                                        <p className="text-bt-muted mb-4 text-xs">{book.author}</p>
                                        <div className="flex items-center gap-3">
                                            <Link
                                                to={`/reading/${book.id}`}
                                                className="text-bt-gold hover:text-bt-muted-dark flex items-center gap-1.5 text-xs font-medium transition-colors"
                                            >
                                                <Timer size={11} /> Read
                                            </Link>
                                            <Link
                                                to={`/books/${book.id}`}
                                                className="text-bt-muted-light hover:text-bt-warm text-xs transition-colors"
                                            >
                                                Details
                                            </Link>
                                        </div>
                                    </div>
                                </div>
                            ))}
                        </div>
                    )}
                </div>

                <div>
                    <h2 className="text-bt-dark mb-4 font-serif text-lg font-semibold">
                        Recently Added
                    </h2>
                    {books.length === 0 ? (
                        <div className="border-bt-dashed bg-bt-surface flex flex-col items-center justify-center rounded-xl border-2 border-dashed px-4 py-14 text-center">
                            <div className="bg-bt-accent-bg mb-3 flex h-10 w-10 items-center justify-center rounded-full">
                                <BookOpen size={16} className="text-bt-gold" strokeWidth={1.5} />
                            </div>
                            <p className="text-bt-dark font-serif text-sm font-medium tracking-tight">
                                Your library is silent
                            </p>
                            <p className="text-bt-muted mt-0.5 text-xs">
                                Books you actively track will appear right here.
                            </p>
                        </div>
                    ) : (
                        <div className="border-bt-border bg-bt-surface overflow-hidden rounded-xl border shadow-xs">
                            {books.slice(0, 5).map((book) => (
                                <Link
                                    key={book.id}
                                    to={`/books/${book.id}`}
                                    className={`group border-bt-border bg-bt-surface flex items-center gap-4 border-b border-l-2 px-5 py-4 transition-all duration-200 last:border-b-0 ${statusBorderClassName(book.status)}`}
                                >
                                    {book.cover_url ? (
                                        <img
                                            src={book.cover_url}
                                            alt={book.title}
                                            className="border-bt-border/50 h-11 w-8 shrink-0 rounded-sm border object-cover shadow-xs transition-transform duration-200 group-hover:scale-102"
                                        />
                                    ) : (
                                        <div className="bg-bt-accent-bg border-bt-dashed flex h-11 w-8 shrink-0 items-center justify-center rounded-sm border">
                                            <BookOpen
                                                size={12}
                                                className="text-bt-placeholder group-hover:text-bt-gold transition-colors duration-200"
                                                strokeWidth={1.5}
                                            />
                                        </div>
                                    )}
                                    <div className="min-w-0 flex-1">
                                        <p className="text-bt-dark group-hover:text-bt-gold truncate font-serif text-sm font-medium tracking-tight transition-colors duration-200">
                                            {book.title}
                                        </p>
                                        <p className="text-bt-muted mt-0.5 truncate text-xs">
                                            by {book.author} ·{' '}
                                            <span className={statusTextClassName(book.status)}>
                                                {statusLabel(book.status)}
                                            </span>
                                        </p>
                                    </div>

                                    <ArrowRight
                                        size={14}
                                        strokeWidth={2}
                                        className="text-bt-placeholder group-hover:text-bt-dark shrink-0 transform transition-all duration-200 group-hover:translate-x-0.5"
                                    />
                                </Link>
                            ))}
                        </div>
                    )}
                </div>
            </div>
        </Layout>
    )
}

export default Dashboard
