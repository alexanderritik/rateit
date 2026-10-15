import { useState, useEffect } from 'react'
import { ALL_BADGES } from '../lib/badges'

interface BadgeInfo { id: string; name: string; icon: string; desc: string; earned_at: string }
interface Level { name: string; icon: string; total_pts: number; next_level_at: number }
interface MeData {
  current_streak: number
  longest_streak: number
  total_points: number
  level: Level
  badges: BadgeInfo[]
}

const LEVELS = [
  { name: 'Rookie',       icon: '🌱', min: 0,    next: 300  },
  { name: 'Explorer',     icon: '🗺️', min: 300,  next: 1000 },
  { name: 'Investigator', icon: '🧐', min: 1000, next: 3000 },
  { name: 'Detective',    icon: '🔍', min: 3000, next: 8000 },
  { name: 'Master',       icon: '🏆', min: 8000, next: -1   },
]

export function Profile({ deviceID }: { deviceID: string }) {
  const [me, setMe] = useState<MeData | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetch('/api/me', { headers: { 'X-Device-ID': deviceID } })
      .then(r => r.json())
      .then(d => { setMe(d); setLoading(false) })
      .catch(() => setLoading(false))
  }, [deviceID])

  if (loading) return <div className="loading"><div className="spinner" /><span>Loading…</span></div>
  if (!me) return <div className="empty-state">Could not load profile.</div>

  const earnedIds = new Set((me.badges || []).map(b => b.id))
  const lvl = me.level
  const progress = lvl.next_level_at === -1 ? 100
    : Math.round(((me.total_points - (LEVELS.find(l => l.name === lvl.name)?.min ?? 0)) /
        (lvl.next_level_at - (LEVELS.find(l => l.name === lvl.name)?.min ?? 0))) * 100)

  return (
    <div>
      {/* Level card */}
      <div className="card">
        <div className="profile-level-row">
          <div className="profile-level-icon">{lvl.icon}</div>
          <div>
            <div className="profile-level-name">{lvl.name}</div>
            <div className="profile-level-pts">{me.total_points} pts total</div>
          </div>
        </div>
        {lvl.next_level_at !== -1 && (
          <>
            <div className="level-progress-bar">
              <div className="level-progress-fill" style={{ width: `${Math.min(progress, 100)}%` }} />
            </div>
            <div className="level-progress-label">
              {me.total_points} / {lvl.next_level_at} pts → {LEVELS.find(l => l.min === lvl.next_level_at)?.icon} {LEVELS.find(l => l.min === lvl.next_level_at)?.name}
            </div>
          </>
        )}
        {lvl.next_level_at === -1 && (
          <div className="level-progress-label" style={{ color: 'var(--warning)' }}>Max level reached 🏆</div>
        )}
      </div>

      {/* Stats */}
      <div className="stats-row">
        <div className="stat-box">
          <div className="stat-number">{me.current_streak}</div>
          <div className="stat-label">Current streak</div>
        </div>
        <div className="stat-box">
          <div className="stat-number">{me.longest_streak}</div>
          <div className="stat-label">Best streak</div>
        </div>
        <div className="stat-box">
          <div className="stat-number">{me.total_points}</div>
          <div className="stat-label">Total points</div>
        </div>
      </div>

      {/* Badges */}
      <div className="card-label" style={{ marginBottom: 12 }}>Badges</div>
      <div className="badge-grid">
        {ALL_BADGES.map(b => {
          const earned = earnedIds.has(b.id)
          return (
            <div key={b.id} className={`badge-item ${earned ? 'earned' : 'locked'}`}>
              <div className="badge-icon">{earned ? b.icon : '🔒'}</div>
              <div className="badge-name">{b.name}</div>
              <div className="badge-desc">{b.desc}</div>
            </div>
          )
        })}
      </div>
    </div>
  )
}
