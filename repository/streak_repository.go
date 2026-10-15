package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ritikmayank/catch-it/model"
)

type StreakRepository struct {
	db *pgxpool.Pool
}

func NewStreakRepository(db *pgxpool.Pool) *StreakRepository {
	return &StreakRepository{db: db}
}

func (r *StreakRepository) Get(ctx context.Context, deviceID string) (*model.Streak, error) {
	var s model.Streak
	err := r.db.QueryRow(ctx,
		`SELECT device_id, current_streak, longest_streak,
		        TO_CHAR(last_played_date, 'YYYY-MM-DD'), total_points, updated_at
		 FROM streaks WHERE device_id=$1`,
		deviceID,
	).Scan(&s.DeviceID, &s.CurrentStreak, &s.LongestStreak, &s.LastPlayedDate, &s.TotalPoints, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// AddPoints adds points to a device's total without touching streak or last_played_date.
func (r *StreakRepository) AddPoints(ctx context.Context, deviceID string, points int) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO streaks (device_id, total_points, updated_at)
		 VALUES ($1, $2, NOW())
		 ON CONFLICT (device_id) DO UPDATE SET
		   total_points = streaks.total_points + $2,
		   updated_at = NOW()`,
		deviceID, points,
	)
	return err
}

// RecordPlay updates the streak for a device after a correct play.
// Returns the new current streak.
func (r *StreakRepository) RecordPlay(ctx context.Context, deviceID string, today time.Time, pointsEarned int) (int, error) {
	todayStr := today.Format("2006-01-02")
	yesterdayStr := today.AddDate(0, 0, -1).Format("2006-01-02")

	var newStreak int
	err := r.db.QueryRow(ctx,
		`INSERT INTO streaks (device_id, current_streak, longest_streak, last_played_date, total_points, updated_at)
		 VALUES ($1, 1, 1, $2, $3, NOW())
		 ON CONFLICT (device_id) DO UPDATE SET
		   current_streak = CASE
		     WHEN streaks.last_played_date = $4::date THEN streaks.current_streak + 1
		     WHEN streaks.last_played_date = $2::date THEN streaks.current_streak  -- already played today
		     ELSE 1
		   END,
		   longest_streak = GREATEST(
		     streaks.longest_streak,
		     CASE
		       WHEN streaks.last_played_date = $4::date THEN streaks.current_streak + 1
		       WHEN streaks.last_played_date = $2::date THEN streaks.current_streak
		       ELSE 1
		     END
		   ),
		   last_played_date = CASE
		     WHEN streaks.last_played_date = $2::date THEN streaks.last_played_date
		     ELSE $2::date
		   END,
		   total_points = streaks.total_points + $3,
		   updated_at = NOW()
		 RETURNING current_streak`,
		deviceID, todayStr, pointsEarned, yesterdayStr,
	).Scan(&newStreak)
	return newStreak, err
}
