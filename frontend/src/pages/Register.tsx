import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import api from '../api/axios'
import { Eye, EyeOff, BookOpen } from 'lucide-react'
import type { AuthResponse } from '../types'
import { useAuth } from '../context/useAuth'
import { getApiError } from '../utils/errors'

const Register = () => {
    const { login } = useAuth()
    const navigate = useNavigate()
    const [formData, setFormData] = useState({
        name: '',
        email: '',
        password: '',
        confirmPassword: '',
    })
    const [error, setError] = useState('')
    const [loading, setLoading] = useState(false)
    const [showPassword, setShowPassword] = useState(false)

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setFormData({ ...formData, [e.target.name]: e.target.value })
    }

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault()
        setError('')

        if (formData.name.trim() === '') {
            setError('Please enter your full name')
            return
        }
        if (formData.password.length < 8) {
            setError('Password must be at least 8 characters')
            return
        }
        if (formData.password !== formData.confirmPassword) {
            setError('Passwords do not match')
            return
        }

        setLoading(true)
        try {
            const res = await api.post<AuthResponse>('/api/auth/register', {
                name: formData.name,
                email: formData.email,
                password: formData.password,
                timezone: Intl.DateTimeFormat().resolvedOptions().timeZone,
            })
            login(res.data.token, res.data.user)
            navigate('/')
        } catch (err) {
            setError(getApiError(err))
        } finally {
            setLoading(false)
        }
    }

    return (
        <div className="bg-bt-cream flex min-h-screen flex-col items-center justify-center px-6 py-16">
            <Link to="/welcome" className="group mb-10 flex items-center gap-3">
                <div className="border-bt-placeholder group-hover:border-bt-gold relative flex h-10 w-10 items-center justify-center rounded-xl border-2 border-dashed bg-transparent transition-all duration-300 group-hover:rotate-12">
                    <div className="relative">
                        <BookOpen
                            size={18}
                            className="text-bt-dark transition-colors duration-300"
                            strokeWidth={1.5}
                        />
                        <span className="bg-bt-dark group-hover:bg-bt-gold absolute bottom-1 left-1/2 h-1.5 w-1.5 -translate-x-1/2 rounded-full transition-colors duration-300" />
                    </div>
                </div>
                <span className="text-bt-dark font-serif text-xl font-normal tracking-wide">
                    <strong className="font-bold">Book</strong>
                    <span className="text-bt-warm group-hover:text-bt-gold ml-0.5 font-sans text-lg font-semibold tracking-wider uppercase transition-colors duration-300">
                        Tracker
                    </span>
                </span>
            </Link>

            <div className="bg-bt-surface border-bt-border w-full max-w-sm rounded-2xl border px-8 py-10">
                <div className="mb-8 text-center">
                    <h2 className="text-bt-dark mb-1.5 font-serif text-2xl font-semibold">
                        Create your account
                    </h2>
                    <p className="text-bt-muted text-sm">Start your reading journey today</p>
                </div>

                {error && (
                    <div className="bg-bt-danger-bg text-bt-danger mb-5 rounded-lg px-4 py-3 text-sm">
                        {error}
                    </div>
                )}

                <form onSubmit={handleSubmit} className="space-y-4">
                    <div>
                        <label className="text-bt-muted mb-2 block text-xs font-medium tracking-widest uppercase">
                            Full Name
                        </label>
                        <input
                            type="text"
                            name="name"
                            value={formData.name}
                            onChange={handleChange}
                            required
                            placeholder="Your name"
                            className="input-field w-full px-4 py-2.5 text-sm transition-all duration-200"
                        />
                    </div>

                    <div>
                        <label className="text-bt-muted mb-2 block text-xs font-medium tracking-widest uppercase">
                            Email
                        </label>
                        <input
                            type="email"
                            name="email"
                            value={formData.email}
                            onChange={handleChange}
                            required
                            placeholder="your@email.com"
                            className="input-field w-full px-4 py-2.5 text-sm transition-all duration-200"
                        />
                    </div>

                    <div>
                        <label className="text-bt-muted mb-2 block text-xs font-medium tracking-widest uppercase">
                            Password
                        </label>
                        <div className="relative">
                            <input
                                type={showPassword ? 'text' : 'password'}
                                name="password"
                                value={formData.password}
                                onChange={handleChange}
                                required
                                placeholder="Min 8 characters"
                                className="input-field w-full px-4 py-2.5 pr-10 text-sm transition-all duration-200"
                            />
                            <button
                                type="button"
                                onClick={() => setShowPassword(!showPassword)}
                                className="text-bt-muted-light hover:text-bt-muted-dark absolute top-1/2 right-3 -translate-y-1/2 cursor-pointer transition-colors"
                            >
                                {showPassword ? <EyeOff size={15} /> : <Eye size={15} />}
                            </button>
                        </div>
                    </div>

                    <div>
                        <label className="text-bt-muted mb-2 block text-xs font-medium tracking-widest uppercase">
                            Confirm Password
                        </label>
                        <input
                            type="password"
                            name="confirmPassword"
                            value={formData.confirmPassword}
                            onChange={handleChange}
                            required
                            placeholder="Repeat your password"
                            className="input-field w-full px-4 py-2.5 text-sm transition-all duration-200"
                        />
                    </div>

                    <button
                        type="submit"
                        disabled={loading}
                        className="bg-bt-dark text-bt-cream hover:bg-bt-gold mt-2 w-full cursor-pointer rounded-lg py-2.5 text-sm font-medium transition-all duration-200 disabled:pointer-events-none disabled:opacity-50"
                    >
                        {loading ? 'Creating account...' : 'Create account'}
                    </button>
                </form>

                <p className="text-bt-muted mt-6 text-center text-sm">
                    Already have an account?{' '}
                    <Link
                        to="/login"
                        className="text-bt-dark hover:text-bt-gold font-medium transition-colors"
                    >
                        Sign in
                    </Link>
                </p>
            </div>

            <p className="text-bt-footer mt-8 text-xs">© 2026 BookTracker</p>
        </div>
    )
}

export default Register
