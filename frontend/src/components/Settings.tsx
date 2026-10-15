import { useState, useEffect } from 'react'
import { getVoicesForLang, isSpeechSupported, VOICE_OPTIONS } from '../lib/speech'
import { THEMES, applyTheme, getSavedThemeID } from '../lib/theme'

const VOICE_KEY  = 'rateit_voice'
const LANG_KEY   = 'rateit_lang'

export function Settings() {
  const [lang, setLang]           = useState(() => localStorage.getItem(LANG_KEY) || 'en-US')
  const [voiceName, setVoiceName] = useState(() => localStorage.getItem(VOICE_KEY) || '')
  const [voices, setVoices]       = useState<SpeechSynthesisVoice[]>([])
  const [themeID, setThemeID]     = useState(getSavedThemeID)
  const speechOK = isSpeechSupported()

  // Voices load async on some browsers
  useEffect(() => {
    function load() { setVoices(getVoicesForLang(lang)) }
    load()
    window.speechSynthesis?.addEventListener('voiceschanged', load)
    return () => window.speechSynthesis?.removeEventListener('voiceschanged', load)
  }, [lang])

  function handleLangChange(code: string) {
    setLang(code)
    setVoiceName('')
    localStorage.setItem(LANG_KEY, code)
    localStorage.removeItem(VOICE_KEY)
    setVoices(getVoicesForLang(code))
  }

  function handleVoiceChange(name: string) {
    setVoiceName(name)
    localStorage.setItem(VOICE_KEY, name)
  }

  function handleTheme(id: string) {
    setThemeID(id)
    applyTheme(id)
  }

  return (
    <div>
      {/* ── Theme ── */}
      <div className="settings-section">
        <div className="settings-title">App theme</div>
        <div className="theme-grid">
          {THEMES.map(t => (
            <button
              key={t.id}
              className={`theme-btn ${themeID === t.id ? 'active' : ''}`}
              onClick={() => handleTheme(t.id)}
            >
              <span className="theme-swatch" style={{ background: t.preview }} />
              <span className="theme-label">{t.label}</span>
              {themeID === t.id && <span className="theme-check">✓</span>}
            </button>
          ))}
        </div>
      </div>

      {/* ── Default playback language ── */}
      {speechOK && (
        <div className="settings-section">
          <div className="settings-title">Default playback language</div>
          <div className="settings-desc">Used when "Auto" is selected on the Today screen.</div>
          <div className="lang-selector-strip" style={{ flexWrap: 'wrap', marginTop: 10 }}>
            {VOICE_OPTIONS.map(v => (
              <button
                key={v.code}
                className={`lang-chip ${lang === v.code ? 'active' : ''}`}
                onClick={() => handleLangChange(v.code)}
              >
                {v.flag} {v.label}
              </button>
            ))}
          </div>
        </div>
      )}

      {/* ── Voice model ── */}
      {speechOK && (
        <div className="settings-section">
          <div className="settings-title">Voice model</div>
          <div className="settings-desc">
            {voices.length === 0
              ? 'No voices found for this language on your device.'
              : `${voices.length} voice${voices.length > 1 ? 's' : ''} available for ${lang}`}
          </div>
          {voices.length > 0 && (
            <div className="voice-list">
              <button
                className={`voice-row ${voiceName === '' ? 'active' : ''}`}
                onClick={() => handleVoiceChange('')}
              >
                <div>
                  <div className="voice-name">Default</div>
                  <div className="voice-sub">Let the browser choose</div>
                </div>
                {voiceName === '' && <span className="voice-check">✓</span>}
              </button>
              {voices.map(v => (
                <button
                  key={v.name}
                  className={`voice-row ${voiceName === v.name ? 'active' : ''}`}
                  onClick={() => handleVoiceChange(v.name)}
                >
                  <div>
                    <div className="voice-name">{v.name}</div>
                    <div className="voice-sub">{v.lang} {v.localService ? '· on-device' : '· online'}</div>
                  </div>
                  {voiceName === v.name && <span className="voice-check">✓</span>}
                </button>
              ))}
            </div>
          )}
        </div>
      )}

      {!speechOK && (
        <div className="empty-state">
          Speech synthesis is not supported in this browser.
        </div>
      )}
    </div>
  )
}
