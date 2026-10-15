import { useState, useEffect } from 'react'

interface Entry {
  rank: number
  device_id: string
  total_rated: number
  total_points: number
  is_me: boolean
}

export function Leaderboard({ deviceID }: { deviceID: string }) {
  const [entries, setEntries] = useState<Entry[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetch('/api/leaderboard', { headers: { 'X-Device-ID': deviceID } })
      .then(r => r.json())
      .then(d => { setEntries(d ?? []); setLoading(false) })
      .catch(() => setLoading(false))
  }, [deviceID])

  if (loading) return <div className="loading"><div className="spinner" /><span>Loading…</span></div>
  if (entries.length === 0) return <div className="empty-state">No scores yet.<br />Be the first to rate today's clips!</div>

  const medals = ['🥇', '🥈', '🥉']

  return (
    <div>
      <div className="card-label" style={{ marginBottom: 12 }}>Top raters</div>
      {entries.map(e => (
        <div key={e.rank} className={`lb-row ${e.is_me ? 'is-me' : ''}`}>
          <div className={`lb-rank ${e.rank <= 3 ? 'top3' : ''}`}>
            {e.rank <= 3 ? medals[e.rank - 1] : `#${e.rank}`}
          </div>
          <div className="lb-device">{e.device_id}{e.is_me ? ' (you)' : ''}</div>
          <div style={{ textAlign: 'right' }}>
            <div className="lb-points">{e.total_points} pts</div>
            <div className="lb-time">{e.total_rated} rated</div>
          </div>
        </div>
      ))}
    </div>
  )
}
