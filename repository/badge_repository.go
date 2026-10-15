package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ritikmayank/catch-it/model"
)

type BadgeRepository struct {
	db *pgxpool.Pool
}

func NewBadgeRepository(db *pgxpool.Pool) *BadgeRepository {
	return &BadgeRepository{db: db}
}

// Award grants a badge to a device. Returns true if newly awarded, false if already had it.
func (r *BadgeRepository) Award(ctx context.Context, deviceID, badgeID string) (bool, error) {
	tag, err := r.db.Exec(ctx,
		`INSERT INTO badges (device_id, badge_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		deviceID, badgeID,
	)
	if err != nil {
		return false, err
	}
	return tag.RowsAffected() > 0, nil
}

// ListByDevice returns all badges a device has earned.
func (r *BadgeRepository) ListByDevice(ctx context.Context, deviceID string) ([]model.BadgeInfo, error) {
	rows, err := r.db.Query(ctx,
		`SELECT badge_id, earned_at FROM badges WHERE device_id=$1 ORDER BY earned_at ASC`,
		deviceID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.BadgeInfo
	for rows.Next() {
		var badgeID string
		var earnedAt time.Time
		if err := rows.Scan(&badgeID, &earnedAt); err != nil {
			return nil, err
		}
		def, ok := model.BadgeCatalog[badgeID]
		if !ok {
			continue
		}
		result = append(result, model.BadgeInfo{
			ID:       def.ID,
			Name:     def.Name,
			Icon:     def.Icon,
			Desc:     def.Desc,
			EarnedAt: earnedAt,
		})
	}
	return result, rows.Err()
}

// CountCorrect returns total correct plays for a device.
func (r *BadgeRepository) CountCorrect(ctx context.Context, deviceID string) (int, error) {
	var count int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM plays WHERE device_id=$1 AND correct=true`,
		deviceID,
	).Scan(&count)
	return count, err
}
