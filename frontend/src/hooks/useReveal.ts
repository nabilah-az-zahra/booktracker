import { useState, useEffect, useRef } from 'react'

export const useReveal = () => {
    const ref = useRef<HTMLDivElement>(null)
    const [visible, setVisible] = useState(false)

    useEffect(() => {
        const el = ref.current
        if (!el) return

        const obs = new IntersectionObserver(
            ([entry]) => {
                if (entry.isIntersecting) setVisible(true)
            },
            { threshold: 0.1 },
        )

        obs.observe(el)
        return () => obs.disconnect()
    }, [])

    return { ref, visible }
}
