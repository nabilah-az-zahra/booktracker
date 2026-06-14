import { useState, useEffect } from 'react'
import { useParams, useNavigate, Link } from 'react-router-dom'
import Layout from '../components/Layout'
import api from '../api/axios'
import { useTimer } from '../hooks/useTimer'
import type { Book, ReadingSession, ReadingProgress } from '../types'
import { Play, Pause, Square, BookOpen, ArrowLeft, Check, BookMarked } from 'lucide-react'
import { getApiError } from '../utils/errors'
import { formatTime } from '../utils/formatUtils'

const ReadingTimer = () => {
    const { id } = useParams<{ id: string }>()
    const navigate = useNavigate()
    const { seconds, isRunning, isPaused, start, pause, resume, reset, restore } = useTimer()

    const [book, setBook] = useState<Book | null>(null)
    const [session, setSession] = useState<ReadingSession | null>(null)
    const [progress, setProgress] = useState<ReadingProgress | null>(null)
    const [loading, setLoading] = useState(true)
    const [showStopModal, setShowStopModal] = useState(false)
    const [currentPageInput, setCurrentPageInput] = useState('')
    const [submitting, setSubmitting] = useState(false)
    const [done, setDone] = useState(false)
    const [finalSeconds, setFinalSeconds] = useState(0)
    const [error, setError] = useState('')
    const [activeBookId, setActiveBookId] = useState<string | null>(null)
    const [showCancelModal, setShowCancelModal] = useState(false)
    const [cancelling, setCancelling] = useState(false)

    useEffect(() => {
        const fetchData = async () => {
            try {
                const [bookRes, activeRes, progressRes] = await Promise.all([
                    api.get(`/api/books/${id}`),
                    api.get('/api/sessions/active'),
                    api.get(`/api/progress/${id}`).catch(() => ({ data: { data: null } })),
                ])
                setBook(bookRes.data.data)
                setProgress(progressRes.data.data)

                const active = activeRes.data.data
                if (active) {
                    if (active.book_id === id) {
                        setSession(active)
                        if (active.status === 'active') {
                            start()
                        } else if (active.status === 'paused') {
                            restore(active.duration_seconds ?? 0)
                        }
                    } else {
                        setActiveBookId(active.book_id)
                    }
                }
            } catch (err) {
                console.error(err)
            } finally {
                setLoading(false)
            }
        }
        fetchData()
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [id])

    useEffect(() => {
        const handleUnload = (e: BeforeUnloadEvent) => {
            if (session && isRunning) {
                e.preventDefault()
            }
        }
        window.addEventListener('beforeunload', handleUnload)
        return () => window.removeEventListener('beforeunload', handleUnload)
    }, [session, isRunning])

    useEffect(() => {
        return () => {
            reset()
        }
    }, [reset])

    const handleStart = async () => {
        try {
            const res = await api.post('/api/sessions/start', { book_id: id })
            setSession(res.data.data)
            start()
        } catch (err) {
            setError(getApiError(err))
        }
    }

    const handlePause = async () => {
        if (!session) return
        pause()
        try {
            await api.patch(`/api/sessions/${session.id}/pause`, {
                duration_seconds: seconds,
            })
        } catch (err) {
            console.error('Failed to sync pause state with server:', err)
        }
    }

    const handleResume = async () => {
        if (!session) return
        resume()
        try {
            await api.patch(`/api/sessions/${session.id}/resume`)
        } catch (err) {
            console.error('Failed to sync resume state with server:', err)
        }
    }

    const handleStop = async () => {
        if (!session || submitting) return
        const currentPage = parseInt(currentPageInput)

        if (isNaN(currentPage) || currentPage < 0) {
            setError('Please enter a valid page number')
            return
        }

        if (currentPage < lastPage) {
            setError(`Page must be at least ${lastPage} (your last saved page)`)
            return
        }

        if (book?.total_pages && currentPage > book.total_pages) {
            setError(`Page number cannot exceed total pages (${book.total_pages})`)
            return
        }

        const previousPage = progress?.current_page ?? 0
        const pagesReadThisSession = Math.max(0, currentPage - previousPage)

        setSubmitting(true)
        try {
            await api.patch(`/api/sessions/${session.id}/stop`, {
                duration_seconds: seconds,
                pages_read: pagesReadThisSession,
            })
            setFinalSeconds(seconds)
            setDone(true)
            setSession(null)
            reset()
        } catch (err) {
            console.error('Stop session error:', err)
            setError('Failed to save session. Please try again.')
        } finally {
            setSubmitting(false)
        }
    }

    const handleCancel = async () => {
        if (!session || cancelling) return
        setCancelling(true)
        try {
            pause()
            await api.delete(`/api/sessions/${session.id}`)
            setSession(null)
            reset()
            setShowCancelModal(false)
            navigate(`/books/${id}`)
        } catch (err) {
            setError(getApiError(err))
        } finally {
            setCancelling(false)
        }
    }

    const getTimerClassName = () => {
        if (!isRunning && !isPaused) return 'text-bt-placeholder'
        if (isPaused) return 'text-bt-gold'
        return 'text-bt-dark'
    }

    const lastPage = progress?.current_page ?? 0

    if (loading)
        return (
            <Layout>
                <div className="flex items-center justify-center py-20">
                    <p className="text-bt-muted-light text-sm">Loading...</p>
                </div>
            </Layout>
        )

    if (!book)
        return (
            <Layout>
                <div className="py-20 text-center">
                    <p className="text-bt-muted text-sm">Book not found</p>
                </div>
            </Layout>
        )

    if (done)
        return (
            <Layout>
                <div className="mx-auto max-w-md py-16 text-center">
                    <div className="bg-bt-status-finished-bg mx-auto mb-6 flex h-16 w-16 items-center justify-center rounded-full">
                        <Check
                            size={28}
                            className="text-bt-status-finished-text"
                            strokeWidth={1.5}
                        />
                    </div>
                    <h2 className="text-bt-dark mb-2 font-serif text-2xl font-semibold">
                        Session saved!
                    </h2>
                    <p className="text-bt-muted mb-1 text-sm">Great reading session. Keep it up!</p>
                    <p className="text-bt-gold my-8 font-serif text-4xl font-semibold">
                        {formatTime(finalSeconds)}
                    </p>
                    <div className="flex justify-center gap-3">
                        <button
                            onClick={() => {
                                setDone(false)
                                setShowStopModal(false)
                                setCurrentPageInput('')
                                setError('')
                            }}
                            className="bg-bt-dark text-bt-cream hover:bg-bt-gold cursor-pointer rounded-lg px-5 py-2.5 text-sm font-medium transition-all duration-200"
                        >
                            Read more
                        </button>
                        <button
                            onClick={() => navigate(`/books/${id}`)}
                            className="bg-bt-cream text-bt-warm border-bt-input-border hover:border-bt-gold cursor-pointer rounded-lg border px-5 py-2.5 text-sm font-medium transition-all duration-200"
                        >
                            View book
                        </button>
                    </div>
                </div>
            </Layout>
        )

    return (
        <Layout>
            <button
                onClick={() => navigate(`/books/${id}`)}
                className="text-bt-muted hover:text-bt-dark mb-8 flex cursor-pointer items-center gap-1.5 text-sm transition-colors"
            >
                <ArrowLeft size={14} /> Back to book
            </button>

            <div className="mx-auto max-w-md">
                <div className="mb-6 flex items-center gap-4">
                    {book.cover_url ? (
                        <img
                            src={book.cover_url}
                            alt={book.title}
                            className="h-16 w-12 shrink-0 rounded object-cover shadow-md"
                        />
                    ) : (
                        <div className="bg-bt-accent-bg flex h-16 w-12 shrink-0 items-center justify-center rounded">
                            <BookOpen size={16} className="text-bt-gold" strokeWidth={1.5} />
                        </div>
                    )}
                    <div className="min-w-0">
                        <p className="text-bt-dark font-serif text-base leading-snug font-semibold">
                            {book.title}
                        </p>
                        <p className="text-bt-muted mt-0.5 text-xs">{book.author}</p>
                    </div>
                </div>

                {lastPage > 0 && (
                    <div className="bg-bt-accent-bg mb-6 flex items-center gap-3 rounded-lg px-4 py-3">
                        <BookMarked size={14} className="text-bt-gold" strokeWidth={1.5} />
                        <p className="text-bt-muted-dark text-xs">
                            Last session ended on{' '}
                            <span className="text-bt-dark font-semibold">page {lastPage}</span>
                            {book.total_pages > 0 && (
                                <span className="text-bt-muted"> of {book.total_pages}</span>
                            )}
                        </p>
                    </div>
                )}

                <div className="bg-bt-surface border-bt-border mb-8 rounded-2xl border py-14 text-center">
                    <p
                        className={`mb-3 font-serif text-7xl font-semibold tracking-tight transition-colors duration-300 ${getTimerClassName()}`}
                        style={{ fontVariantNumeric: 'tabular-nums' }}
                    >
                        {formatTime(seconds)}
                    </p>
                    <p className="text-bt-muted-light text-sm">
                        {!session && !isRunning
                            ? 'Ready to start'
                            : isPaused
                              ? 'Paused'
                              : isRunning
                                ? 'Reading...'
                                : 'Ready'}
                    </p>
                </div>

                {error && !showStopModal && (
                    <div className="border-bt-error-border bg-bt-error-bg text-bt-error mb-4 rounded-md border-l-[3px] px-4 py-3 text-sm">
                        {error}
                    </div>
                )}

                <div className="flex items-center justify-center gap-4">
                    {!session && !isRunning && (
                        <button
                            onClick={handleStart}
                            className="bg-bt-dark text-bt-cream hover:bg-bt-gold flex cursor-pointer items-center gap-2.5 rounded-[10px] px-8 py-3.5 text-sm font-medium transition-all duration-200"
                        >
                            <Play size={16} strokeWidth={2} /> Start Reading
                        </button>
                    )}

                    {isRunning && (
                        <>
                            <button
                                onClick={handlePause}
                                className="bg-bt-surface text-bt-dark border-bt-input-border hover:border-bt-gold flex cursor-pointer items-center gap-2 rounded-[10px] border px-6 py-3 text-sm font-medium transition-all duration-200"
                            >
                                <Pause size={15} strokeWidth={2} /> Pause
                            </button>
                            <button
                                onClick={() => {
                                    setShowStopModal(true)
                                    setError('')
                                }}
                                className="bg-bt-dark text-bt-cream hover:bg-bt-gold flex cursor-pointer items-center gap-2 rounded-[10px] px-6 py-3 text-sm font-medium transition-all duration-200"
                            >
                                <Square size={15} strokeWidth={2} /> Done
                            </button>
                        </>
                    )}

                    {isPaused && (
                        <>
                            <button
                                onClick={handleResume}
                                className="bg-bt-dark text-bt-cream hover:bg-bt-gold flex cursor-pointer items-center gap-2 rounded-[10px] px-6 py-3 text-sm font-medium transition-all duration-200"
                            >
                                <Play size={15} strokeWidth={2} /> Resume
                            </button>
                            <button
                                onClick={() => {
                                    setShowStopModal(true)
                                    setError('')
                                }}
                                className="bg-bt-surface text-bt-warm border-bt-input-border hover:border-bt-gold flex cursor-pointer items-center gap-2 rounded-[10px] border px-6 py-3 text-sm font-medium transition-all duration-200"
                            >
                                <Square size={15} strokeWidth={2} /> Done
                            </button>
                        </>
                    )}

                    {(isRunning || isPaused) && (
                        <button
                            onClick={() => setShowCancelModal(true)}
                            className="text-bt-muted hover:text-bt-danger cursor-pointer text-xs transition-colors"
                        >
                            Cancel session
                        </button>
                    )}
                </div>

                {activeBookId && !session && (
                    <div className="border-bt-banner-border bg-bt-banner-bg mb-6 flex items-center justify-between rounded-lg border px-4 py-3">
                        <p className="text-bt-dark text-sm">
                            You have an active session on another book
                        </p>
                        <Link
                            to={`/reading/${activeBookId}`}
                            className="text-bt-gold hover:text-bt-muted-dark text-xs font-medium transition-colors"
                        >
                            Go there →
                        </Link>
                    </div>
                )}
            </div>

            {showStopModal && (
                <div className="bg-bt-dark/50 fixed inset-0 z-200 flex items-center justify-center">
                    <div className="bg-bt-surface mx-4 w-full max-w-sm overflow-hidden rounded-2xl">
                        <div className="px-6 pt-8 pb-6 text-center">
                            <p className="text-bt-muted-light mb-2 text-xs font-medium tracking-widest uppercase">
                                Session time
                            </p>
                            <p className="text-bt-dark mb-6 font-serif text-5xl font-semibold">
                                {formatTime(seconds)}
                            </p>

                            <p className="text-bt-dark mb-1 font-serif text-lg font-semibold">
                                What page are you on now?
                            </p>

                            {lastPage > 0 && (
                                <p className="text-bt-muted mb-4 text-xs">
                                    You were on page{' '}
                                    <button
                                        onClick={() => setCurrentPageInput(String(lastPage))}
                                        className="text-bt-gold hover:text-bt-muted-dark cursor-pointer font-semibold transition-colors"
                                    >
                                        {lastPage}
                                    </button>{' '}
                                    last time
                                    {book?.total_pages ? ` (of ${book.total_pages})` : ''}
                                </p>
                            )}

                            {!lastPage && (
                                <p className="text-bt-muted mb-4 text-xs">
                                    Enter the page number you stopped at
                                    {book?.total_pages ? ` (of ${book.total_pages})` : ''}
                                </p>
                            )}

                            <input
                                type="number"
                                value={currentPageInput}
                                onChange={(e) => setCurrentPageInput(e.target.value)}
                                placeholder={lastPage > 0 ? `Last: ${lastPage}` : 'e.g. 42'}
                                min={lastPage}
                                max={book?.total_pages || undefined}
                                autoFocus
                                className="input-field mb-2 w-full rounded-[10px] py-3 text-center font-serif text-2xl transition-all duration-200"
                            />

                            {currentPageInput && parseInt(currentPageInput) > lastPage && (
                                <p className="text-bt-status-finished-text mb-4 text-xs">
                                    +{parseInt(currentPageInput) - lastPage} pages read this session
                                </p>
                            )}

                            {error && <p className="text-bt-danger mb-4 text-xs">{error}</p>}

                            <div className="mt-4 flex gap-3">
                                <button
                                    onClick={handleStop}
                                    disabled={submitting || currentPageInput === ''}
                                    className="bg-bt-dark text-bt-cream hover:bg-bt-gold flex-1 cursor-pointer rounded-lg py-3 text-sm font-medium transition-all duration-200 disabled:opacity-50"
                                >
                                    {submitting ? 'Saving...' : 'Save session'}
                                </button>
                                <button
                                    onClick={() => {
                                        setShowStopModal(false)
                                        setError('')
                                    }}
                                    className="bg-bt-cream text-bt-warm border-bt-input-border hover:border-bt-gold flex-1 cursor-pointer rounded-lg border py-3 text-sm font-medium transition-all duration-200"
                                >
                                    Keep reading
                                </button>
                            </div>
                        </div>
                    </div>
                </div>
            )}
            {showCancelModal && (
                <div className="bg-bt-dark/50 fixed inset-0 z-200 flex items-center justify-center">
                    <div className="bg-bt-surface mx-4 w-full max-w-sm overflow-hidden rounded-xl">
                        <div className="border-bt-border border-b px-6 py-5">
                            <h2 className="text-bt-dark font-serif text-lg font-semibold">
                                Cancel session?
                            </h2>
                        </div>
                        <div className="px-6 py-5">
                            <p className="text-bt-warm text-sm leading-relaxed">
                                This session will be discarded and your reading time won't be saved.
                            </p>
                        </div>
                        <div className="flex gap-3 px-6 pb-6">
                            <button
                                onClick={handleCancel}
                                disabled={cancelling}
                                className="bg-bt-danger hover:bg-bt-danger-dark flex-1 cursor-pointer rounded-md py-2.5 text-sm font-medium text-white transition-all duration-200 disabled:opacity-50"
                            >
                                {cancelling ? 'Cancelling...' : 'Yes, cancel'}
                            </button>
                            <button
                                onClick={() => setShowCancelModal(false)}
                                className="bg-bt-cream text-bt-warm border-bt-input-border hover:border-bt-gold flex-1 cursor-pointer rounded-md border py-2.5 text-sm font-medium transition-all duration-200"
                            >
                                Keep reading
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </Layout>
    )
}

export default ReadingTimer
