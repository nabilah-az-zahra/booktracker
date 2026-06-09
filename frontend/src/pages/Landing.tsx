import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { BookOpen, MoveRight, Feather } from 'lucide-react'
import { RevealSection } from '../components/RevealSection'
import { FaqItem } from '../components/FaqItem'
import { Ticker } from '../components/Ticker'
import { MOCK_LIBRARY_BOOKS, FAQ_ITEMS } from '../constants/landingData'

const Landing = () => {
    const [scrolled, setScrolled] = useState(false)

    useEffect(() => {
        const handleScroll = () => setScrolled(window.scrollY > 40)
        window.addEventListener('scroll', handleScroll, { passive: true })
        return () => window.removeEventListener('scroll', handleScroll)
    }, [])

    return (
        <div className="bg-bt-cream text-bt-dark min-h-screen">
            <nav
                className={`fixed inset-x-0 top-0 z-50 transition-all duration-300 ${
                    scrolled
                        ? 'bg-bt-surface border-bt-border border-b shadow-sm'
                        : 'bg-transparent'
                }`}
            >
                <div className="mx-auto flex h-16 max-w-6xl items-center justify-between px-6 md:px-8">
                    <Link to="/" className="group flex items-center gap-3">
                        <div className="bg-bt-dark group-hover:bg-bt-gold flex h-8 w-8 items-center justify-center rounded-sm transition-colors duration-200">
                            <BookOpen
                                size={15}
                                className="text-bt-gold group-hover:text-bt-dark transition-colors duration-200"
                                strokeWidth={1.5}
                            />
                        </div>
                        <span className="text-base font-semibold tracking-tight">
                            Book<span className="text-bt-gold">Tracker</span>
                        </span>
                    </Link>
                    <div className="flex items-center gap-6">
                        <Link
                            to="/login"
                            className="text-bt-muted hover:text-bt-dark py-1 text-sm transition-colors"
                        >
                            Sign in
                        </Link>
                        <Link
                            to="/register"
                            className="border-bt-dark hover:bg-bt-dark hover:text-bt-cream flex items-center gap-2 rounded-sm border px-4 py-2 text-sm transition-colors duration-200"
                        >
                            Get started
                            <Feather size={12} strokeWidth={1.5} />
                        </Link>
                    </div>
                </div>
            </nav>

            <header className="mx-auto max-w-6xl px-6 pt-36 pb-16 md:px-8">
                <div className="max-w-3xl">
                    <blockquote className="mb-6 text-4xl leading-tight font-medium tracking-tight md:text-6xl">
                        "I can never read all the books I want."
                    </blockquote>
                    <div className="mb-16 flex items-center gap-3">
                        <Feather size={11} strokeWidth={1.5} className="text-bt-gold" />
                        <cite className="text-bt-muted text-xs tracking-widest uppercase not-italic">
                            Sylvia Plath
                        </cite>
                    </div>
                    <p className="text-bt-warm mb-10 max-w-md text-base leading-relaxed">
                        So you might as well track the ones you do read. A timer, a library, and
                        some numbers. That's it.
                    </p>
                    <div className="flex items-center gap-5">
                        <Link
                            to="/register"
                            className="bg-bt-dark text-bt-cream hover:bg-bt-gold hover:text-bt-dark inline-flex items-center gap-2 rounded-lg px-6 py-3 text-sm transition-colors duration-200"
                        >
                            Start for free <MoveRight size={14} strokeWidth={1.5} />
                        </Link>
                        <Link
                            to="/login"
                            className="text-bt-muted hover:text-bt-dark decoration-bt-border hover:decoration-bt-gold text-sm underline underline-offset-4 transition-colors"
                        >
                            Sign in
                        </Link>
                    </div>
                </div>
            </header>

            <Ticker />

            <RevealSection className="bg-bt-surface border-bt-border border-y">
                <div className="mx-auto max-w-6xl px-6 py-24 md:px-8">
                    <div className="grid grid-cols-1 items-start gap-16 lg:grid-cols-2">
                        <div className="reveal-child">
                            <p className="text-bt-muted mb-6 text-xs tracking-widest uppercase">
                                Reading timer
                            </p>
                            <div className="bg-bt-cream border-bt-border mb-6 rounded-xl border p-8">
                                <div className="py-8 text-center">
                                    <p className="mb-3 text-6xl font-semibold tracking-tight tabular-nums">
                                        12:34
                                    </p>
                                    <div className="flex items-center justify-center gap-2">
                                        <div className="bg-bt-gold h-1.5 w-1.5 animate-pulse rounded-full" />
                                        <span className="text-bt-muted text-xs">Reading...</span>
                                    </div>
                                </div>
                                <div className="mt-4 flex justify-center gap-3">
                                    <div className="bg-bt-dark text-bt-cream rounded-lg px-4 py-2 text-xs">
                                        Pause
                                    </div>
                                    <div className="border-bt-border text-bt-muted rounded-lg border px-4 py-2 text-xs">
                                        Done
                                    </div>
                                </div>
                            </div>
                            <p className="text-bt-muted text-sm leading-relaxed">
                                Start when you open the book. Pause when your phone goes off. Stop
                                when you're done. Enter your page number. That's the whole flow.
                            </p>
                        </div>

                        <div className="space-y-10">
                            <div className="reveal-child">
                                <p className="text-bt-muted mb-6 text-xs tracking-widest uppercase">
                                    Your library
                                </p>
                                <div className="bg-bt-cream border-bt-border mb-6 overflow-hidden rounded-xl border">
                                    {MOCK_LIBRARY_BOOKS.map((book, i) => (
                                        <div
                                            key={book.title}
                                            className={`flex items-center gap-3 px-5 py-3 ${i < MOCK_LIBRARY_BOOKS.length - 1 ? 'border-bt-border border-b' : ''}`}
                                        >
                                            <div
                                                className={`h-4 w-1 shrink-0 rounded-full ${
                                                    book.status === 'reading'
                                                        ? 'bg-bt-gold'
                                                        : book.status === 'finished'
                                                          ? 'bg-bt-success'
                                                          : 'bg-bt-muted/30'
                                                }`}
                                            />
                                            <span className="text-bt-dark truncate text-sm">
                                                {book.title}
                                            </span>
                                        </div>
                                    ))}
                                </div>
                                <p className="text-bt-muted text-sm leading-relaxed">
                                    Search and add books from a database of millions. Mark them
                                    want-to-read, reading, or finished.
                                </p>
                            </div>

                            <div className="reveal-child">
                                <p className="text-bt-muted mb-6 text-xs tracking-widest uppercase">
                                    Stats
                                </p>
                                <div className="mb-6 grid grid-cols-2 gap-3">
                                    {[
                                        { value: '47h', label: 'Time read' },
                                        { value: '12', label: 'Books finished' },
                                        { value: '3,241', label: 'Pages' },
                                        { value: '18', label: 'Day streak' },
                                    ].map((stat) => (
                                        <div
                                            key={stat.label}
                                            className="bg-bt-cream border-bt-border rounded-lg border px-4 py-3"
                                        >
                                            <p className="mb-0.5 text-xl font-semibold">
                                                {stat.value}
                                            </p>
                                            <p className="text-bt-muted text-[10px] tracking-wider uppercase">
                                                {stat.label}
                                            </p>
                                        </div>
                                    ))}
                                </div>
                                <p className="text-bt-muted text-sm leading-relaxed">
                                    Hours, pages, streaks, yearly goals. Useful if you're the kind
                                    of person who likes to see the numbers go up.
                                </p>
                            </div>
                        </div>
                    </div>
                </div>
            </RevealSection>

            <RevealSection className="bg-bt-cream border-bt-border border-b">
                <div className="mx-auto max-w-6xl px-6 py-24 md:px-8">
                    <div className="grid grid-cols-1 gap-12 md:grid-cols-3 md:gap-16">
                        <div className="reveal-child">
                            <div className="bg-bt-gold mb-6 h-px w-6" />
                            <h2 className="mb-3 text-2xl font-semibold">Questions</h2>
                            <p className="text-bt-muted mb-8 text-sm leading-relaxed">
                                If yours isn't here, make an account and poke around. It's free.
                            </p>
                            <Link
                                to="/register"
                                className="text-bt-gold hover:text-bt-dark inline-flex items-center gap-1.5 text-xs font-medium transition-colors"
                            >
                                Get started <MoveRight size={11} strokeWidth={1.5} />
                            </Link>
                        </div>
                        <div className="reveal-child bg-bt-surface border-bt-border divide-bt-border divide-y rounded-xl border px-6 py-2 md:col-span-2">
                            {FAQ_ITEMS.map((item) => (
                                <FaqItem key={item.q} q={item.q} a={item.a} />
                            ))}
                        </div>
                    </div>
                </div>
            </RevealSection>

            <RevealSection className="bg-bt-dark text-bt-cream">
                <div className="mx-auto max-w-6xl px-6 py-20 md:px-8">
                    <div className="reveal-child max-w-sm">
                        <p className="mb-3 text-2xl leading-snug font-semibold">
                            Free to use.
                            <br />
                            No nonsense.
                        </p>
                        <p className="text-bt-warm mb-8 text-sm">
                            Make an account. Add a book. Start a session.
                        </p>
                        <Link
                            to="/register"
                            className="bg-bt-gold text-bt-dark hover:bg-bt-cream inline-flex items-center gap-2 rounded-lg px-6 py-3 text-sm font-medium transition-colors duration-200"
                        >
                            Get started <MoveRight size={14} strokeWidth={1.5} />
                        </Link>
                    </div>
                </div>
            </RevealSection>

            <footer className="bg-bt-dark text-bt-warm border-bt-warm/10 border-t">
                <div className="mx-auto max-w-6xl px-6 py-16 md:px-8">
                    <div className="border-bt-warm/10 flex flex-col items-start justify-between gap-10 border-b pb-12 md:flex-row">
                        <div className="max-w-xs">
                            <Link to="/" className="group mb-4 flex items-center gap-3">
                                <div className="border-bt-gold group-hover:bg-bt-gold flex h-8 w-8 items-center justify-center rounded-sm border transition-colors duration-200">
                                    <BookOpen
                                        size={14}
                                        className="text-bt-gold group-hover:text-bt-dark transition-colors duration-200"
                                        strokeWidth={1.5}
                                    />
                                </div>
                                <span className="text-bt-cream text-sm font-semibold">
                                    Book<span className="text-bt-gold">Tracker</span>
                                </span>
                            </Link>
                            <p className="text-bt-warm/70 text-xs leading-relaxed">
                                Built because spreadsheets got annoying.
                            </p>
                        </div>
                        <div className="flex gap-16">
                            <div>
                                <p className="text-bt-cream mb-4 text-[10px] font-semibold tracking-widest uppercase">
                                    Account
                                </p>
                                <div className="space-y-2">
                                    <Link
                                        to="/login"
                                        className="hover:text-bt-cream block text-xs transition-colors"
                                    >
                                        Sign in
                                    </Link>
                                    <Link
                                        to="/register"
                                        className="hover:text-bt-cream block text-xs transition-colors"
                                    >
                                        Create account
                                    </Link>
                                </div>
                            </div>
                            <div>
                                <p className="text-bt-cream mb-4 text-[10px] font-semibold tracking-widest uppercase">
                                    Explore
                                </p>
                                <div className="space-y-2">
                                    <Link
                                        to="/books/search"
                                        className="hover:text-bt-cream block text-xs transition-colors"
                                    >
                                        Find books
                                    </Link>
                                    <Link
                                        to="/register"
                                        className="hover:text-bt-cream block text-xs transition-colors"
                                    >
                                        Start tracking
                                    </Link>
                                </div>
                            </div>
                        </div>
                    </div>
                    <div className="text-bt-warm/40 flex flex-col items-start justify-between gap-3 pt-6 text-[10px] tracking-wider uppercase md:flex-row md:items-center">
                        <div className="flex items-center gap-2">
                            <Feather size={11} strokeWidth={1.5} />
                            <span>Built for readers</span>
                        </div>
                        <p>© {new Date().getFullYear()} BookTracker</p>
                    </div>
                </div>
            </footer>
        </div>
    )
}

export default Landing
