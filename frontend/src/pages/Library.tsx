import { useState, useEffect, useMemo } from 'react'
import { Link } from 'react-router-dom'
import Layout from '../components/Layout'
import PageState from '../components/PageState'
import api from '../api/axios'
import type { Book } from '../types'
import { BookOpen, Timer, Search, Plus, BookMarked, Check } from 'lucide-react'
import { statusLabel, statusClassName } from '../utils/bookUtils'

type StatusFilter = 'all' | 'reading' | 'want_to_read' | 'finished'

const Library = () => {
    const [books, setBooks] = useState<Book[]>([])
    const [loading, setLoading] = useState(true)
    const [filter, setFilter] = useState<StatusFilter>('all')
    const [search, setSearch] = useState('')
    const [error, setError] = useState(false)

    const fetchBooks = async () => {
        try {
            setError(false)
            const res = await api.get('/api/books')
            setBooks(res.data.data || [])
        } catch (err) {
            console.error(err)
            setError(true)
        } finally {
            setLoading(false)
        }
    }

    useEffect(() => {
        // eslint-disable-next-line react-hooks/set-state-in-effect
        fetchBooks()
    }, [])

    const filtered = useMemo(
        () =>
            books.filter((book) => {
                const matchesFilter = filter === 'all' || book.status === filter
                const matchesSearch =
                    book.title.toLowerCase().includes(search.toLowerCase()) ||
                    book.author.toLowerCase().includes(search.toLowerCase())
                return matchesFilter && matchesSearch
            }),
        [books, filter, search],
    )

    const counts = useMemo(
        () => ({
            all: books.length,
            reading: books.filter((b) => b.status === 'reading').length,
            want_to_read: books.filter((b) => b.status === 'want_to_read').length,
            finished: books.filter((b) => b.status === 'finished').length,
        }),
        [books],
    )

    const filters: { key: StatusFilter; label: string }[] = [
        { key: 'all', label: `All (${counts.all})` },
        { key: 'reading', label: `Reading (${counts.reading})` },
        { key: 'want_to_read', label: `Want to Read (${counts.want_to_read})` },
        { key: 'finished', label: `Finished (${counts.finished})` },
    ]

    if (loading || error) {
        return (
            <PageState
                loading={loading}
                error={error}
                loadingText="Loading your library..."
                onRetry={fetchBooks}
            />
        )
    }

    return (
        <Layout>
            <div className="mb-8 flex items-start justify-between">
                <div>
                    <h1 className="text-bt-dark mb-1 font-serif text-2xl font-semibold">
                        My Library
                    </h1>
                    <p className="text-bt-muted text-sm">
                        {books.length} book{books.length !== 1 ? 's' : ''} in your collection
                    </p>
                </div>
                <Link
                    to="/books/search"
                    className="bg-bt-dark text-bt-cream hover:bg-bt-gold flex cursor-pointer items-center gap-2 rounded-md px-4 py-2 text-sm font-medium transition-colors duration-200"
                >
                    <Plus size={14} />
                    Add Book
                </Link>
            </div>

            <div className="mb-6 flex flex-col gap-4 md:flex-row md:items-center">
                <div className="bg-bt-accent-bg flex flex-wrap gap-1 rounded-lg p-1">
                    {filters.map((f) => (
                        <button
                            key={f.key}
                            onClick={() => setFilter(f.key)}
                            className={`cursor-pointer rounded-md px-3 py-1.5 text-xs font-medium transition-colors duration-200 ${
                                filter === f.key
                                    ? 'bg-bt-surface text-bt-dark shadow-sm'
                                    : 'text-bt-muted bg-transparent'
                            }`}
                        >
                            {f.label}
                        </button>
                    ))}
                </div>
                <div className="relative md:ml-auto">
                    <Search
                        size={13}
                        className="text-bt-muted-light absolute top-1/2 left-3 -translate-y-1/2"
                    />
                    <input
                        type="text"
                        value={search}
                        onChange={(e) => setSearch(e.target.value)}
                        placeholder="Search your library..."
                        className="input-field w-56 py-2 pr-4 pl-8 text-sm transition-colors duration-200"
                    />
                </div>
            </div>

            {filtered.length === 0 ? (
                <div className="border-bt-dashed bg-bt-surface rounded-xl border border-dashed py-16 text-center">
                    <BookMarked
                        size={28}
                        strokeWidth={1.5}
                        className="text-bt-placeholder mx-auto mb-3"
                    />
                    <p className="text-bt-muted mb-1 text-sm">
                        {search ? `No books matching "${search}"` : 'No books here yet'}
                    </p>
                    <Link
                        to="/books/search"
                        className="text-bt-gold hover:text-bt-muted-dark text-xs font-medium transition-colors"
                    >
                        Search and add books →
                    </Link>
                </div>
            ) : (
                <div className="grid grid-cols-1 gap-3 md:grid-cols-2 lg:grid-cols-3">
                    {filtered.map((book) => (
                        <Link
                            key={book.id}
                            to={`/books/${book.id}`}
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
                            <div className="flex min-w-0 flex-1 flex-col justify-between">
                                <div>
                                    <p className="text-bt-dark mb-1 line-clamp-2 text-sm leading-snug font-medium">
                                        {book.title}
                                    </p>
                                    <p className="text-bt-muted text-xs">{book.author}</p>
                                </div>
                                <div className="mt-3 flex items-center justify-between">
                                    <span
                                        className={`rounded-full px-2 py-0.5 text-[11px] ${statusClassName(book.status)}`}
                                    >
                                        {statusLabel(book.status)}
                                    </span>
                                    {book.status === 'reading' && (
                                        <Link
                                            to={`/reading/${book.id}`}
                                            onClick={(e) => e.stopPropagation()}
                                            className="text-bt-gold hover:text-bt-muted-dark flex items-center gap-1 text-xs font-medium transition-colors"
                                        >
                                            <Timer size={11} /> Read
                                        </Link>
                                    )}
                                    {book.status === 'finished' && (
                                        <Check size={13} className="text-bt-status-finished-text" />
                                    )}
                                </div>
                            </div>
                        </Link>
                    ))}
                </div>
            )}
        </Layout>
    )
}

export default Library
