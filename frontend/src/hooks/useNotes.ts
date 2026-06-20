import { useState, useCallback } from 'react'
import api from '../api/axios'
import type { SessionNote } from '../types'

interface NoteForm {
    chapter: string
    pages: string
    note: string
}

const emptyForm: NoteForm = { chapter: '', pages: '', note: '' }

export const useNotes = (sessionId: string | null) => {
    const [notes, setNotes] = useState<SessionNote[]>([])
    const [form, setForm] = useState<NoteForm>(emptyForm)
    const [submitting, setSubmitting] = useState(false)
    const [error, setError] = useState('')

    const updateForm = useCallback((field: keyof NoteForm, value: string) => {
        setForm((prev) => ({ ...prev, [field]: value }))
    }, [])

    const addNote = useCallback(async () => {
        if (!sessionId || !form.note.trim()) return
        setSubmitting(true)
        setError('')
        try {
            const res = await api.post(`/api/sessions/${sessionId}/notes`, {
                chapter: form.chapter.trim(),
                pages: form.pages.trim(),
                note: form.note.trim(),
            })
            setNotes((prev) => [...prev, res.data.data])
            setForm(emptyForm)
        } catch {
            setError('Failed to save note. Try again.')
        } finally {
            setSubmitting(false)
        }
    }, [sessionId, form])

    const deleteNote = useCallback(async (noteId: string) => {
        try {
            await api.delete(`/api/notes/${noteId}`)
            setNotes((prev) => prev.filter((n) => n.id !== noteId))
        } catch {
            setError('Failed to delete note')
        }
    }, [])

    const clearNotes = useCallback(() => {
        setNotes([])
        setForm(emptyForm)
        setError('')
    }, [])

    return {
        notes,
        form,
        submitting,
        error,
        updateForm,
        addNote,
        deleteNote,
        clearNotes,
    }
}
