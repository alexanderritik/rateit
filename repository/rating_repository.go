package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ritikmayank/catch-it/model"
)

type RatingRepository struct {
	db *pgxpool.Pool
}

func NewRatingRepository(db *pgxpool.Pool) *RatingRepository {
	return &RatingRepository{db: db}
}

func (r *RatingRepository) Upsert(ctx context.Context, deviceID string, clipID, score int, listened bool) (*model.Rating, error) {
	var rt model.Rating
	err := r.db.QueryRow(ctx,
		`INSERT INTO ratings (device_id, clip_id, score, listened)
		 VALUES ($1, $2, $3, $4)
		 ON CONFLICT (device_id, clip_id) DO UPDATE
		   SET score=$3, listened=ratings.listened OR $4, rated_at=NOW()
		 RETURNING id, device_id, clip_id, score, listened, rated_at`,
		deviceID, clipID, score, listened,
	).Scan(&rt.ID, &rt.DeviceID, &rt.ClipID, &rt.Score, &rt.Listened, &rt.RatedAt)
	return &rt, err
}

func (r *RatingRepository) GetByDeviceAndClip(ctx context.Context, deviceID string, clipID int) (*model.Rating, error) {
	var rt model.Rating
	err := r.db.QueryRow(ctx,
		`SELECT id, device_id, clip_id, score, listened, rated_at
		 FROM ratings WHERE device_id=$1 AND clip_id=$2`,
		deviceID, clipID,
	).Scan(&rt.ID, &rt.DeviceID, &rt.ClipID, &rt.Score, &rt.Listened, &rt.RatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &rt, err
}

// RatedClipIDs returns the set of clip IDs already rated by a device (for a list of clip IDs).
func (r *RatingRepository) RatedClipIDs(ctx context.Context, deviceID string, clipIDs []int) (map[int]model.Rating, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, device_id, clip_id, score, listened, rated_at
		 FROM ratings WHERE device_id=$1 AND clip_id = ANY($2)`,
		deviceID, clipIDs,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	m := make(map[int]model.Rating)
	for rows.Next() {
		var rt model.Rating
		if err := rows.Scan(&rt.ID, &rt.DeviceID, &rt.ClipID, &rt.Score, &rt.Listened, &rt.RatedAt); err != nil {
			return nil, err
		}
		m[rt.ClipID] = rt
	}
	return m, rows.Err()
}

// CountRatedToday returns how many clips a device has rated from a given list.
func (r *RatingRepository) CountRatedToday(ctx context.Context, deviceID string, clipIDs []int) (int, error) {
	var count int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM ratings WHERE device_id=$1 AND clip_id = ANY($2)`,
		deviceID, clipIDs,
	).Scan(&count)
	return count, err
}

// CountTotal returns total ratings by a device ever.
func (r *RatingRepository) CountTotal(ctx context.Context, deviceID string) (int, error) {
	var count int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM ratings WHERE device_id=$1`, deviceID,
	).Scan(&count)
	return count, err
}

// AllListenedToday returns true if all clips for a day were listened to before rating.
func (r *RatingRepository) AllListenedToday(ctx context.Context, deviceID string, clipIDs []int) (bool, error) {
	var count int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM ratings
		 WHERE device_id=$1 AND clip_id = ANY($2) AND listened=true`,
		deviceID, clipIDs,
	).Scan(&count)
	return count == len(clipIDs), err
}

type LeaderboardRow struct {
	DeviceID    string
	TotalRated  int
	TotalPoints int
}

// Leaderboard returns top N raters by total points.
func (r *RatingRepository) Leaderboard(ctx context.Context, limit int) ([]LeaderboardRow, error) {
	rows, err := r.db.Query(ctx,
		`SELECT s.device_id, COUNT(rt.id)::int, s.total_points
		 FROM streaks s
		 LEFT JOIN ratings rt ON rt.device_id = s.device_id
		 GROUP BY s.device_id, s.total_points
		 ORDER BY s.total_points DESC
		 LIMIT $1`,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []LeaderboardRow
	for rows.Next() {
		var entry LeaderboardRow
		if err := rows.Scan(&entry.DeviceID, &entry.TotalRated, &entry.TotalPoints); err != nil {
			return nil, err
		}
		result = append(result, entry)
	}
	return result, rows.Err()
}
