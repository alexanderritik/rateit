export interface Theme {
  id: string
  label: string
  preview: string   // swatch color
  vars: Record<string, string>
}

export const THEMES: Theme[] = [
  {
    id: 'purple',
    label: 'Violet',
    preview: '#6c63ff',
    vars: {
      '--bg':           '#0f0f1a',
      '--surface':      '#1a1a2e',
      '--border':       '#2a2a45',
      '--accent':       '#6c63ff',
      '--accent-light': '#8b85ff',
    },
  },
  {
    id: 'ocean',
    label: 'Ocean',
    preview: '#0ea5e9',
    vars: {
      '--bg':           '#020d1a',
      '--surface':      '#0a1929',
      '--border':       '#0f2d45',
      '--accent':       '#0ea5e9',
      '--accent-light': '#38bdf8',
    },
  },
  {
    id: 'forest',
    label: 'Forest',
    preview: '#22c55e',
    vars: {
      '--bg':           '#030f0a',
      '--surface':      '#071a0f',
      '--border':       '#0e3320',
      '--accent':       '#22c55e',
      '--accent-light': '#4ade80',
    },
  },
  {
    id: 'ember',
    label: 'Ember',
    preview: '#f97316',
    vars: {
      '--bg':           '#140a03',
      '--surface':      '#1f1007',
      '--border':       '#3d200a',
      '--accent':       '#f97316',
      '--accent-light': '#fb923c',
    },
  },
  {
    id: 'light',
    label: 'Light',
    preview: '#f8f8ff',
    vars: {
      '--bg':           '#f0f0f8',
      '--surface':      '#ffffff',
      '--border':       '#d4d4e8',
      '--accent':       '#6c63ff',
      '--accent-light': '#4f46e5',
      '--text':         '#1a1a2e',
      '--text-muted':   '#6060a0',
    },
  },
]

const THEME_KEY = 'rateit_theme'

export function getSavedThemeID(): string {
  return localStorage.getItem(THEME_KEY) || 'purple'
}

export function applyTheme(id: string): void {
  const theme = THEMES.find(t => t.id === id) || THEMES[0]
  const root = document.documentElement
  Object.entries(theme.vars).forEach(([k, v]) => root.style.setProperty(k, v))
  // Reset overrides not in this theme (e.g. --text for dark themes)
  if (id !== 'light') {
    root.style.setProperty('--text', '#e8e8f0')
    root.style.setProperty('--text-muted', '#9090a8')
  }
  localStorage.setItem(THEME_KEY, id)
}
