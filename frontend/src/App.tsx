import { useState, useEffect } from 'react'
import { Today } from './components/Today'
import { Leaderboard } from './components/Leaderboard'
import { Profile } from './components/Profile'
import { Settings } from './components/Settings'
import { getOrCreateDeviceID } from './lib/device'

type Tab = 'today' | 'leaderboard' | 'profile' | 'settings'

interface MeData {
  current_streak: number
  total_points: number
  level: { name: string; icon: string }
}

export default function App() {
  const [tab, setTab] = useState<Tab>('today')
  const [me, setMe] = useState<MeData | null>(null)
  const deviceID = getOrCreateDeviceID()

  const refreshMe = () => {
    fetch('/api/me', { headers: { 'X-Device-ID': deviceID } })
      .then(r => r.json())
      .then(d => setMe(d))
      .catch(() => {})
  }

  useEffect(() => { refreshMe() }, [deviceID])

  return (
    <div className="app">
      <header className="header">
        <div className="header-logo">Rate<span>It</span></div>
        <div className="header-right">
          {me?.level && (
            <div className="level-badge">{me.level.icon} {me.level.name}</div>
          )}
          {me && me.current_streak > 0 && (
            <div className="streak-badge">🔥 {me.current_streak}</div>
          )}
        </div>
      </header>

      <div className="tabs">
        <button className={`tab ${tab === 'today'       ? 'active' : ''}`} onClick={() => setTab('today')}>Today</button>
        <button className={`tab ${tab === 'leaderboard' ? 'active' : ''}`} onClick={() => setTab('leaderboard')}>Board</button>
        <button className={`tab ${tab === 'profile'     ? 'active' : ''}`} onClick={() => setTab('profile')}>Me</button>
        <button className={`tab ${tab === 'settings'    ? 'active' : ''}`} onClick={() => setTab('settings')}>⚙</button>
      </div>

      {tab === 'today'       && <Today       deviceID={deviceID} onComplete={refreshMe} />}
      {tab === 'leaderboard' && <Leaderboard deviceID={deviceID} />}
      {tab === 'profile'     && <Profile     deviceID={deviceID} />}
      {tab === 'settings'    && <Settings />}
    </div>
  )
}
