import axios from 'axios'

const api = axios.create({
    baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080',
    withCredentials: true,
})

api.interceptors.request.use((config) => {
    const token = localStorage.getItem('token')
    if (token) {
        config.headers.Authorization = `Bearer ${token}`
    }
    return config
})

let isRefreshing = false
let failedQueue: Array<{
    resolve: (token: string) => void
    reject: (error: unknown) => void
}> = []

const processQueue = (error: unknown, token: string | null = null) => {
    failedQueue.forEach(({ resolve, reject }) => {
        if (error) {
            reject(error)
        } else {
            resolve(token!)
        }
    })
    failedQueue = []
}

api.interceptors.response.use(
    (response) => response,
    async (error) => {
        const originalRequest = error.config
        if (error.response?.status !== 401 || originalRequest._retry) {
            return Promise.reject(error)
        }
        if (originalRequest.url?.includes('/auth/refresh')) {
            localStorage.removeItem('token')
            localStorage.removeItem('user')
            const currentPath = window.location.pathname
            window.location.href = `/login?redirect=${encodeURIComponent(currentPath)}`
            return Promise.reject(error)
        }
        if (isRefreshing) {
            return new Promise((resolve, reject) => {
                failedQueue.push({ resolve, reject })
            })
                .then((token) => {
                    originalRequest.headers.Authorization = `Bearer ${token}`
                    return api(originalRequest)
                })
                .catch((err) => {
                    return Promise.reject(err)
                })
        }
        originalRequest._retry = true
        isRefreshing = true

        try {
            const res = await api.post('/api/auth/refresh')
            const newToken = res.data.token
            localStorage.setItem('token', newToken)
            localStorage.setItem('user', JSON.stringify(res.data.user))
            api.defaults.headers.common['Authorization'] = `Bearer ${newToken}`
            originalRequest.headers.Authorization = `Bearer ${newToken}`
            processQueue(null, newToken)
            return api(originalRequest)
        } catch (refreshError) {
            processQueue(refreshError, null)
            localStorage.removeItem('token')
            localStorage.removeItem('user')
            const currentPath = window.location.pathname
            window.location.href = `/login?redirect=${encodeURIComponent(currentPath)}`
            return Promise.reject(refreshError)
        } finally {
            isRefreshing = false
        }
    },
)

export default api
