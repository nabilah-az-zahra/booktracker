import { useState, useEffect } from 'react'
import { useParams, useNavigate, Link } from 'react-router-dom'
import Layout from '../components/Layout'
import PageState from '../components/PageState'
import api from '../api/axios'
import type { Book, ReadingSession, ReadingProgress } from '../types'
import {
    BookOpen,
    Timer,
    Trash2,
    ArrowLeft,
    Clock,
    BookMarked,
    Star,
    Check,
    ChevronDown,
} from 'lucide-react'
import { statusClassName } from '../utils/bookUtils'
import { formatTime } from '../utils/formatUtils'

const formatDate = (dateStr: string) =>
    new Date(dateStr).toLocaleDateString('en-US', {
        month: 'short',
        day: 'numeric',
        year: 'numeric',
    })

const BookDetail = () => {
    const { id } = useParams<{ id: string }>()
    const navigate = useNavigate()
    const [book, setBook] = useState<Book | null>(null)
    const [sessions, setSessions] = useState<ReadingSession[]>([])
    const [progress, setProgress] = useState<ReadingProgress | null>(null)
    const [loading, setLoading] = useState(true)
    const [showDeleteModal, setShowDeleteModal] = useState(false)
    const [showStatusDropdown, setShowStatusDropdown] = useState(false)
    const [deleteError, setDeleteError] = useState('')
    const [isDeleting, setIsDeleting] = useState(false)
    const [updateError, setUpdateError] = useState('')

    useEffect(() => {
        const fetchData = async () => {
            setLoading(true)
            setBook(null)
            setSessions([])
            setProgress(null)
            try {
                const bookRes = await api.get(`/api/books/${id}`)
                if (bookRes?.data?.data) {
                    setBook(bookRes.data.data)
                }

                const sessionsRes = await api
                    .get(`/api/sessions/book/${id}`)
                    .catch(() => ({ data: { data: [] } }))
                setSessions(sessionsRes.data.data || [])

                const progressRes = await api
                    .get(`/api/progress/${id}`)
                    .catch(() => ({ data: { data: null } }))
                setProgress(progressRes.data.data)
            } catch (err) {
                console.error('BookDetail error:', err)
            } finally {
                setLoading(false)
            }
        }
        if (id) fetchData()
    }, [id])

    const handleBookUpdate = async (patch: Partial<Pick<Book, 'status' | 'rating'>>) => {
        if (!book) return
        const previous = book
        setUpdateError('')
        setShowStatusDropdown(false)
        setBook({ ...book, ...patch })
        try {
            const res = await api.put(`/api/books/${id}`, {
                title: book.title,
                author: book.author,
                cover_url: book.cover_url,
                total_pages: book.total_pages,
                status: book.status,
                rating: book.rating,
                ...patch,
            })
            setBook(res.data.data)
        } catch (err) {
            console.error(err)
            setBook(previous)
            setUpdateError('Failed to update book. Please try again.')
        }
    }

    const handleDelete = async () => {
        setDeleteError('')
        setIsDeleting(true)
        try {
            await api.delete(`/api/books/${id}`)
            navigate('/library')
        } catch (err) {
            console.error(err)
            setDeleteError('Failed to remove book. Please try again.')
        } finally {
            setIsDeleting(false)
        }
    }

    const totalReadingTime = sessions.reduce((sum, s) => sum + (s.duration_seconds ?? 0), 0)
    const totalPagesRead = progress?.current_page || 0
    const progressPercent = book?.total_pages
        ? Math.min(Math.round((totalPagesRead / book.total_pages) * 100), 100)
        : 0

    const statusOptions = [
        { value: 'want_to_read', label: 'Want to Read' },
        { value: 'reading', label: 'Currently Reading' },
        { value: 'finished', label: 'Finished' },
    ]

    if (loading) {
        return <PageState loading loadingText="Loading book..." />
    }

    if (!book) {
        return (
            <Layout>
                <div className="py-20 text-center">
                    <p className="text-bt-muted mb-3 text-sm">Book not found</p>
                    <Link
                        to="/library"
                        className="text-bt-gold hover:text-bt-muted-dark text-xs transition-colors"
                    >
                        ← Back to library
                    </Link>
                </div>
            </Layout>
        )
    }

    return (
        <Layout>
            <button
                onClick={() => navigate('/library')}
                className="text-bt-muted hover:text-bt-dark mb-6 flex cursor-pointer items-center gap-1.5 text-sm transition-colors"
            >
                <ArrowLeft size={14} /> Back to Library
            </button>

            {updateError && (
                <div className="bg-bt-danger-bg text-bt-danger mb-6 flex items-start justify-between rounded-lg px-4 py-3 text-xs font-medium">
                    <span>{updateError}</span>
                    <button
                        onClick={() => setUpdateError('')}
                        className="cursor-pointer text-sm leading-none font-bold hover:opacity-70"
                    >
                        ×
                    </button>
                </div>
            )}

            <div className="grid grid-cols-1 gap-8 lg:grid-cols-3">
                <div className="lg:col-span-1">
                    <div className="bg-bt-surface border-bt-border sticky top-8 rounded-xl border p-6">
                        <div className="mb-6 flex justify-center">
                            {book.cover_url ? (
                                <img
                                    src={book.cover_url}
                                    alt={book.title}
                                    className="h-44 w-32 rounded-md object-cover shadow-md"
                                />
                            ) : (
                                <div className="bg-bt-accent-bg flex h-44 w-32 items-center justify-center rounded-md">
                                    <BookOpen
                                        size={32}
                                        className="text-bt-gold"
                                        strokeWidth={1.5}
                                    />
                                </div>
                            )}
                        </div>

                        <h1 className="text-bt-dark mb-1 text-center font-serif text-lg leading-snug font-semibold">
                            {book.title}
                        </h1>
                        <p className="text-bt-muted mb-4 text-center text-sm">{book.author}</p>

                        <div className="relative mb-4">
                            <button
                                onClick={() => setShowStatusDropdown(!showStatusDropdown)}
                                className={`flex w-full cursor-pointer items-center justify-between rounded-lg border border-transparent px-3 py-2 text-sm transition-colors duration-200 ${statusClassName(book.status)}`}
                            >
                                <span className="font-medium">
                                    {statusOptions.find((s) => s.value === book.status)?.label}
                                </span>
                                <ChevronDown
                                    size={14}
                                    className={`transition-transform duration-200 ${showStatusDropdown ? 'rotate-180' : ''}`}
                                />
                            </button>
                            {showStatusDropdown && (
                                <>
                                    <div
                                        className="fixed inset-0 z-5"
                                        onClick={() => setShowStatusDropdown(false)}
                                    />
                                    <div className="bg-bt-surface border-bt-border absolute top-full right-0 left-0 z-10 mt-1 rounded-lg border py-1 shadow-md">
                                        {statusOptions.map((opt) => (
                                            <button
                                                key={opt.value}
                                                onClick={() =>
                                                    handleBookUpdate({
                                                        status: opt.value as Book['status'],
                                                    })
                                                }
                                                className="text-bt-dark hover:bg-bt-cream flex w-full cursor-pointer items-center justify-between px-4 py-2 text-left text-sm transition-colors"
                                            >
                                                {opt.label}
                                                {book.status === opt.value && (
                                                    <Check size={13} className="text-bt-gold" />
                                                )}
                                            </button>
                                        ))}
                                    </div>
                                </>
                            )}
                        </div>

                        <div className="mb-5 flex items-center justify-center gap-1">
                            {[1, 2, 3, 4, 5].map((star) => (
                                <button
                                    key={star}
                                    onClick={() => handleBookUpdate({ rating: star })}
                                    className="cursor-pointer transition-transform duration-150 hover:scale-110"
                                >
                                    <Star
                                        size={18}
                                        strokeWidth={1.5}
                                        className={
                                            star <= book.rating ? 'text-bt-gold' : 'text-bt-dashed'
                                        }
                                        style={{
                                            fill:
                                                star <= book.rating
                                                    ? 'var(--color-bt-gold)'
                                                    : 'none',
                                        }}
                                    />
                                </button>
                            ))}
                        </div>

                        {book.total_pages > 0 && (
                            <div className="mb-5">
                                <div className="text-bt-muted mb-2 flex items-center justify-between text-xs">
                                    <span>Progress</span>
                                    <span>
                                        {totalPagesRead} / {book.total_pages} pages
                                    </span>
                                </div>
                                <div className="bg-bt-track h-1.5 w-full rounded-full">
                                    <div
                                        className="bg-bt-gold h-full rounded-full transition-all duration-500"
                                        style={{ width: `${progressPercent}%` }}
                                    />
                                </div>
                                <p className="text-bt-muted-light mt-1 text-right text-xs">
                                    {progressPercent}%
                                </p>
                            </div>
                        )}

                        {book.status === 'reading' && (
                            <Link
                                to={`/reading/${book.id}`}
                                className="bg-bt-dark text-bt-cream hover:bg-bt-gold mb-3 flex w-full items-center justify-center gap-2 rounded-lg py-2.5 text-sm font-medium transition-colors duration-200"
                            >
                                <Timer size={14} /> Start Reading
                            </Link>
                        )}

                        <button
                            onClick={() => setShowDeleteModal(true)}
                            className="text-bt-subtle hover:text-bt-error-border flex w-full cursor-pointer items-center justify-center gap-2 py-2 text-xs transition-colors"
                        >
                            <Trash2 size={13} /> Remove from library
                        </button>
                    </div>
                </div>

                <div className="space-y-6 lg:col-span-2">
                    <div className="grid grid-cols-3 gap-3">
                        <div className="bg-bt-surface border-bt-border rounded-xl border px-4 py-4">
                            <Clock size={15} strokeWidth={1.5} className="text-bt-gold mb-2" />
                            <p className="text-bt-dark mb-0.5 font-serif text-xl font-semibold">
                                {formatTime(totalReadingTime)}
                            </p>
                            <p className="text-bt-muted text-xs">Time Read</p>
                        </div>
                        <div className="bg-bt-surface border-bt-border rounded-xl border px-4 py-4">
                            <BookMarked size={15} strokeWidth={1.5} className="text-bt-gold mb-2" />
                            <p className="text-bt-dark mb-0.5 font-serif text-xl font-semibold">
                                {sessions.length}
                            </p>
                            <p className="text-bt-muted text-xs">Sessions</p>
                        </div>
                        <div className="bg-bt-surface border-bt-border rounded-xl border px-4 py-4">
                            <BookOpen size={15} strokeWidth={1.5} className="text-bt-gold mb-2" />
                            <p className="text-bt-dark mb-0.5 font-serif text-xl font-semibold">
                                {totalPagesRead}
                            </p>
                            <p className="text-bt-muted text-xs">Pages Read</p>
                        </div>
                    </div>

                    <div className="bg-bt-surface border-bt-border overflow-hidden rounded-xl border">
                        <div className="border-bt-border border-b px-6 py-4">
                            <h2 className="text-bt-dark font-serif text-base font-semibold">
                                Reading Sessions
                            </h2>
                        </div>

                        {sessions.length === 0 ? (
                            <div className="py-12 text-center">
                                <Timer
                                    size={24}
                                    strokeWidth={1.5}
                                    className="text-bt-placeholder mx-auto mb-3"
                                />
                                <p className="text-bt-muted text-sm">No reading sessions yet</p>
                                {book.status === 'reading' && (
                                    <Link
                                        to={`/reading/${book.id}`}
                                        className="text-bt-gold hover:text-bt-muted-dark mt-2 inline-block text-xs font-medium transition-colors"
                                    >
                                        Start your first session →
                                    </Link>
                                )}
                            </div>
                        ) : (
                            <div className="divide-bt-border divide-y">
                                {sessions.map((session, index) => (
                                    <div
                                        key={session.id}
                                        className="flex items-center justify-between px-6 py-4"
                                    >
                                        <div className="flex items-center gap-4">
                                            <div className="bg-bt-accent-bg text-bt-gold flex h-8 w-8 shrink-0 items-center justify-center rounded-full text-xs font-semibold">
                                                {sessions.length - index}
                                            </div>
                                            <div>
                                                <p className="text-bt-dark text-sm font-medium">
                                                    {formatTime(session.duration_seconds ?? 0)}
                                                </p>
                                                <p className="text-bt-muted mt-0.5 text-xs">
                                                    {formatDate(session.started_at)}
                                                </p>
                                            </div>
                                        </div>
                                        <div className="text-right">
                                            <p className="text-bt-dark text-sm font-medium">
                                                {session.pages_read} pages
                                            </p>
                                            <p className="text-bt-muted-light mt-0.5 text-xs">
                                                {session.started_at
                                                    ? new Date(
                                                          session.started_at,
                                                      ).toLocaleTimeString('en-US', {
                                                          hour: '2-digit',
                                                          minute: '2-digit',
                                                      })
                                                    : '—'}
                                            </p>
                                        </div>
                                    </div>
                                ))}
                            </div>
                        )}
                    </div>
                </div>
            </div>

            {showDeleteModal && (
                <div className="bg-bt-dark/50 fixed inset-0 z-200 flex items-center justify-center">
                    <div className="bg-bt-surface mx-4 w-full max-w-sm overflow-hidden rounded-xl">
                        <div className="border-bt-border border-b px-6 py-5">
                            <h2 className="text-bt-dark font-serif text-lg font-semibold">
                                Remove book?
                            </h2>
                        </div>
                        <div className="px-6 py-5">
                            <p className="text-bt-warm text-sm leading-relaxed">
                                <span className="text-bt-dark font-medium">{book.title}</span> will
                                be removed from your library along with all reading sessions. This
                                cannot be undone.
                            </p>
                        </div>
                        {deleteError && (
                            <div className="px-6 pb-4">
                                <p className="text-bt-error bg-bt-error-bg border-bt-error-border rounded-md border px-3 py-2 text-xs font-medium">
                                    {deleteError}
                                </p>
                            </div>
                        )}
                        <div className="flex gap-3 px-6 pb-6">
                            <button
                                onClick={() => {
                                    setShowDeleteModal(false)
                                    setDeleteError('')
                                }}
                                className="bg-bt-cream text-bt-warm border-bt-input-border hover:border-bt-gold flex-1 cursor-pointer rounded-md border py-2.5 text-sm font-medium transition-colors duration-200"
                            >
                                Cancel
                            </button>
                            <button
                                onClick={handleDelete}
                                disabled={isDeleting}
                                className="bg-bt-danger hover:bg-bt-danger-dark flex-1 cursor-pointer rounded-md py-2.5 text-sm font-medium text-white transition-colors duration-200"
                            >
                                {isDeleting ? 'Removing...' : 'Yes, remove'}
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </Layout>
    )
}

export default BookDetail
