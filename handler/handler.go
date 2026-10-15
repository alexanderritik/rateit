package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/ritikmayank/catch-it/model"
	"github.com/ritikmayank/catch-it/repository"
)

const clipsPerDay = 5
const pointsPerRating = 20
const bonusAllDone = 50

type Handler struct {
	clips   *repository.ClipRepository
	ratings *repository.RatingRepository
	streaks *repository.StreakRepository
	badges  *repository.BadgeRepository
}

func New(
	clips *repository.ClipRepository,
	ratings *repository.RatingRepository,
	streaks *repository.StreakRepository,
	badges *repository.BadgeRepository,
) *Handler {
	return &Handler{clips: clips, ratings: ratings, streaks: streaks, badges: badges}
}

func json200(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

func jsonErr(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func deviceID(r *http.Request) string {
	id := strings.TrimSpace(r.Header.Get("X-Device-ID"))
	if id == "" {
		return "anonymous"
	}
	return id
}

// GET /api/today  — returns today's clips merged with the user's existing ratings
func (h *Handler) Today(w http.ResponseWriter, r *http.Request) {
	today := todayUTC()
	clips, err := h.clips.ListByDate(r.Context(), today)
	if err != nil || len(clips) == 0 {
		jsonErr(w, http.StatusNotFound, "no clips today")
		return
	}

	devID := deviceID(r)
	clipIDs := make([]int, len(clips))
	for i, c := range clips {
		clipIDs[i] = c.ID
	}

	ratingMap, _ := h.ratings.RatedClipIDs(r.Context(), devID, clipIDs)
	rated := 0

	result := make([]model.ClipWithRating, len(clips))
	for i, c := range clips {
		cwr := model.ClipWithRating{
			ID:        c.ID,
			ClipIndex: c.ClipIndex,
			Text:      c.Text,
			Language:  c.Language,
		}
		if rt, ok := ratingMap[c.ID]; ok {
			s := rt.Score
			cwr.Score = &s
			cwr.Listened = rt.Listened
			rated++
		}
		result[i] = cwr
	}

	json200(w, model.DayProgress{
		Date:    today.Format("2006-01-02"),
		Clips:   result,
		Total:   len(clips),
		Rated:   rated,
		DayDone: rated >= len(clips),
	})
}

// POST /api/clips/{id}/rate
func (h *Handler) Rate(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	clipID := parseID(idStr)
	if clipID == 0 {
		jsonErr(w, http.StatusBadRequest, "invalid clip id")
		return
	}

	var req model.RateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonErr(w, http.StatusBadRequest, "invalid body")
		return
	}
	if req.Score < 1 || req.Score > 5 {
		jsonErr(w, http.StatusBadRequest, "score must be 1–5")
		return
	}

	// Verify clip belongs to today
	clip, err := h.clips.GetByID(r.Context(), clipID)
	if err != nil {
		jsonErr(w, http.StatusNotFound, "clip not found")
		return
	}
	today := todayUTC()
	if !sameDay(clip.PlayDate, today) {
		jsonErr(w, http.StatusBadRequest, "clip is not from today")
		return
	}

	devID := deviceID(r)
	if _, err := h.ratings.Upsert(r.Context(), devID, clipID, req.Score, req.Listened); err != nil {
		log.Error().Err(err).Msg("upsert rating")
		jsonErr(w, http.StatusInternalServerError, "could not save rating")
		return
	}

	// Check day completion
	allClips, _ := h.clips.ListByDate(r.Context(), today)
	allIDs := make([]int, len(allClips))
	for i, c := range allClips {
		allIDs[i] = c.ID
	}
	ratedCount, _ := h.ratings.CountRatedToday(r.Context(), devID, allIDs)
	dayDone := ratedCount >= len(allClips)

	// Points: 20 per rating, +50 bonus when day is done for the first time
	points := pointsPerRating
	if dayDone && ratedCount == len(allClips) {
		points += bonusAllDone
	}

	// Update streak only when the full day is done
	currentStreak := 0
	totalPoints := 0
	if dayDone {
		currentStreak, _ = h.streaks.RecordPlay(r.Context(), devID, today, points)
	} else {
		// Still record partial points
		h.streaks.AddPoints(r.Context(), devID, pointsPerRating)
	}
	if s, err := h.streaks.Get(r.Context(), devID); err == nil {
		totalPoints = s.TotalPoints
		currentStreak = s.CurrentStreak
	}

	newBadges := h.awardBadges(r, devID, dayDone, allIDs, totalPoints)

	json200(w, model.RateResponse{
		ClipID:    clipID,
		Score:     req.Score,
		DayDone:   dayDone,
		Points:    points,
		Streak:    currentStreak,
		Level:     model.ComputeLevel(totalPoints),
		NewBadges: newBadges,
	})
}

func (h *Handler) awardBadges(r *http.Request, devID string, dayDone bool, allIDs []int, totalPoints int) []model.BadgeInfo {
	var newBadges []model.BadgeInfo

	tryAward := func(id string) {
		def, ok := model.BadgeCatalog[id]
		if !ok {
			return
		}
		awarded, err := h.badges.Award(r.Context(), devID, id)
		if err != nil || !awarded {
			return
		}
		newBadges = append(newBadges, model.BadgeInfo{
			ID: def.ID, Name: def.Name, Icon: def.Icon, Desc: def.Desc,
			EarnedAt: time.Now().UTC(),
		})
	}

	total, _ := h.ratings.CountTotal(r.Context(), devID)
	if total >= 1 {
		tryAward("first_rating")
	}
	if total >= 50 {
		tryAward("sharp_rater")
	}
	if total >= 200 {
		tryAward("voice_expert")
	}

	if dayDone {
		tryAward("day_one")

		allListened, _ := h.ratings.AllListenedToday(r.Context(), devID, allIDs)
		if allListened {
			tryAward("careful_ear")
		}

		if s, err := h.streaks.Get(r.Context(), devID); err == nil {
			if s.CurrentStreak >= 3 {
				tryAward("streak_3")
			}
			if s.CurrentStreak >= 7 {
				tryAward("streak_7")
			}
			if s.CurrentStreak >= 30 {
				tryAward("streak_30")
			}
		}
	}
	return newBadges
}

// GET /api/leaderboard
func (h *Handler) Leaderboard(w http.ResponseWriter, r *http.Request) {
	lbRows, err := h.ratings.Leaderboard(r.Context(), 10)
	if err != nil {
		log.Error().Err(err).Msg("leaderboard")
		jsonErr(w, http.StatusInternalServerError, "unavailable")
		return
	}

	devID := deviceID(r)
	entries := make([]model.LeaderboardEntry, len(lbRows))
	for i, row := range lbRows {
		entries[i] = model.LeaderboardEntry{
			Rank:        i + 1,
			DeviceID:    mask(row.DeviceID),
			TotalRated:  row.TotalRated,
			TotalPoints: row.TotalPoints,
			IsMe:        row.DeviceID == devID,
		}
	}
	json200(w, entries)
}

// GET /api/me
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	devID := deviceID(r)
	s, err := h.streaks.Get(r.Context(), devID)
	if err != nil {
		json200(w, model.Streak{DeviceID: devID, Level: model.ComputeLevel(0), Badges: []model.BadgeInfo{}})
		return
	}
	earned, _ := h.badges.ListByDevice(r.Context(), devID)
	if earned == nil {
		earned = []model.BadgeInfo{}
	}
	s.Level = model.ComputeLevel(s.TotalPoints)
	s.Badges = earned
	json200(w, s)
}

func Health(w http.ResponseWriter, r *http.Request) {
	json200(w, map[string]string{"status": "ok"})
}

// ── helpers ───────────────────────────────────────────────────────────────────

func todayUTC() time.Time {
	return time.Now().UTC().Truncate(24 * time.Hour)
}

func sameDay(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}

func parseID(s string) int {
	var id int
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0
		}
		id = id*10 + int(c-'0')
	}
	return id
}

func mask(id string) string {
	if len(id) <= 4 {
		return id
	}
	return id[:4] + "****"
}
