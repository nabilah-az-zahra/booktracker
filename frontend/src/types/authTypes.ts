import type { User } from './index'

export interface AuthContextType {
    user: User | null
    token: string | null
    login: (token: string, user: User) => void
    logout: () => void
    isAuthenticated: boolean
    isLoading: boolean
}
