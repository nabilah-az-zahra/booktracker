export const formatTime = (seconds: number): string => {
    const hours = Math.floor(seconds / 3600)
    const minutes = Math.floor((seconds % 3600) / 60)
    const secs = seconds % 60

    const mm = String(minutes).padStart(2, '0')
    const ss = String(secs).padStart(2, '0')

    if (hours > 0) {
        const hh = String(hours).padStart(2, '0')
        return `${hh}:${mm}:${ss}`
    }

    return `${mm}:${ss}`
}
