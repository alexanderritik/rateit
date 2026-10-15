import { useState, useRef } from 'react'
import { speak, stopSpeaking, VOICE_OPTIONS } from '../lib/speech'

interface Clip {
  id: number
  clip_index: number
  text: string
  language: string
  score: number | null
  listened: boolean
}

interface Props {
  clip: Clip
  deviceID: string
  overrideLang: string | null   // null = use clip's own language
  onRated: (clipID: number, score: number, listened: boolean) => void
}

const LABELS: Record<number, { short: string; color: string }> = {
  1: { short: 'Very unnatural', color: '#ef4444' },
  2: { short: 'Mostly unnatural', color: '#f97316' },
  3: { short: 'OK',             color: '#eab308' },
  4: { short: 'Mostly natural', color: '#84cc16' },
  5: { short: 'Very natural',   color: '#22c55e' },
}

export function RatingCard({ clip, deviceID, overrideLang, onRated }: Props) {
  const [score, setScore] = useState<number | null>(clip.score)
  const [listened, setListened] = useState(clip.listened)
  const [speaking, setSpeaking] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const [submitted, setSubmitted] = useState(clip.score !== null)
  const listenedRef = useRef(clip.listened)

  function handlePlay() {
    if (speaking) {
      stopSpeaking()
      setSpeaking(false)
      return
    }
    setSpeaking(true)
    if (!listenedRef.current) {
      listenedRef.current = true
      setListened(true)
    }
    const voiceName = localStorage.getItem('rateit_voice') || null
    speak(clip.text, overrideLang ?? clip.language, voiceName, () => setSpeaking(false))
  }

  async function handleSubmit(s: number) {
    if (submitting) return
    setScore(s)
    setSubmitting(true)
    try {
      await fetch(`/api/clips/${clip.id}/rate`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json', 'X-Device-ID': deviceID },
        body: JSON.stringify({ score: s, listened: listenedRef.current }),
      })
      setSubmitted(true)
      onRated(clip.id, s, listenedRef.current)
    } catch {
      setSubmitting(false)
    }
  }

  const label = score ? LABELS[score] : null
  const activeLang = overrideLang ?? clip.language
  const voiceEntry = VOICE_OPTIONS.find(v => v.code === activeLang)
  const voiceLabel = voiceEntry ? `${voiceEntry.flag} ${voiceEntry.label}` : activeLang

  return (
    <div className={`rating-card ${submitted ? 'done' : ''}`}>
      {/* Clip index + language pill */}
      <div className="clip-meta">
        <span className="clip-index">#{clip.clip_index}</span>
        <span className="lang-pill">{voiceLabel}</span>
        {submitted && <span className="done-check">✓ Rated</span>}
      </div>

      {/* Text */}
      <p className="clip-text">{clip.text}</p>

      {/* Play button */}
      <button
        className={`play-btn ${speaking ? 'playing' : ''}`}
        onClick={handlePlay}
        disabled={submitted && !speaking}
      >
        {speaking ? (
          <><span className="play-icon">⏹</span> Stop</>
        ) : (
          <><span className="play-icon">▶</span> {listened ? 'Play again' : 'Play'}</>
        )}
      </button>

      {!submitted && (
        <>
          <div className="rating-label-row">
            <span className="rating-prompt">How natural does this sound?</span>
            {label && <span className="rating-label-val" style={{ color: label.color }}>{label.short}</span>}
          </div>

          {/* Star buttons */}
          <div className="stars">
            {[1, 2, 3, 4, 5].map(n => (
              <button
                key={n}
                className={`star-btn ${score === n ? 'selected' : ''} ${score && score >= n ? 'filled' : ''}`}
                onClick={() => setScore(n)}
                disabled={submitting}
                title={LABELS[n].short}
              >
                {n}
              </button>
            ))}
          </div>

          <button
            className="submit-rating-btn"
            onClick={() => score && handleSubmit(score)}
            disabled={score === null || submitting || !listened}
            title={!listened ? 'Listen to the clip first' : ''}
          >
            {submitting ? 'Saving…' : !listened ? '🔊 Listen first' : 'Submit rating'}
          </button>
        </>
      )}

      {submitted && score && (
        <div className="submitted-score" style={{ color: LABELS[score].color }}>
          Your rating: {'★'.repeat(score)}{'☆'.repeat(5 - score)} — {LABELS[score].short}
        </div>
      )}
    </div>
  )
}
