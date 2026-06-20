import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import Layout from '../components/Layout'
import PageState from '../components/PageState'
import api from '../api/axios'
import type { DailyReading, StatsData } from '../types'
import {
    BookOpen,
    Clock,
    Flame,
    Target,
    TrendingUp,
    Award,
    BookMarked,
    ArrowRight,
    BarChart2,
    Search,
} from 'lucide-react'
import { formatTime } from '../utils/formatUtils'
import { BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer, CartesianGrid } from 'recharts'

type HistoryWindow = 7 | 30 | 90
type HistoryMetric = 'pages' | 'minutes'

const Stats = () => {
    const [stats, setStats] = useState<StatsData | null>(null)
    const [history, setHistory] = useState<DailyReading[]>([])
    const [historyWindow, setHistoryWindow] = useState<HistoryWindow>(30)
    const [historyMetric, setHistoryMetric] = useState<HistoryMetric>('pages')
    const [historyLoading, setHistoryLoading] = useState(false)
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState(false)

    const fetchStats = async () => {
        try {
            setError(false)
            const res = await api.get('/api/stats')
            setStats(res.data.data)
        } catch (err) {
            console.error(err)
            setError(true)
        } finally {
            setLoading(false)
        }
    }

    const fetchHistory = async (days: HistoryWindow) => {
        setHistoryLoading(true)
        try {
            const res = await api.get(`/api/stats/history?days=${days}`)
            setHistory(res.data.data || [])
        } catch (err) {
            console.error(err)
        } finally {
            setHistoryLoading(false)
        }
    }

    useEffect(() => {
        // eslint-disable-next-line react-hooks/set-state-in-effect
        fetchStats()
    }, [])

    useEffect(() => {
        // eslint-disable-next-line react-hooks/set-state-in-effect
        fetchHistory(historyWindow)
    }, [historyWindow])

    const handleWindowChange = (days: HistoryWindow) => {
        setHistoryWindow(days)
    }

    const yearlyProgress = stats?.yearly_goal
        ? Math.min(Math.round((stats.yearly_finished / stats.yearly_goal) * 100), 100)
        : 0

    if (loading || error) {
        return (
            <PageState
                loading={loading}
                error={error}
                loadingText="Loading stats..."
                onRetry={fetchStats}
            />
        )
    }

    if (!stats)
        return (
            <Layout>
                <div className="py-20 text-center">
                    <p className="text-bt-muted text-sm">Failed to load stats</p>
                </div>
            </Layout>
        )

    const finishedPct =
        stats.total_books > 0 ? Math.round((stats.finished_books / stats.total_books) * 100) : 0
    const readingPct =
        stats.total_books > 0
            ? Math.round(((stats.total_books - stats.finished_books) / stats.total_books) * 100)
            : 0

    const chartData = history.map((d) => ({
        date: d.date.slice(5),
        value: historyMetric === 'pages' ? d.pages : Math.round(d.seconds / 60),
        fullDate: d.date,
    }))

    const hasAnyData = history.some((d) => d.pages > 0 || d.seconds > 0)
    const CustomTooltip = ({
        active,
        payload,
        label,
    }: {
        active?: boolean
        payload?: Array<{ value: number }>
        label?: string
    }) => {
        if (active && payload && payload.length) {
            return (
                <div className="bg-bt-surface border-bt-border rounded-lg border px-3 py-2 shadow-sm">
                    <p className="text-bt-muted mb-1 text-xs">{label}</p>
                    <p className="text-bt-dark text-sm font-semibold">
                        {payload[0].value}{' '}
                        <span className="text-bt-muted font-normal">
                            {historyMetric === 'pages' ? 'pages' : 'min'}
                        </span>
                    </p>
                </div>
            )
        }
        return null
    }

    return (
        <Layout>
            <div className="mb-8">
                <h1 className="text-bt-dark mb-1 font-serif text-2xl font-semibold">Statistics</h1>
                <p className="text-bt-muted text-sm">Your reading in numbers</p>
            </div>

            {stats.total_books === 0 ? (
                <div className="border-bt-dashed bg-bt-surface rounded-xl border border-dashed py-16 text-center">
                    <TrendingUp
                        size={32}
                        strokeWidth={1.5}
                        className="text-bt-placeholder mx-auto mb-4"
                    />
                    <p className="text-bt-dark mb-2 font-serif text-lg font-semibold">
                        No stats yet
                    </p>
                    <p className="text-bt-muted mb-6 text-sm">
                        Add books and start reading sessions to see your stats here
                    </p>
                    <Link
                        to="/books/search"
                        className="text-bt-gold hover:text-bt-muted-dark text-sm font-medium transition-colors"
                    >
                        Find your first book
                        <Search
                            size={13}
                            strokeWidth={2}
                            className="ml-1 inline-block align-text-bottom"
                        />
                    </Link>
                </div>
            ) : (
                <div className="space-y-6">
                    <div className="grid grid-cols-2 gap-3 md:grid-cols-4">
                        <div className="border-bt-border bg-bt-surface hover:border-bt-gold rounded-xl border px-5 py-5 transition-colors duration-200">
                            <BookOpen size={16} strokeWidth={1.5} className="text-bt-gold mb-3" />
                            <p className="text-bt-dark mb-1 font-serif text-3xl font-semibold">
                                {stats.finished_books}
                            </p>
                            <p className="text-bt-warm mb-0.5 text-xs font-medium">Books Read</p>
                            <p className="text-bt-muted-light text-xs">
                                of {stats.total_books} total
                            </p>
                        </div>
                        <div className="border-bt-border bg-bt-surface hover:border-bt-gold rounded-xl border px-5 py-5 transition-colors duration-200">
                            <Clock size={16} strokeWidth={1.5} className="text-bt-gold mb-3" />
                            <p className="text-bt-dark mb-1 font-serif text-3xl font-semibold">
                                {formatTime(stats.total_reading_time_seconds)}
                            </p>
                            <p className="text-bt-warm mb-0.5 text-xs font-medium">Reading Time</p>
                            <p className="text-bt-muted-light text-xs">across all sessions</p>
                        </div>
                        <div className="border-bt-border bg-bt-surface hover:border-bt-gold rounded-xl border px-5 py-5 transition-colors duration-200">
                            <BookMarked size={16} strokeWidth={1.5} className="text-bt-gold mb-3" />
                            <p className="text-bt-dark mb-1 font-serif text-3xl font-semibold">
                                {stats.total_pages_read.toLocaleString()}
                            </p>
                            <p className="text-bt-warm mb-0.5 text-xs font-medium">Pages Read</p>
                            <p className="text-bt-muted-light text-xs">total pages</p>
                        </div>
                        <div className="border-bt-border bg-bt-surface hover:border-bt-gold rounded-xl border px-5 py-5 transition-colors duration-200">
                            <Flame size={16} strokeWidth={1.5} className="text-bt-gold mb-3" />
                            <p className="text-bt-dark mb-1 font-serif text-3xl font-semibold">
                                {stats.current_streak}d
                            </p>
                            <p className="text-bt-warm mb-0.5 text-xs font-medium">
                                Current Streak
                            </p>
                            <p className="text-bt-muted-light text-xs">
                                {stats.current_streak > 0 ? 'keep it up' : 'start today'}
                            </p>
                        </div>
                    </div>

                    <div className="border-bt-border bg-bt-surface rounded-xl border p-6">
                        <div className="mb-5 flex items-start justify-between">
                            <div className="flex items-center gap-3">
                                <Target size={18} strokeWidth={1.5} className="text-bt-gold" />
                                <div>
                                    <p className="text-bt-dark font-serif text-base font-semibold">
                                        {new Date().getFullYear()} Goal
                                    </p>
                                    <p className="text-bt-muted mt-0.5 text-xs">
                                        {stats.yearly_goal
                                            ? `${stats.yearly_finished} of ${stats.yearly_goal} books`
                                            : 'No goal set'}
                                    </p>
                                </div>
                            </div>
                            {stats.yearly_goal > 0 && (
                                <span
                                    className={`font-serif text-sm font-semibold ${yearlyProgress >= 100 ? 'text-bt-status-finished-text' : 'text-bt-gold'}`}
                                >
                                    {yearlyProgress}%
                                </span>
                            )}
                        </div>

                        {stats.yearly_goal > 0 ? (
                            <>
                                <div className="bg-bt-track mb-3 h-2 w-full rounded-full">
                                    <div
                                        className={`h-full rounded-full transition-all duration-700 ${yearlyProgress >= 100 ? 'bg-bt-status-finished-text' : 'bg-bt-gold'}`}
                                        style={{ width: `${yearlyProgress}%` }}
                                    />
                                </div>
                                {yearlyProgress >= 100 ? (
                                    <div className="flex items-center gap-2">
                                        <Award size={14} className="text-bt-status-finished-text" />
                                        <p className="text-bt-status-finished-text text-xs font-medium">
                                            Goal reached.
                                        </p>
                                    </div>
                                ) : (
                                    <p className="text-bt-muted-light text-xs">
                                        {stats.yearly_goal - stats.yearly_finished} book
                                        {stats.yearly_goal - stats.yearly_finished !== 1
                                            ? 's'
                                            : ''}{' '}
                                        to go
                                    </p>
                                )}
                            </>
                        ) : (
                            <div className="border-bt-dashed bg-bt-cream flex items-center justify-between rounded-lg border border-dashed px-4 py-3">
                                <p className="text-bt-muted text-sm">
                                    Set a yearly goal to track progress
                                </p>
                                <GoalSetter
                                    onSaved={(goal) =>
                                        setStats((prev) =>
                                            prev ? { ...prev, yearly_goal: goal } : prev,
                                        )
                                    }
                                />
                            </div>
                        )}
                    </div>

                    <div className="border-bt-border bg-bt-surface rounded-xl border p-6">
                        <div className="mb-5 flex items-center gap-3">
                            <TrendingUp size={18} strokeWidth={1.5} className="text-bt-gold" />
                            <p className="text-bt-dark font-serif text-base font-semibold">
                                Library Breakdown
                            </p>
                        </div>
                        <div className="space-y-4">
                            <div>
                                <div className="mb-1.5 flex items-center justify-between">
                                    <p className="text-bt-warm text-xs font-medium">Finished</p>
                                    <p className="text-bt-muted text-xs">
                                        {stats.finished_books} book
                                        {stats.finished_books !== 1 ? 's' : ''} ({finishedPct}%)
                                    </p>
                                </div>
                                <div className="bg-bt-track h-1.5 w-full rounded-full">
                                    <div
                                        className="bg-bt-status-finished-text h-full rounded-full transition-all duration-700"
                                        style={{ width: `${finishedPct}%` }}
                                    />
                                </div>
                            </div>
                            <div>
                                <div className="mb-1.5 flex items-center justify-between">
                                    <p className="text-bt-warm text-xs font-medium">
                                        Currently Reading
                                    </p>
                                    <p className="text-bt-muted text-xs">
                                        {stats.total_books - stats.finished_books} book
                                        {stats.total_books - stats.finished_books !== 1
                                            ? 's'
                                            : ''}{' '}
                                        ({readingPct}%)
                                    </p>
                                </div>
                                <div className="bg-bt-track h-1.5 w-full rounded-full">
                                    <div
                                        className="bg-bt-gold h-full rounded-full transition-all duration-700"
                                        style={{ width: `${readingPct}%` }}
                                    />
                                </div>
                            </div>
                        </div>
                    </div>

                    <div className="border-bt-border bg-bt-surface rounded-xl border p-6">
                        <div className="mb-5 flex items-start justify-between">
                            <div className="flex items-center gap-3">
                                <BarChart2 size={18} strokeWidth={1.5} className="text-bt-gold" />
                                <p className="text-bt-dark font-serif text-base font-semibold">
                                    Reading History
                                </p>
                            </div>
                            <div className="flex items-center gap-2">
                                <div className="bg-bt-accent-bg flex gap-1 rounded-lg p-1">
                                    <button
                                        onClick={() => setHistoryMetric('pages')}
                                        className={`cursor-pointer rounded-md px-2.5 py-1 text-xs font-medium transition-colors duration-200 ${
                                            historyMetric === 'pages'
                                                ? 'bg-bt-surface text-bt-dark shadow-sm'
                                                : 'text-bt-muted'
                                        }`}
                                    >
                                        Pages
                                    </button>
                                    <button
                                        onClick={() => setHistoryMetric('minutes')}
                                        className={`cursor-pointer rounded-md px-2.5 py-1 text-xs font-medium transition-colors duration-200 ${
                                            historyMetric === 'minutes'
                                                ? 'bg-bt-surface text-bt-dark shadow-sm'
                                                : 'text-bt-muted'
                                        }`}
                                    >
                                        Minutes
                                    </button>
                                </div>
                                <div className="bg-bt-accent-bg flex gap-1 rounded-lg p-1">
                                    {([7, 30, 90] as HistoryWindow[]).map((d) => (
                                        <button
                                            key={d}
                                            onClick={() => handleWindowChange(d)}
                                            className={`cursor-pointer rounded-md px-2.5 py-1 text-xs font-medium transition-colors duration-200 ${
                                                historyWindow === d
                                                    ? 'bg-bt-surface text-bt-dark shadow-sm'
                                                    : 'text-bt-muted'
                                            }`}
                                        >
                                            {d}d
                                        </button>
                                    ))}
                                </div>
                            </div>
                        </div>

                        {historyLoading ? (
                            <div className="flex h-48 items-center justify-center">
                                <p className="text-bt-muted-light text-sm">Loading...</p>
                            </div>
                        ) : !hasAnyData ? (
                            <div className="flex h-48 flex-col items-center justify-center">
                                <BarChart2
                                    size={24}
                                    strokeWidth={1.5}
                                    className="text-bt-placeholder mb-3"
                                />
                                <p className="text-bt-muted text-sm">
                                    No reading sessions in the last {historyWindow} days
                                </p>
                            </div>
                        ) : (
                            <ResponsiveContainer width="100%" height={200}>
                                <BarChart
                                    data={chartData}
                                    margin={{ top: 4, right: 4, left: -24, bottom: 0 }}
                                >
                                    <CartesianGrid
                                        strokeDasharray="3 3"
                                        stroke="var(--color-bt-border)"
                                        vertical={false}
                                    />
                                    <XAxis
                                        dataKey="date"
                                        tick={{ fontSize: 10, fill: 'var(--color-bt-muted)' }}
                                        tickLine={false}
                                        axisLine={false}
                                        interval={
                                            historyWindow === 7 ? 0 : historyWindow === 30 ? 6 : 14
                                        }
                                    />
                                    <YAxis
                                        tick={{ fontSize: 10, fill: 'var(--color-bt-muted)' }}
                                        tickLine={false}
                                        axisLine={false}
                                    />
                                    <Tooltip
                                        // eslint-disable-next-line react-hooks/static-components
                                        content={<CustomTooltip />}
                                        cursor={{ fill: 'var(--color-bt-order)', opacity: 0.5 }}
                                    />
                                    <Bar
                                        dataKey="value"
                                        fill="var(--color-bt-gold)"
                                        radius={[3, 3, 0, 0]}
                                        maxBarSize={
                                            historyWindow === 7 ? 40 : historyWindow === 30 ? 16 : 8
                                        }
                                    />
                                </BarChart>
                            </ResponsiveContainer>
                        )}
                    </div>

                    {stats.total_reading_time_seconds > 0 && stats.total_pages_read > 0 && (
                        <div className="border-bt-border bg-bt-surface rounded-xl border p-6">
                            <div className="mb-4 flex items-center gap-3">
                                <TrendingUp size={18} strokeWidth={1.5} className="text-bt-gold" />
                                <p className="text-bt-dark font-serif text-base font-semibold">
                                    Reading Speed
                                </p>
                            </div>
                            <div className="grid grid-cols-2 gap-4">
                                <div className="bg-bt-cream rounded-xl px-4 py-4 text-center">
                                    <p className="text-bt-dark mb-1 font-serif text-3xl font-semibold">
                                        {Math.round(
                                            stats.total_pages_read /
                                                (stats.total_reading_time_seconds / 3600),
                                        )}
                                    </p>
                                    <p className="text-bt-muted text-xs">pages per hour</p>
                                </div>
                                <div className="bg-bt-cream rounded-xl px-4 py-4 text-center">
                                    <p className="text-bt-dark mb-1 font-serif text-3xl font-semibold">
                                        {Math.round(
                                            stats.total_reading_time_seconds /
                                                60 /
                                                Math.max(stats.total_pages_read, 1),
                                        )}
                                        m
                                    </p>
                                    <p className="text-bt-muted text-xs">per page</p>
                                </div>
                            </div>
                        </div>
                    )}

                    <div className="bg-bt-dark flex items-center justify-between rounded-xl px-6 py-5">
                        <div>
                            <p className="mb-0.5 font-serif text-base font-semibold text-white">
                                {stats.current_streak > 0
                                    ? `${stats.current_streak} days running`
                                    : 'Nothing in progress'}
                            </p>
                            <p className="text-bt-muted-dark text-xs">
                                {stats.current_streak > 0
                                    ? "Don't break the chain."
                                    : 'Pick something up.'}
                            </p>
                        </div>
                        <Link
                            to="/library"
                            className="bg-bt-gold text-bt-dark hover:bg-bt-cream-light flex items-center gap-2 rounded-md px-4 py-2 text-xs font-medium transition-colors duration-200"
                        >
                            Library <ArrowRight size={12} />
                        </Link>
                    </div>
                </div>
            )}
        </Layout>
    )
}

