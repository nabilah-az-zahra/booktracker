export const statusLabel = (status: string) => {
    if (status === 'want_to_read') return 'Want to read'
    if (status === 'reading') return 'Reading'
    return 'Finished'
}

export const statusClassName = (status: string) => {
    if (status === 'finished') return 'bg-bt-status-finished-bg text-bt-status-finished-text'
    if (status === 'reading') return 'bg-bt-accent-bg text-bt-gold'
    return 'bg-bt-status-want-bg text-bt-muted'
}
