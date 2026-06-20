export const statusLabel = (status: string) => {
    if (status === 'want_to_read') return 'Want to read'
    if (status === 'reading') return 'Reading'
    return 'Finished'
}

export const statusClassName = (status: string) => {
    if (status === 'finished') return 'bg-bt-status-finished-bg text-bt-status-finished-text'
    if (status === 'reading') return 'bg-bt-accent-bg text-bt-gold'
    return 'bg-bt-accent-bg text-bt-muted'
}

export const statusBorderClassName = (status: string) => {
    if (status === 'finished') return 'border-l-bt-status-finished-text'
    if (status === 'reading') return 'border-l-bt-gold'
    return 'border-l-bt-border'
}

export const statusTextClassName = (status: string) => {
    if (status === 'finished') return 'text-bt-status-finished-text'
    if (status === 'reading') return 'text-bt-gold'
    return 'text-bt-muted'
}
