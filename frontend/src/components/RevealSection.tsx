import { useReveal } from '../hooks/useReveal'

interface RevealSectionProps {
    children: React.ReactNode
    className?: string
}

export const RevealSection = ({ children, className = '' }: RevealSectionProps) => {
    const { ref, visible } = useReveal()

    return (
        <section
            ref={ref}
            className={`transition-opacity duration-700 ${
                visible ? 'is-revealed opacity-100' : 'opacity-0'
            } ${className}`}
        >
            {children}
        </section>
    )
}
