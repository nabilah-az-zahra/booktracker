import { Routes, Route, Navigate } from 'react-router-dom'
import { useAuth } from './context/useAuth'
import type { ReactNode } from 'react'

import Landing from './pages/Landing'
import Login from './pages/Login'
import Register from './pages/Register'
import Dashboard from './pages/Dashboard'
import Library from './pages/Library'
import BookDetail from './pages/BookDetail'
import BookSearch from './pages/BookSearch'
import ReadingTimer from './pages/ReadingTimer'
import Stats from './pages/Stats'

const ProtectedRoute = ({ children }: { children: ReactNode }) => {
    const { isAuthenticated, isLoading } = useAuth()
    if (isLoading) return null
    return isAuthenticated ? <>{children}</> : <Navigate to="/welcome" replace />
}

const GuestRoute = ({ children }: { children: ReactNode }) => {
    const { isAuthenticated, isLoading } = useAuth()
    if (isLoading) return null
    return isAuthenticated ? <Navigate to="/" replace /> : <>{children}</>
}

const App = () => {
    return (
        <Routes>
            <Route
                path="/welcome"
                element={
                    <GuestRoute>
                        <Landing />
                    </GuestRoute>
                }
            />
            <Route
                path="/login"
                element={
                    <GuestRoute>
                        <Login />
                    </GuestRoute>
                }
            />
            <Route
                path="/register"
                element={
                    <GuestRoute>
                        <Register />
                    </GuestRoute>
                }
            />
            <Route
                path="/"
                element={
                    <ProtectedRoute>
                        <Dashboard />
                    </ProtectedRoute>
                }
            />
            <Route
                path="/library"
                element={
                    <ProtectedRoute>
                        <Library />
                    </ProtectedRoute>
                }
            />
            <Route
                path="/books/search"
                element={
                    <ProtectedRoute>
                        <BookSearch />
                    </ProtectedRoute>
                }
            />
            <Route
                path="/books/:id"
                element={
                    <ProtectedRoute>
                        <BookDetail />
                    </ProtectedRoute>
                }
            />
            <Route
                path="/reading/:id"
                element={
                    <ProtectedRoute>
                        <ReadingTimer />
                    </ProtectedRoute>
                }
            />
            <Route
                path="/stats"
                element={
                    <ProtectedRoute>
                        <Stats />
                    </ProtectedRoute>
                }
            />
            <Route path="*" element={<Navigate to="/welcome" replace />} />
        </Routes>
    )
}

export default App
