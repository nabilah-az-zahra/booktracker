import { useState } from 'react'
import { NavLink, useNavigate, Link } from 'react-router-dom'
import { useAuth } from '../context/useAuth'
import { BookOpen, BookMarked, Search, BarChart2, LogOut, Menu, X, ChevronDown } from 'lucide-react'

const navItems = [
    { path: '/', label: 'Home', icon: BookOpen },
    { path: '/library', label: 'Library', icon: BookMarked },
    { path: '/books/search', label: 'Find Books', icon: Search },
    { path: '/stats', label: 'Stats', icon: BarChart2 },
]

const Layout = ({ children }: { children: React.ReactNode }) => {
    const { user, logout } = useAuth()
    const navigate = useNavigate()
    const [mobileOpen, setMobileOpen] = useState(false)
    const [showUserMenu, setShowUserMenu] = useState(false)
    const [showLogoutModal, setShowLogoutModal] = useState(false)

    const handleLogout = () => {
        logout()
        navigate('/welcome')
    }

    return (
        <div className="bg-bt-cream min-h-screen font-serif">
            <header className="bg-bt-surface border-bt-border sticky top-0 z-40 border-b">
                <div className="mx-auto max-w-5xl px-6">
                    <div className="flex h-14 items-center justify-between">
                        <Link to="/" className="flex shrink-0 items-center gap-2">
                            <div className="bg-bt-gold flex h-7 w-7 items-center justify-center rounded-md">
                                <BookOpen size={14} className="text-white" strokeWidth={1.5} />
                            </div>
                            <span className="text-bt-dark font-serif text-base font-semibold">
                                BookTracker
                            </span>
                        </Link>
                        <nav className="hidden items-center gap-1 md:flex">
                            {navItems.map((item) => (
                                <NavLink
                                    key={item.path}
                                    to={item.path}
                                    end={item.path === '/'}
                                    className={({ isActive }) =>
                                        `relative px-4 py-1.5 text-sm transition-all duration-150 ${isActive ? 'text-bt-dark font-medium' : 'text-bt-muted hover:text-bt-dark'}`
                                    }
                                >
                                    {({ isActive }) => (
                                        <>
                                            {item.label}
                                            {isActive && (
                                                <span className="bg-bt-gold absolute right-4 bottom-0 left-4 h-0.5 rounded-full" />
                                            )}
                                        </>
                                    )}
                                </NavLink>
                            ))}
                        </nav>
                        <div className="flex items-center gap-3">
                            <div className="relative hidden md:block">
                                <button
                                    onClick={() => setShowUserMenu(!showUserMenu)}
                                    className="hover:bg-bt-cream flex cursor-pointer items-center gap-2 rounded-lg px-2 py-1.5 transition-all duration-150"
                                >
                                    <div className="bg-bt-accent-bg text-bt-gold flex h-6 w-6 items-center justify-center rounded-full text-xs font-semibold">
                                        {user?.name?.charAt(0).toUpperCase()}
                                    </div>
                                    <span className="text-bt-dark text-xs font-medium">
                                        {user?.name?.split(' ')[0]}
                                    </span>
                                    <ChevronDown
                                        size={12}
                                        className={`text-bt-muted transition-transform ${showUserMenu ? 'rotate-180' : ''}`}
                                    />
                                </button>
                                {showUserMenu && (
                                    <>
                                        <div
                                            className="fixed inset-0 z-30"
                                            onClick={() => setShowUserMenu(false)}
                                        />
                                        <div className="bg-bt-surface border-bt-border absolute right-0 z-40 mt-1 w-48 rounded-xl border py-1 shadow-sm">
                                            <div className="border-bt-border border-b px-4 py-3">
                                                <p className="text-bt-dark text-xs font-semibold">
                                                    {user?.name}
                                                </p>
                                                <p className="text-bt-muted-light mt-0.5 truncate text-xs">
                                                    {user?.email}
                                                </p>
                                            </div>
                                            <button
                                                onClick={() => {
                                                    setShowUserMenu(false)
                                                    setShowLogoutModal(true)
                                                }}
                                                className="text-bt-muted hover:bg-bt-cream hover:text-bt-gold flex w-full cursor-pointer items-center gap-2.5 px-4 py-2.5 text-xs transition-colors"
                                            >
                                                <LogOut size={13} />
                                                Sign out
                                            </button>
                                        </div>
                                    </>
                                )}
                            </div>
                            <button
                                onClick={() => setMobileOpen(true)}
                                className="text-bt-warm cursor-pointer md:hidden"
                            >
                                <Menu size={20} />
                            </button>
                        </div>
                    </div>
                </div>
            </header>

            {mobileOpen && (
                <>
                    <div
                        className="bg-bt-dark/40 fixed inset-0 z-50 md:hidden"
                        onClick={() => setMobileOpen(false)}
                    />
                    <div className="bg-bt-surface fixed top-0 right-0 bottom-0 z-60 flex w-64 flex-col md:hidden">
                        <div className="border-bt-border flex items-center justify-between border-b px-5 py-4">
                            <span className="text-bt-dark font-serif text-base font-semibold">
                                Menu
                            </span>
                            <button
                                onClick={() => setMobileOpen(false)}
                                className="text-bt-muted cursor-pointer"
                            >
                                <X size={18} />
                            </button>
                        </div>
                        <div className="border-bt-border flex items-center gap-3 border-b px-5 py-4">
                            <div className="bg-bt-accent-bg text-bt-gold flex h-8 w-8 shrink-0 items-center justify-center rounded-full text-sm font-semibold">
                                {user?.name?.charAt(0).toUpperCase()}
                            </div>
                            <div className="min-w-0">
                                <p className="text-bt-dark truncate text-sm font-medium">
                                    {user?.name}
                                </p>
                                <p className="text-bt-muted-light truncate text-xs">
                                    {user?.email}
                                </p>
                            </div>
                        </div>
                        <nav className="flex-1 px-3 py-3">
                            {navItems.map((item) => {
                                const Icon = item.icon
                                return (
                                    <NavLink
                                        key={item.path}
                                        to={item.path}
                                        end={item.path === '/'}
                                        onClick={() => setMobileOpen(false)}
                                        className={({ isActive }) =>
                                            `mb-0.5 flex items-center gap-3 rounded-lg px-3 py-2.5 text-sm transition-all duration-150 ${isActive ? 'bg-bt-accent-bg text-bt-gold font-medium' : 'text-bt-warm hover:bg-bt-cream'}`
                                        }
                                    >
                                        <Icon size={16} strokeWidth={1.75} />
                                        {item.label}
                                    </NavLink>
                                )
                            })}
                        </nav>
                        <div className="border-bt-border border-t p-4">
                            <button
                                onClick={() => {
                                    setMobileOpen(false)
                                    setShowLogoutModal(true)
                                }}
                                className="text-bt-muted hover:bg-bt-cream hover:text-bt-gold flex w-full cursor-pointer items-center gap-2 rounded-lg px-3 py-2.5 text-sm transition-colors"
                            >
                                <LogOut size={15} />
                                Sign out
                            </button>
                        </div>
                    </div>
                </>
            )}

            <main className="mx-auto max-w-5xl px-6 py-8">{children}</main>

            {showLogoutModal && (
                <div className="bg-bt-dark/50 fixed inset-0 z-200 flex items-center justify-center">
                    <div className="bg-bt-surface mx-4 w-full max-w-sm overflow-hidden rounded-xl">
                        <div className="border-bt-border flex items-center justify-between border-b px-6 py-5">
                            <h2 className="text-bt-dark font-serif text-lg font-semibold">
                                Sign out?
                            </h2>
                            <button
                                onClick={() => setShowLogoutModal(false)}
                                className="text-bt-muted-light hover:text-bt-dark cursor-pointer transition-colors"
                            >
                                <X size={18} />
                            </button>
                        </div>
                        <div className="px-6 py-5">
                            <p className="text-bt-warm text-sm leading-relaxed">
                                Are you sure you want to sign out? You will need to sign in again to
                                access your library.
                            </p>
                        </div>
                        <div className="flex gap-3 px-6 pb-6">
                            <button
                                onClick={() => setShowLogoutModal(false)}
                                className="bg-bt-cream text-bt-warm border-bt-input-border hover:border-bt-gold flex-1 cursor-pointer rounded-md border py-2.5 text-sm font-medium transition-all duration-200"
                            >
                                Cancel
                            </button>
                            <button
                                onClick={handleLogout}
                                className="bg-bt-dark text-bt-cream hover:bg-bt-gold flex-1 cursor-pointer rounded-md py-2.5 text-sm font-medium transition-all duration-200"
                            >
                                Yes, sign out
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    )
}

export default Layout