const GoalSetter = ({ onSaved }: { onSaved: (goal: number) => void }) => {
    const [editing, setEditing] = useState(false)
    const [value, setValue] = useState('')
    const [saving, setSaving] = useState(false)
    const [saveError, setSaveError] = useState('')

    const handleSave = async () => {
        const goal = parseInt(value)
        if (isNaN(goal) || goal < 1) return
        setSaving(true)
        setSaveError('')
        try {
            await api.patch('/api/profile/goal', { yearly_goal: goal })
            onSaved(goal)
            setEditing(false)
        } catch (err) {
            console.error(err)
            setSaveError('Failed to save goal')
        } finally {
            setSaving(false)
        }
    }

    if (!editing)
        return (
            <button
                onClick={() => {
                    setEditing(true)
                    setSaveError('')
                }}
                className="text-bt-gold hover:text-bt-muted-dark shrink-0 cursor-pointer text-xs font-medium transition-colors"
            >
                Set goal
                <Target size={13} strokeWidth={2} className="mb-0.5 ml-1 inline-block" />
            </button>
        )

    return (
        <div className="flex flex-col items-end">
            <div className="flex shrink-0 items-center gap-2">
                <input
                    type="number"
                    value={value}
                    onChange={(e) => setValue(e.target.value)}
                    placeholder="e.g. 12"
                    min="1"
                    autoFocus
                    className="input-field w-20 px-2 py-1 text-center text-xs transition-colors duration-200"
                    onKeyDown={(e) => e.key === 'Enter' && handleSave()}
                />
                <button
                    onClick={handleSave}
                    disabled={saving || !value}
                    className="text-bt-gold hover:text-bt-muted-dark cursor-pointer text-xs font-medium transition-colors disabled:opacity-50"
                >
                    {saving ? 'Saving...' : 'Save'}
                </button>
                <button
                    onClick={() => setEditing(false)}
                    className="text-bt-muted-light hover:text-bt-warm cursor-pointer text-xs transition-colors"
                >
                    Cancel
                </button>
            </div>
            {saveError && <p className="text-bt-danger mt-1 text-xs">{saveError}</p>}
        </div>
    )
}

export default Stats
