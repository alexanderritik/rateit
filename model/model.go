package model

import "time"

// ── Clips & Ratings ──────────────────────────────────────────────────────────

type Clip struct {
	ID        int       `json:"id"`
	PlayDate  time.Time `json:"-"`
	ClipIndex int       `json:"clip_index"`
	Text      string    `json:"text"`
	Language  string    `json:"language"`
	CreatedAt time.Time `json:"-"`
}

// ClipWithRating is what the client receives — includes the user's existing score if any
type ClipWithRating struct {
	ID        int    `json:"id"`
	ClipIndex int    `json:"clip_index"`
	Text      string `json:"text"`
	Language  string `json:"language"`
	Score     *int   `json:"score"` // nil = not yet rated
	Listened  bool   `json:"listened"`
}

type RateRequest struct {
	Score    int  `json:"score"`    // 1–5
	Listened bool `json:"listened"` // did the user play the clip?
}

type RateResponse struct {
	ClipID    int  `json:"clip_id"`
	Score     int  `json:"score"`
	DayDone   bool `json:"day_done"`   // all 5 clips rated today
	Points    int  `json:"points"`
	Streak    int  `json:"streak"`
	Level     Level       `json:"level"`
	NewBadges []BadgeInfo `json:"new_badges"`
}

type DayProgress struct {
	Date     string           `json:"date"`
	Clips    []ClipWithRating `json:"clips"`
	Total    int              `json:"total"`    // clips today
	Rated    int              `json:"rated"`    // how many user has rated
	DayDone  bool             `json:"day_done"`
}

type Rating struct {
	ID       int       `json:"id"`
	DeviceID string    `json:"device_id"`
	ClipID   int       `json:"clip_id"`
	Score    int       `json:"score"`
	Listened bool      `json:"listened"`
	RatedAt  time.Time `json:"rated_at"`
}

type LeaderboardEntry struct {
	Rank        int    `json:"rank"`
	DeviceID    string `json:"device_id"`
	TotalRated  int    `json:"total_rated"`
	TotalPoints int    `json:"total_points"`
	IsMe        bool   `json:"is_me"`
}

// ── Streaks ───────────────────────────────────────────────────────────────────

type Streak struct {
	DeviceID       string      `json:"device_id"`
	CurrentStreak  int         `json:"current_streak"`
	LongestStreak  int         `json:"longest_streak"`
	LastPlayedDate *string     `json:"last_played_date"`
	TotalPoints    int         `json:"total_points"`
	UpdatedAt      time.Time   `json:"updated_at"`
	Level          Level       `json:"level"`
	Badges         []BadgeInfo `json:"badges"`
}

// ── Levels ────────────────────────────────────────────────────────────────────

type Level struct {
	Name       string `json:"name"`
	Icon       string `json:"icon"`
	TotalPts   int    `json:"total_pts"`
	NextLevelAt int   `json:"next_level_at"` // -1 = max
}

func ComputeLevel(totalPoints int) Level {
	switch {
	case totalPoints >= 8000:
		return Level{"Master", "🏆", totalPoints, -1}
	case totalPoints >= 3000:
		return Level{"Expert", "🎓", totalPoints, 8000}
	case totalPoints >= 1000:
		return Level{"Analyst", "🔬", totalPoints, 3000}
	case totalPoints >= 300:
		return Level{"Listener", "👂", totalPoints, 1000}
	default:
		return Level{"Rookie", "🌱", totalPoints, 300}
	}
}

// ── Badges ────────────────────────────────────────────────────────────────────

type BadgeInfo struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Icon     string    `json:"icon"`
	Desc     string    `json:"desc"`
	EarnedAt time.Time `json:"earned_at"`
}

type BadgeDef struct {
	ID   string
	Name string
	Icon string
	Desc string
}

var BadgeCatalog = map[string]BadgeDef{
	"first_rating": {"first_rating", "First Rating",  "⭐", "Rated your first clip"},
	"day_one":      {"day_one",      "Day One",        "🌅", "Completed your first full day"},
	"streak_3":     {"streak_3",     "On a Roll",      "🔥", "3 days in a row"},
	"streak_7":     {"streak_7",     "Week Warrior",   "💪", "7 days in a row"},
	"streak_30":    {"streak_30",    "Unstoppable",    "🌟", "30 days in a row"},
	"careful_ear":  {"careful_ear",  "Careful Ear",    "👂", "Listened before rating every clip in a day"},
	"sharp_rater":  {"sharp_rater",  "Sharp Rater",    "🎯", "Rated 50 clips total"},
	"voice_expert": {"voice_expert", "Voice Expert",   "🎤", "Rated 200 clips total"},
}
