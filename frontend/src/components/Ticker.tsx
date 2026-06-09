import MarqueeComponent from 'react-fast-marquee'
import { TICKER_ITEMS } from '../constants/landingData'
import { Sparkle } from 'lucide-react'

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const Marquee = (MarqueeComponent as any).default || MarqueeComponent

export const Ticker = () => {
    return (
        <div className="bg-bt-dark border-bt-dark/10 border-y py-3">
            <Marquee speed={40} gradient={false}>
                {TICKER_ITEMS.map((item, index) => (
                    <span key={index} className="flex items-center">
                        <span className="text-bt-gold px-6 text-[10px] font-medium tracking-widest uppercase">
                            {item}
                        </span>
                        <Sparkle size={12} className="text-bt-gold/40 shrink-0" />
                    </span>
                ))}
            </Marquee>
        </div>
    )
}
