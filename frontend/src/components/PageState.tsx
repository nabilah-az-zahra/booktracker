import Layout from './Layout'

interface PageStateProps {
    loading?: boolean
    error?: boolean
    loadingText?: string
    onRetry?: () => void
}

const PageState = ({ loading, error, loadingText = 'Loading...', onRetry }: PageStateProps) => {
    if (loading) {
        return (
            <Layout>
                <div className="flex items-center justify-center py-20">
                    <p className="text-bt-muted-light text-sm">{loadingText}</p>
                </div>
            </Layout>
        )
    }

    if (error) {
        return (
            <Layout>
                <div className="flex flex-col items-center justify-center py-20 text-center">
                    <p className="text-bt-dark mb-2 font-medium">Something went wrong</p>
                    <p className="text-bt-muted mb-6 text-sm">
                        We couldn't load your data right now.
                    </p>
                    {onRetry && (
                        <button
                            onClick={onRetry}
                            className="bg-bt-dark text-bt-cream hover:bg-bt-gold rounded-md px-4 py-2 text-sm transition-colors"
                        >
                            Try again
                        </button>
                    )}
                </div>
            </Layout>
        )
    }

    return null
}

export default PageState
