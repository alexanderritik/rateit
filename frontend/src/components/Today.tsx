import { useState, useEffect } from 'react'
import { RatingCard, RateResp } from './RatingCard'
import { VOICE_OPTIONS } from '../lib/speech'

interface Clip {
  id: number
  clip_index: number
  text: string
  language: string
  score: number | null
  listened: boolean
}

interface DayProgress {
  date: string
  clips: Clip[]
  total: number
  rated: number
  day_done: boolean
}

const LANG_KEY = 'rateit_lang'

const SELECTOR = [
  { code: 'auto', label: 'Auto', flag: '🔄' },
  ...VOICE_OPTIONS,
]

export function Today({ deviceID, onComplete }: { deviceID: string; onComplete: () => void }) {
  const [progress, setProgress] = useState<DayProgress | null>(null)
  const [loading, setLoading]   = useState(true)
  const [lastResp, setLastResp] = useState<RateResp | null>(null)
  const [lang, setLang]         = useState(() => localStorage.getItem(LANG_KEY) || 'auto')

  useEffect(() => {
    fetch('/api/today', { headers: { 'X-Device-ID': deviceID } })
      .then(r => r.json())
      .then(d => { setProgress(d); setLoading(false) })
      .catch(() => setLoading(false))
  }, [deviceID])

  function handleLangChange(code: string) {
    setLang(code)
    localStorage.setItem(LANG_KEY, code)
  }

  // Called by RatingCard with the full server response — no second fetch
  function handleRated(resp: RateResp) {
    setLastResp(resp)
    if (resp.day_done) onComplete()

    setProgress(prev => {
      if (!prev) return prev
      const updated = prev.clips.map(c => c.id === resp.clip_id ? { ...c, score: resp.score } : c)
      const rated   = updated.filter(c => c.score !== null).length
      return { ...prev, clips: updated, rated, day_done: rated >= prev.total }
    })
  }

  if (loading)   return <div className="loading"><div className="spinner" /><span>Loading today's clips…</span></div>
  if (!progress) return <div className="empty-state">No clips today. Come back soon!</div>

  const pct        = Math.round((progress.rated / progress.total) * 100)
  const activeLang = SELECTOR.find(s => s.code === lang) || SELECTOR[0]

  return (
    <div>
      {/* Progress bar */}
      <div className="day-progress-wrap">
        <div className="day-progress-label">
          <span>{progress.rated} of {progress.total} rated</span>
          {progress.day_done && <span className="done-badge">Day complete ✓</span>}
        </div>
        <div className="day-progress-bar">
          <div className="day-progress-fill" style={{ width: `${pct}%` }} />
        </div>
      </div>

      {/* Language selector */}
      <div className="lang-selector-wrap">
        <span className="lang-selector-label">Listen in</span>
        <div className="lang-selector-strip">
          {SELECTOR.map(s => (
            <button
              key={s.code}
              className={`lang-chip ${lang === s.code ? 'active' : ''}`}
              onClick={() => handleLangChange(s.code)}
              title={s.label}
            >
              {s.flag} <span>{s.code === 'auto' ? 'Auto' : s.code.split('-')[0].toUpperCase()}</span>
            </button>
          ))}
        </div>
      </div>

      {/* New badges */}
      {lastResp && lastResp.new_badges.length > 0 && (
        <div className="card" style={{ marginBottom: 16 }}>
          <div className="card-label" style={{ marginBottom: 10 }}>New badges!</div>
          <div className="badge-grid">
            {lastResp.new_badges.map(b => (
              <div key={b.id} className="badge-item new">
                <div className="badge-icon">{b.icon}</div>
                <div className="badge-name">{b.name}</div>
                <div className="badge-desc">{b.desc}</div>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Day done banner */}
      {progress.day_done && lastResp && (
        <div className="day-done-banner">
          <div className="day-done-icon">🎉</div>
          <div>
            <div className="day-done-title">All done for today!</div>
            <div className="day-done-sub">
              {lastResp.level.icon} {lastResp.level.name} · 🔥 {lastResp.streak} day streak
            </div>
          </div>
        </div>
      )}

      {/* Clip cards */}
      <div className="clips-list">
        {progress.clips.map(clip => (
          <RatingCard
            key={clip.id}
            clip={clip}
            deviceID={deviceID}
            overrideLang={lang === 'auto' ? null : lang}
            onRated={handleRated}
          />
        ))}
      </div>

      <p className="voice-info-note">
        {lang === 'auto'
          ? 'Each clip plays in its own language.'
          : `All clips play in ${activeLang.flag} ${activeLang.label}.`}
      </p>
    </div>
  )
}
