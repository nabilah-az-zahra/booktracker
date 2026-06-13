import { useEffect, useReducer, useState, useCallback } from 'react'
import type { ReactNode } from 'react'
import type { User } from '../types'
import { AuthContext } from './AuthContext'
import type { AuthContextType } from '../types/authTypes'
import api from '../api/axios'

interface AuthState {
    user: User | null
    token: string | null
}

type AuthAction =
    | { type: 'LOGIN'; token: string; user: User }
    | { type: 'LOGOUT' }
    | { type: 'RESTORE'; token: string; user: User }

const authReducer = (state: AuthState, action: AuthAction): AuthState => {
    switch (action.type) {
        case 'LOGIN':
        case 'RESTORE':
            return { token: action.token, user: action.user }
        case 'LOGOUT':
            return { token: null, user: null }
        default:
            return state
    }
}

export const AuthProvider = ({ children }: { children: ReactNode }) => {
    const [state, dispatch] = useReducer(authReducer, {
        user: null,
        token: null,
    })
    const [isReady, setIsReady] = useState(false)

    useEffect(() => {
        try {
            const savedToken = localStorage.getItem('token')
            const savedUser = localStorage.getItem('user')
            if (savedToken && savedUser) {
                dispatch({
                    type: 'RESTORE',
                    token: savedToken,
                    user: JSON.parse(savedUser) as User,
                })
            }
        } catch {
            localStorage.removeItem('token')
            localStorage.removeItem('user')
        }
        setTimeout(() => setIsReady(true), 0)
    }, [])

    const login = useCallback((token: string, user: User) => {
        localStorage.setItem('token', token)
        localStorage.setItem('user', JSON.stringify(user))
        dispatch({ type: 'LOGIN', token, user })
    }, [])

    const logout = useCallback(async () => {
        try {
            await api.post('/api/auth/logout')
        } catch {
            //
        }
        localStorage.removeItem('token')
        localStorage.removeItem('user')
        dispatch({ type: 'LOGOUT' })
    }, [])

    const value: AuthContextType = {
        user: state.user,
        token: state.token,
        login,
        logout,
        isAuthenticated: !!state.token,
        isLoading: !isReady,
    }

    return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}
