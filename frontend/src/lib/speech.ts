export interface VoiceOption {
  code: string
  label: string
  flag: string
}

export const VOICE_OPTIONS: VoiceOption[] = [
  { code: 'en-US', label: 'English (US)',    flag: '🇺🇸' },
  { code: 'en-GB', label: 'English (UK)',    flag: '🇬🇧' },
  { code: 'en-IN', label: 'English (India)', flag: '🇮🇳' },
  { code: 'hi-IN', label: 'Hindi',           flag: '🇮🇳' },
  { code: 'es-ES', label: 'Spanish',         flag: '🇪🇸' },
  { code: 'fr-FR', label: 'French',          flag: '🇫🇷' },
  { code: 'ar-SA', label: 'Arabic',          flag: '🇸🇦' },
]

// Returns all voices available on this device for a given lang code.
// Matches exact (en-US) then prefix (en-*).
export function getVoicesForLang(lang: string): SpeechSynthesisVoice[] {
  if (!isSpeechSupported()) return []
  const all = window.speechSynthesis.getVoices()
  const exact  = all.filter(v => v.lang === lang)
  const prefix = all.filter(v => v.lang.startsWith(lang.split('-')[0]) && v.lang !== lang)
  return [...exact, ...prefix]
}

export function getAllVoices(): SpeechSynthesisVoice[] {
  if (!isSpeechSupported()) return []
  return window.speechSynthesis.getVoices()
}

export function speak(
  text: string,
  lang: string,
  voiceName: string | null,
  onEnd?: () => void,
): void {
  if (!window.speechSynthesis) return
  window.speechSynthesis.cancel()

  const utter = new SpeechSynthesisUtterance(text)
  utter.lang = lang
  utter.rate = 0.88
  utter.pitch = 1

  const voices = window.speechSynthesis.getVoices()
  if (voiceName) {
    const match = voices.find(v => v.name === voiceName)
    if (match) utter.voice = match
  } else {
    const match = voices.find(v => v.lang === lang)
      || voices.find(v => v.lang.startsWith(lang.split('-')[0]))
    if (match) utter.voice = match
  }

  if (onEnd) utter.onend = onEnd
  window.speechSynthesis.speak(utter)
}

export function stopSpeaking(): void {
  window.speechSynthesis?.cancel()
}

export function isSpeechSupported(): boolean {
  return typeof window !== 'undefined' && 'speechSynthesis' in window
}
