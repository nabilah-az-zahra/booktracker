import { useState, useRef, useCallback } from 'react'

interface UseTimerReturn {
    seconds: number
    isRunning: boolean
    isPaused: boolean
    start: () => void
    pause: () => void
    resume: () => void
    reset: () => void
    restore: (savedSeconds: number) => void
}

export const useTimer = (initialSeconds: number = 0): UseTimerReturn => {
    const [seconds, setSeconds] = useState(initialSeconds)
    const [isRunning, setIsRunning] = useState(false)
    const [isPaused, setIsPaused] = useState(false)
    const intervalRef = useRef<ReturnType<typeof setInterval> | null>(null)

    const start = useCallback(() => {
        if (intervalRef.current) return
        setIsRunning(true)
        setIsPaused(false)
        intervalRef.current = setInterval(() => {
            setSeconds((prev) => prev + 1)
        }, 1000)
    }, [])

    const pause = useCallback(() => {
        if (intervalRef.current) {
            clearInterval(intervalRef.current)
            intervalRef.current = null
        }
        setIsRunning(false)
        setIsPaused(true)
    }, [])

    const resume = start

    const reset = useCallback(() => {
        if (intervalRef.current) {
            clearInterval(intervalRef.current)
            intervalRef.current = null
        }
        setSeconds(0)
        setIsRunning(false)
        setIsPaused(false)
    }, [])

    const restore = useCallback((savedSeconds: number) => {
        if (intervalRef.current) {
            clearInterval(intervalRef.current)
            intervalRef.current = null
        }
        setSeconds(savedSeconds)
        setIsRunning(false)
        setIsPaused(true)
    }, [])

    return { seconds, isRunning, isPaused, start, pause, resume, reset, restore }
}
