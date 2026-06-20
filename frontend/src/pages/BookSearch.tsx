import { useState, useEffect, useCallback, useRef } from 'react'
import { useNavigate } from 'react-router-dom'
import Layout from '../components/Layout'
import api from '../api/axios'
import type { GoogleBook, Book } from '../types'
import {
    Search,
    BookOpen,
    Plus,
    Check,
    Loader2,
    ChevronLeft,
    ChevronRight,
    Library,
} from 'lucide-react'

const RESULTS_PER_PAGE = 9

const BookSearch = () => {
    const navigate = useNavigate()
    const [query, setQuery] = useState('')
    const [submittedQuery, setSubmittedQuery] = useState('')
    const [results, setResults] = useState<GoogleBook[]>([])
    const [totalItems, setTotalItems] = useState(0)
    const [currentPage, setCurrentPage] = useState(1)
    const [searching, setSearching] = useState(false)
    const [addingId, setAddingId] = useState<string | null>(null)
    const [addedIds, setAddedIds] = useState<Set<string>>(new Set())
    const [existingTitles, setExistingTitles] = useState<Set<string>>(new Set())
    const [error, setError] = useState('')
    const lastSearchTimeRef = useRef(0)

    useEffect(() => {
        const fetchLibrary = async () => {
            try {
                const res = await api.get('/api/books')
                const books: Book[] = res.data.data || []
                setExistingTitles(new Set(books.map((b) => b.title.toLowerCase().trim())))
            } catch (err) {
                console.error(err)
            }
        }
        fetchLibrary()
    }, [])

    const searchBooks = useCallback(async (searchQuery: string, page: number) => {
        if (!searchQuery.trim()) return
        setSearching(true)
        setError('')
        try {
            const startIndex = (page - 1) * RESULTS_PER_PAGE
            const res = await fetch(
                `https://www.googleapis.com/books/v1/volumes?q=${encodeURIComponent(searchQuery)}&maxResults=${RESULTS_PER_PAGE}&startIndex=${startIndex}&key=${import.meta.env.VITE_GOOGLE_BOOKS_API_KEY}`,
            )
            if (!res.ok) {
                throw new Error(`Google Books API error: ${res.status}`)
            }
            const data = await res.json()
            setResults(data.items || [])
            setTotalItems(Math.min(data.totalItems || 0, 100))
        } catch {
            setError('Failed to search books. Please try again.')
        } finally {
            setSearching(false)
        }
    }, [])

    const handleSubmit = useCallback(
        async (e: React.FormEvent) => {
            e.preventDefault()
            if (!query.trim()) return
            const now = Date.now()
            if (now - lastSearchTimeRef.current < 1000) return
            lastSearchTimeRef.current = now
            setCurrentPage(1)
            setSubmittedQuery(query)
            setResults([])
            await searchBooks(query, 1)
        },
        [query, searchBooks],
    )

    const handlePageChange = useCallback(
        async (page: number) => {
            setCurrentPage(page)
            await searchBooks(submittedQuery, page)
            window.scrollTo({ top: 0, behavior: 'smooth' })
        },
        [submittedQuery, searchBooks],
    )

    const isAlreadyInLibrary = (book: GoogleBook) =>
        existingTitles.has(book.volumeInfo.title.toLowerCase().trim())

    const addBook = async (book: GoogleBook) => {
        if (isAlreadyInLibrary(book)) return
        setAddingId(book.id)
        try {
            await api.post('/api/books', {
                title: book.volumeInfo.title,
                author: book.volumeInfo.authors?.join(', ') || 'Unknown',
                cover_url:
                    book.volumeInfo.imageLinks?.thumbnail?.replace('http://', 'https://') || '',
                total_pages: book.volumeInfo.pageCount || 0,
                status: 'want_to_read',
            })
            setAddedIds((prev) => new Set(prev).add(book.id))
            setExistingTitles((prev) =>
                new Set(prev).add(book.volumeInfo.title.toLowerCase().trim()),
            )
        } catch {
            setError('Failed to add book. Please try again.')
        } finally {
            setAddingId(null)
        }
    }

    const totalPages = Math.ceil(totalItems / RESULTS_PER_PAGE)

    const getPageNumbers = () => {
        const pages: (number | string)[] = []

        if (totalPages <= 7) {
            for (let i = 1; i <= totalPages; i++) pages.push(i)
        } else {
            pages.push(1)
            if (currentPage > 4) {
                pages.push('...')
            }
            const start = Math.max(2, currentPage - 1)
            const end = Math.min(totalPages - 1, currentPage + 1)
            for (let i = start; i <= end; i++) {
                pages.push(i)
            }
            if (currentPage < totalPages - 3) {
                pages.push('...')
            }
            pages.push(totalPages)
        }
        return pages
    }

    return (
        <Layout>
            <div className="mb-8">
                <h1 className="text-bt-dark mb-1 font-serif text-2xl font-semibold">Find Books</h1>
                <p className="text-bt-muted text-sm">
                    Search millions of books and add them to your library
                </p>
            </div>

            <form onSubmit={handleSubmit} className="mb-8">
                <div className="flex max-w-xl gap-3">
                    <div className="relative flex-1">
                        <Search
                            size={15}
                            className="text-bt-muted-light absolute top-1/2 left-3.5 -translate-y-1/2"
                        />
                        <input
                            type="text"
                            value={query}
                            onChange={(e) => setQuery(e.target.value)}
                            placeholder="Search by title, author, or ISBN..."
                            className="input-field w-full py-2.5 pr-4 pl-10 text-sm transition-all duration-200"
                        />
                    </div>
                    <button
                        type="submit"
                        disabled={searching || !query.trim()}
                        className="bg-bt-dark text-bt-cream hover:bg-bt-gold flex cursor-pointer items-center gap-2 rounded-lg px-5 py-2.5 text-sm font-medium transition-all duration-200 disabled:opacity-50"
                    >
                        {searching ? (
                            <Loader2 size={14} className="animate-spin" />
                        ) : (
                            <Search size={14} />
                        )}
                        Search
                    </button>
                </div>

                {totalItems > 0 && !searching && (
                    <p className="text-bt-muted-light mt-2 text-xs">
                        About {totalItems} results for "{submittedQuery}" — page {currentPage} of{' '}
                        {totalPages}
                    </p>
                )}
            </form>

            {error && (
                <div className="border-bt-error-border bg-bt-error-bg text-bt-error mb-6 max-w-xl rounded-md border-l-[3px] px-4 py-3 text-sm">
                    {error}
                </div>
            )}

            {searching ? (
                <div className="flex items-center justify-center py-20">
                    <div className="flex flex-col items-center gap-3">
                        <Loader2 size={24} className="text-bt-gold animate-spin" />
                        <p className="text-bt-muted text-sm">Searching books...</p>
                    </div>
                </div>
            ) : results.length > 0 ? (
                <>
                    <div className="mb-8 grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
                        {results.map((book) => {
                            const isAdded = addedIds.has(book.id) || isAlreadyInLibrary(book)
                            const isAdding = addingId === book.id
                            return (
                                <div
                                    key={book.id}
                                    className={`bg-bt-surface flex gap-4 rounded-xl border p-4 transition-all duration-200 ${isAdded ? 'border-bt-gold' : 'border-bt-border'}`}
                                >
                                    {book.volumeInfo.imageLinks?.thumbnail ? (
                                        <img
                                            src={book.volumeInfo.imageLinks.thumbnail.replace(
                                                'http://',
                                                'https://',
                                            )}
                                            alt={book.volumeInfo.title}
                                            className="h-20 w-14 shrink-0 rounded object-cover"
                                        />
                                    ) : (
                                        <div className="bg-bt-accent-bg flex h-20 w-14 shrink-0 items-center justify-center rounded">
                                            <BookOpen
                                                size={18}
                                                className="text-bt-gold"
                                                strokeWidth={1.5}
                                            />
                                        </div>
                                    )}
                                    <div className="flex min-w-0 flex-1 flex-col justify-between">
                                        <div>
                                            <p className="text-bt-dark mb-1 line-clamp-2 text-sm leading-snug font-medium">
                                                {book.volumeInfo.title}
                                            </p>
                                            <p className="text-bt-muted mb-1 text-xs">
                                                {book.volumeInfo.authors?.join(', ') ||
                                                    'Unknown author'}
                                            </p>
                                            {book.volumeInfo.pageCount && (
                                                <p className="text-bt-muted-light text-xs">
                                                    {book.volumeInfo.pageCount} pages
                                                </p>
                                            )}
                                        </div>
                                        <button
                                            onClick={() => !isAdded && addBook(book)}
                                            disabled={isAdded || isAdding}
                                            className={`mt-3 flex cursor-pointer items-center gap-1.5 text-xs font-medium transition-all duration-200 disabled:cursor-default ${isAdded ? 'text-bt-gold' : 'text-bt-warm hover:text-bt-gold'}`}
                                        >
                                            {isAdding ? (
                                                <Loader2 size={12} className="animate-spin" />
                                            ) : isAdded ? (
                                                <Check size={12} />
                                            ) : (
                                                <Plus size={12} />
                                            )}
                                            {isAdding
                                                ? 'Adding...'
                                                : isAdded
                                                  ? 'Already in library'
                                                  : 'Add to library'}
                                        </button>
                                    </div>
                                </div>
                            )
                        })}
                    </div>

                    {totalPages > 1 && (
                        <div className="flex items-center justify-center gap-2">
                            <button
                                onClick={() => handlePageChange(currentPage - 1)}
                                disabled={currentPage === 1 || searching}
                                className="border-bt-border bg-bt-surface text-bt-warm hover:border-bt-gold flex h-8 w-8 cursor-pointer items-center justify-center rounded-md border transition-all duration-200 disabled:cursor-not-allowed disabled:opacity-30"
                            >
                                <ChevronLeft size={14} />
                            </button>

                            {getPageNumbers().map((page, index) => {
                                if (page === '...') {
                                    return (
                                        <span
                                            key={`gap-${index}`}
                                            className="text-bt-muted-light flex h-8 w-8 items-center justify-center text-xs select-none"
                                        >
                                            ...
                                        </span>
                                    )
                                }
                                return (
                                    <button
                                        key={page}
                                        onClick={() => handlePageChange(page as number)}
                                        disabled={searching}
                                        className={`flex h-8 w-8 cursor-pointer items-center justify-center rounded-md border text-xs font-medium transition-all duration-200 disabled:opacity-50 ${
                                            currentPage === page
                                                ? 'bg-bt-dark text-bt-cream border-bt-dark'
                                                : 'bg-bt-surface text-bt-warm border-bt-border hover:border-bt-gold'
                                        }`}
                                    >
                                        {page}
                                    </button>
                                )
                            })}

                            <button
                                onClick={() => handlePageChange(currentPage + 1)}
                                disabled={currentPage === totalPages || searching}
                                className="border-bt-border bg-bt-surface text-bt-warm hover:border-bt-gold flex h-8 w-8 cursor-pointer items-center justify-center rounded-md border transition-all duration-200 disabled:cursor-not-allowed disabled:opacity-30"
                            >
                                <ChevronRight size={14} />
                            </button>
                        </div>
                    )}
                </>
            ) : submittedQuery && !searching ? (
                <div className="border-bt-dashed bg-bt-surface rounded-xl border border-dashed py-16 text-center">
                    <BookOpen
                        size={28}
                        strokeWidth={1.5}
                        className="text-bt-placeholder mx-auto mb-3"
                    />
                    <p className="text-bt-muted text-sm">No books found for "{submittedQuery}"</p>
                    <p className="text-bt-muted-light mt-1 text-xs">Try a different search term</p>
                </div>
            ) : (
                <div className="border-bt-dashed bg-bt-surface rounded-xl border border-dashed py-16 text-center">
                    <Search
                        size={28}
                        strokeWidth={1.5}
                        className="text-bt-placeholder mx-auto mb-3"
                    />
                    <p className="text-bt-muted text-sm">Search for a book to get started</p>
                    <p className="text-bt-muted-light mt-1 text-xs">Powered by Google Books</p>
                </div>
            )}

            {addedIds.size > 0 && (
                <div className="mt-6 text-center">
                    <button
                        onClick={() => navigate('/library')}
                        className="group bg-bt-surface border-bt-border hover:border-bt-gold flex cursor-pointer items-center gap-4 rounded-xl border py-2 pr-4 pl-3.5 shadow-xs transition-all duration-300"
                    >
                        <div className="border-bt-border flex items-center gap-2 border-r pr-3">
                            <Library size={16} className="text-bt-warm" strokeWidth={1.5} />
                            <span className="text-bt-dark font-serif text-sm tracking-tight">
                                <strong className="text-bt-gold font-sans font-bold">
                                    {addedIds.size}
                                </strong>{' '}
                                book{addedIds.size > 1 ? 's' : ''} added
                            </span>
                        </div>
                        <span className="text-bt-warm group-hover:text-bt-dark flex items-center gap-1 font-sans text-xs font-bold tracking-wider uppercase transition-colors duration-200">
                            <span>View Library</span>
                            <span className="text-bt-gold transition-transform duration-200 group-hover:translate-x-0.5">
                                <ChevronRight size={14} />
                            </span>
                        </span>
                    </button>
                </div>
            )}
        </Layout>
    )
}

export default BookSearch
