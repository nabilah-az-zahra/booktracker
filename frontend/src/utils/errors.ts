export const getApiError = (err: unknown): string => {
    if (err instanceof Object && 'response' in err) {
        const message = (err as { response?: { data?: { message?: string } } }).response?.data
            ?.message
        if (message) return message
    }
    return 'Something went wrong. Please try again.'
}
