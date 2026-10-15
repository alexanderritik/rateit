Ref: KL-7Q4

# Catch It — Daily Error Hunt

A daily micro-challenge app. You get three short sentences. One of them contains a deliberate mistake. Tap it, submit, and collect points. Come back tomorrow for a fresh challenge.

Built with Go + React for the Kenpath Labs Full Stack Engineer assessment.

---

## The idea

Every day, a short passage appears on screen. The sentences are split into three tappable blocks. Exactly one block contains a factual, logical, or numerical error. The faster you find it, the more points you earn.

Points:
- Correct answer: 100 base points
- Speed bonus: up to 100 additional points (reduces by 1 per 500ms)
- Wrong answer: 0 points

Returning daily builds a streak. The leaderboard resets each day.

No account needed. Your progress is stored against a device ID in localStorage.

---

## Running locally

**Prerequisites:** Go 1.21+, Node.js 18+, Docker

### 1. Clone and enter the project

```bash
git clone <repo-url>
cd catch-it
```

### 2. Start the database

```bash
docker compose up -d
```

This starts a PostgreSQL instance on port 5432.

### 3. Configure environment

```bash
cp .env.example .env
```

Defaults in `.env.example` work with the Docker Compose setup out of the box.

### 4. Build the frontend

```bash
cd frontend
npm install
npm run build
cd ..
```

This builds the React app and writes it to `web/dist/`, where the Go binary embeds it.

### 5. Run the server

```bash
go run .
```

Server starts on [http://localhost:8080](http://localhost:8080).

Migrations run automatically on startup. Challenges are seeded for the next 30 days.

---

## Development mode (hot reload)

Run the Go server and Vite dev server side by side:

```bash
# Terminal 1 — API server
go run .

# Terminal 2 — frontend dev server with proxy
cd frontend
npm run dev
```

The Vite dev server proxies `/api` requests to `:8080`, so both work together on `:5173`.

---

## Deployment

The app is a single Go binary with the React build embedded. It needs only a PostgreSQL database.

**Railway (recommended free tier):**
1. Create a new project on Railway
2. Add a PostgreSQL plugin — Railway provides `DATABASE_URL` automatically
3. Connect the repository and deploy
4. Set `PORT` env var if needed (Railway sets it automatically)

**Docker:**
```bash
docker build -t catch-it .
docker run -e DATABASE_URL=<your-db-url> -p 8080:8080 catch-it
```

---

## API

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/health` | Health check |
| GET | `/api/challenge/today` | Today's challenge (no answer) |
| POST | `/api/challenge/:id/submit` | Submit answer `{segment_id, time_ms}` |
| GET | `/api/leaderboard/today` | Top 10 correct plays today |
| GET | `/api/me` | Current device's streak and points |

Pass `X-Device-ID: <uuid>` in all requests to identify the player.

---

## Project structure

```
catch-it/
├── main.go               # Server startup, routing
├── migrate.go            # Migration runner
├── config/               # Environment config
├── db/                   # PostgreSQL connection pool
├── migrations/           # SQL schema files
├── model/                # Shared data types
├── repository/           # Database access (challenges, plays, streaks)
├── handler/              # HTTP handlers (challenge, leaderboard, me)
├── seed/                 # 30 pre-written daily challenges
├── web/                  # Embeds the React build into the binary
└── frontend/             # React + Vite app
    └── src/
        ├── App.tsx
        ├── lib/device.ts       # Device ID persistence
        └── components/
            ├── Challenge.tsx   # Main game screen
            └── Leaderboard.tsx # Daily top 10
```

---

## What I built and why

**The concept** came from bridging two things: the core mechanic of the production monitoring tool I'd been working on (run something → check if it passes or fails → record the result) and the Kenpath brief (short daily task + rewards).

"Catch It" is that mechanic applied to everyday knowledge: three sentences, one mistake, find it before time runs out. It works for ESL users because the sentences are short and the task is visual (tap a block) rather than text input. It loads fast because the Go binary serves the React build as static files — no CDN latency, no separate API server to cold-start.

**What I left out:**
- User accounts / auth — device ID covers 80% of the use case without any login friction. The 20% case (multiple devices, leaderboard identity) is a good v2 problem.
- Push notifications for streak reminders — important for retention but out of scope for this task.
- Admin panel for adding challenges — challenges are seeded from code for now. A simple CMS endpoint would be the next step.
- Badges — designed for but not implemented. The streak + leaderboard loop is enough to demonstrate the reward system.

**AI tools used:**
- Claude Code (Anthropic) for code generation and architecture decisions throughout the build.

---

## License

MIT
