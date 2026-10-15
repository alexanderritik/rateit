Ref: KL-7Q4

# RateIt — Daily TTS Quality Rater

A daily voice evaluation app. Each day you get 5 short sentences synthesised by the device's text-to-speech engine. Listen to each one and rate how natural it sounds on a 1–5 scale. Complete all 5 to keep your streak alive and earn points.

Built with Go + React for the Kenpath Labs Full Stack Engineer assessment.

---

## The idea

Kenpath builds Svara TTS (82 languages) and Lokah, a human data platform that collects speech quality ratings. RateIt is what that task looks like when it's gamified for daily use.

Each session takes under two minutes:
1. Open the app — 5 clips are ready
2. Press **Play** on each clip — the browser synthesises it aloud
3. Rate it: *1 Very unnatural → 5 Very natural*
4. Submit all 5 → streak updates, points credited, badges awarded

Choose your playback language from a strip at the top of the screen (Auto, or any of 7 languages). Switch TTS voice model and app theme from the Settings tab.

No account needed. Progress is tied to a device ID stored in localStorage.

---

## Rewards

| What | How |
|------|-----|
| Points | 20 pts per clip rated + 50 bonus for completing all 5 in a day |
| Streak | Increments when you complete a full day of ratings |
| Levels | Rookie → Listener → Analyst → Expert → Master (based on total points) |
| Badges | 8 badges: first rating, day completion, streaks, careful listening, volume milestones |

---

## Running locally

**Prerequisites:** Go 1.21+, Node.js 18+, Docker

### 1. Clone

```bash
git clone https://github.com/alexanderritik/rateit
cd rateit
```

### 2. Start the database

```bash
docker compose up -d
```

Starts PostgreSQL on port 5432.

### 3. Configure environment

```bash
cp .env.example .env
```

Defaults work out of the box with Docker Compose.

### 4. Build the frontend

```bash
cd frontend
npm install
npm run build
cd ..
```

Builds the React app into `web/dist/` — the Go binary embeds this at compile time.

### 5. Run the server

```bash
go run .
```

Server starts at [http://localhost:8080](http://localhost:8080).

Migrations run automatically on startup. 150 clips (5 per day × 30 days) are seeded at first run.

---

## Development mode (hot reload)

```bash
# Terminal 1 — API
go run .

# Terminal 2 — frontend with HMR
cd frontend && npm run dev
```

Vite proxies `/api` to `:8080`. Frontend runs on `:5173`.

---

## Deployment

Single binary with embedded frontend — only needs a PostgreSQL database.

**Railway (recommended):**
1. Create a new project → add a PostgreSQL plugin (sets `DATABASE_URL` automatically)
2. Connect the GitHub repo and deploy
3. Railway sets `PORT` automatically — no extra config needed

**Docker:**
```bash
docker build -t rateit .
docker run -e DATABASE_URL=<your-db-url> -p 8080:8080 rateit
```

---

## API

All endpoints accept `X-Device-ID: <uuid>` to identify the player (no auth required).

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/health` | Health check |
| GET | `/api/today` | Today's 5 clips with user's existing ratings |
| POST | `/api/clips/:id/rate` | Submit rating `{score: 1–5, listened: bool}` |
| GET | `/api/leaderboard` | Top 10 by total points |
| GET | `/api/me` | Current device's streak, level, badges, total points |

---

## Project structure

```
rateit/
├── main.go               # Server startup, routing, graceful shutdown
├── migrate.go            # Auto-migration runner (golang-migrate)
├── config/               # Environment variable loading
├── db/                   # PostgreSQL connection pool (pgx/v5)
├── migrations/           # SQL schema (clips, ratings, streaks, badges)
├── model/                # Shared types + level/badge catalog
├── repository/           # DB access — clips, ratings, streaks, badges
├── handler/              # HTTP handlers + badge award logic
├── seed/                 # 30-day clip pool (mixed languages)
├── web/                  # Embeds React build into Go binary at compile time
└── frontend/             # React + Vite + TypeScript
    └── src/
        ├── App.tsx
        ├── lib/
        │   ├── device.ts   # Device ID (localStorage)
        │   ├── speech.ts   # Web Speech API wrapper
        │   ├── badges.ts   # Badge catalog
        │   └── theme.ts    # Theme definitions + applicator
        └── components/
            ├── Today.tsx       # Daily clip list + language selector
            ├── RatingCard.tsx  # Single clip: play + rate + submit
            ├── Leaderboard.tsx # Top 10 raters
            ├── Profile.tsx     # Streak, level progress, badge grid
            └── Settings.tsx    # Voice model picker + theme switcher
```

---

## What I built and why

Kenpath's products sit at the intersection of speech synthesis and human data collection. Their Lokah platform pays workers to complete annotation tasks — including rating how natural synthesised speech sounds. That's a real workflow called a Mean Opinion Score (MOS) test, and it's how every major TTS lab benchmarks quality.

RateIt is that workflow turned into a daily habit app. Users listen to clips, rate them 1–5, earn points and streaks. The data they generate is exactly the kind Kenpath sells to AI companies via Lokah. The gamification layer is what makes users come back without being paid per session.

The app runs as a single Go binary with the React frontend embedded at compile time. No CDN, no separate static host — one deployment on Railway with a Postgres plugin and it's live. Appropriate for the target audience on mobile data: the JS bundle is ~49 KB gzipped.

**Choices made:**
- Web Speech API for TTS — zero cost, zero latency, works offline, available on Android Chrome. The natural upgrade path is to swap in Svara TTS Turbo via API once a key is available.
- Device ID over auth — no login friction. The tradeoff is that leaderboard identities aren't portable across devices, which is acceptable for a v1.
- 5 clips per day — enough to feel substantial, short enough to finish on a commute.

**What I left out:**
- Svara TTS API integration — used browser TTS as a stand-in. The `speak()` function in `speech.ts` is the only place that needs to change.
- Push notifications — important for streak retention, straightforward to add with a service worker.
- Admin panel — clips are seeded from code. A simple authenticated endpoint to add clips per day is the obvious next step.
- Cross-device sync — would require real auth (Supabase or similar). The device ID approach handles the solo-user case well.

**AI tools used:**
- Claude Code (Anthropic) — used for code generation, architecture decisions, and debugging throughout the build.

---

## License

MIT
