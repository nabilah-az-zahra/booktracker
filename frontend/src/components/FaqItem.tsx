import { useState } from 'react'
import { Plus, Minus } from 'lucide-react'

interface FaqItemProps {
    q: string
    a: string
}

export const FaqItem = ({ q, a }: FaqItemProps) => {
    const [open, setOpen] = useState(false)

    return (
        <div
            className="border-bt-border cursor-pointer border-b"
            onClick={() => setOpen((o) => !o)}
        >
            <div className="flex items-center justify-between gap-4 py-5">
                <span
                    className={`font-serif text-sm font-medium transition-colors duration-200 ${open ? 'text-bt-gold' : 'text-bt-dark'}`}
                >
                    {q}
                </span>
                <span
                    className={`shrink-0 transition-colors duration-200 ${open ? 'text-bt-gold' : 'text-bt-muted'}`}
                >
                    {open ? (
                        <Minus size={14} strokeWidth={1.5} />
                    ) : (
                        <Plus size={14} strokeWidth={1.5} />
                    )}
                </span>
            </div>
            <div
                className={`grid transition-all duration-500 ${open ? 'grid-rows-[1fr]' : 'grid-rows-[0fr]'}`}
            >
                <div className="overflow-hidden">
                    <p className="text-bt-muted pb-5 font-serif text-sm leading-relaxed">{a}</p>
                </div>
            </div>
        </div>
    )
}
